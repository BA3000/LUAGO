package state

func (lState *luaState) PC() int {
	return lState.stack.pc
}

func (lState *luaState) AddPC(n int) {
	lState.stack.pc += n
}

// Fetch 从函数原型的指令表中提取出当前指令，并把PC加1
func (lState *luaState) Fetch() uint32 {
	i := lState.stack.closure.proto.Code[lState.stack.pc]
	lState.stack.pc++
	return i
}

// GetConst 从函数原型的常量表中提取出一个常量值，然后推入栈
func (lState *luaState) GetConst(idx int) {
	c := lState.stack.closure.proto.Constants[idx]
	lState.stack.push(c)
}

// GetRK 看情况调用GetConst把常量入栈，或者调用PushValue把栈值入栈
func (lState *luaState) GetRK(rk int) {
	if rk > 0xFF {
		// constant
		lState.GetConst(rk & 0xFF)
	} else {
		// register
		lState.PushValue(rk + 1)
	}
}

// RegisterCount 返回当前 Lua 函数所操作的寄存器数量
func (lState *luaState) RegisterCount() int {
	return int(lState.stack.closure.proto.MaxStackSize)
}

// LoadVararg 把传递给当前 Lua 函数的变长参数推入栈顶（多退少补）
func (lState *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(lState.stack.varargs)
	}
	lState.stack.check(n)
	lState.stack.pushN(lState.stack.varargs, n)
}

func (lState *luaState) LoadProto(idx int) {
	proto := lState.stack.closure.proto.Protos[idx]
	closure := newLuaClosure(proto)
	lState.stack.push(closure)
}
