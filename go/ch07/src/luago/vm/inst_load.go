package vm

import . "luago/api"

// loadNil 给连续n个局部变量设置Nil
func loadNil(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	vm.PushNil()
	for i := a; i <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

// loadBool 给单个寄存器设置布尔值，索引用操作数A指定，布尔值用寄存器B指定（0 false，非0 true），C非0那么跳过下一条指令
func loadBool(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	vm.PushBoolean(b != 0)
	vm.Replace(a)
	if c != 0 {
		vm.AddPC(1)
	}
}

// loadK 将常量表中某个常量加载到指定寄存器，寄存器索引用操作数A指定，常量表索引用操作数Bx指定
func loadK(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	vm.GetConst(bx)
	vm.Replace(a)
}

// loadKx 加载常量，只是范围更广，借助Ax 26位操作数支持更大的常量表
func loadKx(i Instruction, vm LuaVM) {
	a, _ := i.ABx()
	a += 1
	ax := Instruction(vm.Fetch()).Ax()

	vm.GetConst(ax)
	vm.Replace(a)
}
