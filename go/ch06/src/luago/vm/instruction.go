package vm

import "luago/api"

type Instruction uint32

const MAXARG_Bx = 1<<18 - 1			// 2^18 - 1 = 262143
const MAXARG_sBx = MAXARG_Bx >> 1	// 262143 / 2 = 131071

// Opcode 提取操作码
func (i Instruction) Opcode() int {
	return int(i & 0x3F)
}

// ABC 从iABC模式中提取参数
func (i Instruction) ABC() (a, b, c int) {
	a = int(i >> 6 & 0xFF)
	c = int (i >> 14 & 0x1FF)
	b = int(i >> 23 & 0x1FF)
	return
}

// ABx 从iABx模式指令中提取参数
func (i Instruction) ABx() (a, bx int) {
	a = int(i >> 6 & 0xFF)
	bx = int(i >> 14)
	return
}

// AsBx 从iAsBx模式指令中提取参数
func (i Instruction) AsBx() (a, sbx int) {
	a, bx := i.ABx()
	return a, bx - MAXARG_sBx
}

// Ax 从iAx模式指令中提取参数
func (i Instruction) Ax() int {
	return int(i >> 6)
}

// OpName 获取操作码名字
func (i Instruction) OpName() string {
	return opcodes[i.Opcode()].name
}

// OpMode 获取编码模式
func (i Instruction) OpMode() byte {
	return opcodes[i.Opcode()].opMode
}

// BMode 获取操作数B的使用模式
func (i Instruction) BMode() byte {
	return opcodes[i.Opcode()].argBMode
}

// CMode 获取操作数C的使用模式
func (i Instruction) CMode() byte {
	return opcodes[i.Opcode()].argCMode
}

// Execute 从指令提取操作码并查找对应的指令实现方法，最后调用指令实现方法执行指令
func (i Instruction) Execute(vm api.LuaVM) {
	action := opcodes[i.Opcode()].action
	if action != nil {
		action(i, vm)
	} else {
		panic(i.OpName())
	}
}
