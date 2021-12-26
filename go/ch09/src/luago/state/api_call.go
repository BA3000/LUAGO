package state

import "luago/binchunk"
import "luago/vm"

// Load 加载 chunk 用
func (lState *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk)
	c := newLuaClosure(proto)
	lState.stack.push(c)
	return 0
}

// Call 调用Lua函数，nArgs 参数数量，nResults 返回值数量
func (lState *luaState) Call(nArgs, nResults int)  {
	val := lState.stack.get(-(nArgs + 1))
	if c, ok := val.(*closure); ok {
		if c.proto != nil {
			// 调用Lua函数
			lState.callLuaClosure(nArgs, nResults, c)
		} else {
			// 调用Go函数
			lState.callGoClosure(nArgs, nResults, c)
		}
	} else {
		panic("try to call a not function!")
	}
}

func (lState *luaState) callLuaClosure(nArgs, nResults int, c *closure)  {
	// 执行函数所需要的寄存器数量
	nRegs := int(c.proto.MaxStackSize)
	// 定义函数时声明的固定参数数量
	nParams := int(c.proto.NumParams)
	// 是不是 varargs 函数
	isVararg := c.proto.IsVararg == 1

	newStack := newLuaStack(nRegs + 20)
	newStack.closure = c

	// 弹出函数和函数的参数
	funcAndArgs := lState.stack.popN(nArgs + 1)
	// 传入参数
	newStack.pushN(funcAndArgs[1:], nParams);
	newStack.top = nRegs
	if nArgs > nParams && isVararg {
		// 如果是 varargs 函数并且传入参数比给定的参数数量多，那么就要记录下剩下的参数
		newStack.varargs = funcAndArgs[nParams + 1:]
	}

	// 压入新调用帧
	lState.pushLuaStack(newStack)
	// 执行被调函数
	lState.runLuaClosure()
	// 执行完成，弹出
	lState.popLuaStack()

	if nResults != 0 {
		// 弹出执行结果
		results := newStack.popN(newStack.top - nRegs)
		// 检查空间，看还够不够空间入栈
		lState.stack.check(len(results))
		// 将结果推入当前栈帧
		lState.stack.pushN(results, nResults)
	}
}

// callGoClosure 与callLuaClosure类似，先创建新的调用帧，然后把参数值从主调帧里弹出，推入被调帧，Go闭包从主调用帧弹出并抛弃，参数传递完成
// 之后被调帧推入调用栈成为当前帧，最后执行Go函数。执行完毕后弹出被调帧，让主调帧重新成为当前帧。有返回值的话还会把返回值从被调帧弹出的推入主调帧
func (lState *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	// 创建新的调用帧
	newStack := newLuaStack(nArgs + 20)
	newStack.closure = c

	// 弹出参数
	args := lState.stack.popN(nArgs)
	// 将参数推入到新的调用帧
	newStack.pushN(args, nArgs)
	// 弹出Go闭包
	lState.stack.pop()

	// 推入被调用帧
	lState.pushLuaStack(newStack)
	// 执行go函数
	r := c.goFunc(lState)
	// 弹出被调帧
	lState.popLuaStack()

	if nResults != 0 {
		// 有返回值的情况，从被调帧弹出返回值，并且压入到主调帧
		results := newStack.popN(r)
		lState.stack.check(len(results))
		lState.stack.pushN(results, nResults)
	}
}

// runLuaClosure 运行 Lua 闭包代码
func (lState *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(lState.Fetch())
		inst.Execute(lState)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

// popN 从栈连续弹出 n 个值
func (lStack *luaStack) popN(n int) []luaValue {
	vals := make([]luaValue, n)
	for i := n-1; i >= 0; i-- {
		vals[i] = lStack.pop()
	}
	return vals
}

// pushN vals 数组，n 要压入的元素数量，遵循多退少补原则
func (lStack *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 { n = nVals }
	for i := 0; i < n; i++ {
		if i < nVals {
			lStack.push(vals[i])
		} else {
			lStack.push(nil)
		}
	}
}
