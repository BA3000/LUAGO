package vm

import . "luago/api"

const LFIELDS_PER_FLUSH = 50

// newTable 创建 table R(A) := {} (size = B, C)
func newTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.CreateTable(Fb2int(b), Fb2int(c))
	// 将新的 table 放入到指定索引处
	vm.Replace(a)
}

// getTable 根据键从表里取值，并放入目标寄存器 R(A) := R(B)[RK(C)]
func getTable(i Instruction, vm LuaVM)  {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// setTable 根据键往表里赋值，表在寄存器，操作数 A 指定表的索引，键值可能在寄存器，可能在常量表，由操作数 B、C 指定
func setTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(a)
}

// R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
func setList(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	if c > 0 {
		c = c - 1
	} else {
		c = Instruction(vm.Fetch()).Ax()
	}

	idx := int64(c * LFIELDS_PER_FLUSH)
	for j := 1; j <= b; j++ {
		idx++
		// 因为数据紧跟着表，连续排列，所以可以 a + j 获取当前要设置的数据，并推入栈顶
		vm.PushValue(a + j)
		// 将栈顶的数据设置到表的指定索引处
		vm.SetI(a, idx)
	}
}
