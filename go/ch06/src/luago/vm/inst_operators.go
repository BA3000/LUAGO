package vm

import . "luago/api"

// _binaryArith 用来实现二元算数运算指令
func _binaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.Arith(op)
	vm.Replace(a)
}

// _unaryArith 用来实现一元算数运算指令
func _unaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, _ := i.ABC()
	a += 1; b += 1

	vm.PushValue(b)
	vm.Arith(op)
	vm.Replace(a)
}

func add(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPADD)
}

func sub(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPSUB)
}

func mul(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPMUL)
}

func mod(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPMOD)
}

func pow(i Instruction, vm LuaVM)  {
	_binaryArith(i, vm, LUA_OPPOW)
}

func div(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPDIV)
}

func idiv(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPIDIV)
}

func band(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPBAND)
}

func bor(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPBOR)
}

func bxor(i Instruction, vm LuaVM)  {
	_binaryArith(i, vm, LUA_OPBXOR)
}

func shl(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPSHL)
}

func shr(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPSHR)
}

func unm(i Instruction, vm LuaVM) {
	_unaryArith(i, vm, LUA_OPUNM)
}

func bnot(i Instruction, vm LuaVM) {
	_unaryArith(i, vm, LUA_OPBNOT)
}

// length Lua中的#
func length(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1; b+= 1

	vm.Len(b)
	vm.Replace(a)
}

// Lua的 ..
func concat(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1; b += 1; c += 1

	n := c - b + 1
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		vm.PushValue(i)
	}
	vm.Concat(n)
	vm.Replace(a)
}

func _compare(i Instruction, vm LuaVM, op CompareOp) {
	a, b, c := i.ABC()

	vm.GetRK(b)
	vm.GetRK(c)
	if vm.Compare(-2, -1, op) != (a != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2)
}

// eq ==
func eq(i Instruction, vm LuaVM) {
	_compare(i, vm, LUA_OPEQ)
}

// lt <
func lt(i Instruction, vm LuaVM) {
	_compare(i, vm, LUA_OPLT)
}

// le <=
func le(i Instruction, vm LuaVM) {
	_compare(i, vm, LUA_OPLE)
}

// not 非运算 !
func not(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1; b += 1

	vm.PushBoolean(!vm.ToBoolean(b))
	vm.Replace(a)
}

// testSet 判断b索引上的值是否和给定的布尔值c相同，如果是那么将b的值复制到a索引指向的地址，如果不是那么PC++，跳过下一条指令
func testSet(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1; b += 1

	if vm.ToBoolean(b) == (c != 0) {
		vm.Copy(b, a)
	} else {
		vm.AddPC(1)
	}
}

// test 判断寄存器A中的值转换为布尔值后是否和操作数C表示的布尔值一致，不一致则跳过下一条指令
func test(i Instruction, vm LuaVM) {
	a, _, c := i.ABC()
	a += 1

	if vm.ToBoolean(a) != (c != 0) {
		vm.AddPC(1)
	}
}

func forPrep(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1

	// R(A) -= R(A+2)
	vm.PushValue(a)
	vm.PushValue(a + 2)
	vm.Arith(LUA_OPSUB)
	vm.Replace(a)
	// pc += sBx
	vm.AddPC(sBx)
}

func forLoop(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1

	// R(A) += R(A+2);
	vm.PushValue(a + 2)
	vm.PushValue(a)
	vm.Arith(LUA_OPADD)
	vm.Replace(a)

	// R(A) <?= R(A+1)
	isPositiveStep := vm.ToNumber(a + 2) >= 0
	if isPositiveStep && vm.Compare(a, a+1, LUA_OPLE) || !isPositiveStep && vm.Compare(a+1, a, LUA_OPLE) {
		vm.AddPC(sBx)	// pc+=sBx
		vm.Copy(a, a+3)	// R(A+3)=R(A), 数值拷贝给用户定义的局部变量
	}
}
