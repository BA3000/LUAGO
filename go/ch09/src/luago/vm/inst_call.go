package vm

import . "luago/api"

// closure 加载子函数原型
func closure(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	vm.LoadProto(bx)
	vm.Replace(a)
}

// call 实现 CALL 指令
func call(i Instruction, vm LuaVM)  {
	a, b, c := i.ABC()
	a += 1

	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c - 1)
	// 因为返回值已经在栈顶，所以可以直接调用此函数将返回值移动到指定的寄存器中
	_popResults(a, c, vm)
}

// _pushFuncAndArgs 将函数和参数压入栈
func _pushFuncAndArgs(a, b int, vm LuaVM) (nArgs int) {
	if b >= 1 {
		// b-1 args
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b-1
	} else {
		// b == 0
		_fixStack(a, vm)
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

func _popResults(a, c int, vm LuaVM) {
	if c == 1 {
		// no results
	} else if c >1 {
		// c - 1 results
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else {
		// c 为0，吧所有的返回值返回，直接先留在栈顶，留下 a 标记要移动到哪些寄存器中
		// 例如 f(1, 2, g()) 里面 g 就是所有的返回值都要留下，传给 f
		vm.CheckStack(1)
		vm.PushInteger(int64(a))
	}
}

// _fixStack 处理需要把所有返回值作为参数传入到另一个函数的情况，先将函数和前部分的参数压入到前一个函数的返回值上方（即后半部分的参数）
// 然后旋转栈顶，整理顺序，a 是寄存器索引，函数 proto 所在的位置
func _fixStack(a int, vm LuaVM) {
	// 取出栈顶的值，CALL 执行后压入的返回值的数量
	x := int(vm.ToInteger(-1))
	vm.Pop(1)

	vm.CheckStack(x - 1)
	// 从索引 a 开始复制 x 个参数，压入栈
	// 实际上是用来复制函数和前半部分参数
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	// RegisterCount返回当前函数需要使用的寄存器数量，这部操作是将最下方的后半部分参数移动到最上方，整理顺序
	// 最终变成从底部到上方：函数-前半部分参数-后半部分参数
	vm.Rotate(vm.RegisterCount()+1, x-a)
}

func _return(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b == 1 {
		// no return values
	} else if b > 1 {
		// b-1 return value
		vm.CheckStack(b-1)
		for i := a; i <= a+b-2; i++ {
			vm.PushValue(i)
		}
		// 一部分返回值已经在栈顶，调用此函数将另一部分推入即可
		_fixStack(a, vm)
	}
}

// vararg 把传递给当前函数的变长参数加载到连续多个寄存器中
func vararg(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b != 1 {
		// b > 1的时候会把 b - 1 个 vararg 参数用 LoadVararg 复制到寄存器
		// 如果等于 0，那么会把全部 vararg 参数用 LoadVararg 入栈
		vm.LoadVararg(b-1)
		_popResults(a, b, vm)
	}
}

// tailCall TODO 尾递归优化
func tailCall(i Instruction, vm LuaVM) {
	// TODO 优化尾递归，这里只是简单的进行了函数调用
	a, b, _ := i.ABC()
	a += 1
	c := 0

	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

// self 把对象和方法拷贝到相邻的两个目标寄存器中
func self(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1; b += 1

	// 把表copy到a+1索引处
	vm.Copy(b, a + 1)
	// 方法名在常量表，所以用GetRK获取，取得的常量入栈
	vm.GetRK(c)
	// 根据栈顶的方法名从b指定的表获取值，值放入栈顶
	vm.GetTable(b)
	// 栈顶的 obj.f 放到 a 处
	vm.Replace(a)
	// 最后从顶部开始向底部内容的顺序变成 obj （a + 1处）, obj.f （a 处）
}
