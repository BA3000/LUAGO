package vm

type Instruction uint32

const MAXARG_Bx = 1<<18 - 1			// 2^18 - 1 = 262143
const MAXARG_sBx = MAXARG_Bx >> 1	// 262143 / 2 = 131071

// Opcode 提取操作码
func (self Instruction) Opcode() int {
	return int(self & 0x3F)
}

// ABC 从iABC模式中提取参数
func (self Instruction) ABC() (a, b, c int) {
	a = int(self >> 6 & 0xFF)
	c = int (self >> 14 & 0x1FF)
	b = int(self >> 23 & 0x1FF)
	return
}

// ABx 从iABx模式指令中提取参数
func (self Instruction) ABx() (a, bx int) {
	a = int(self >> 6 & 0xFF)
	bx = int(self >> 14)
	return
}

// AsBx 从iAsBx模式指令中提取参数
func (self Instruction) AsBx() (a, sbx int) {
	a, bx := self.ABx()
	return a, bx - MAXARG_sBx
}

// Ax 从iAx模式指令中提取参数
func (self Instruction) Ax() int {
	return int(self >> 6)
}

// OpName 获取操作码名字
func (self Instruction) OpName() string {
	return opcodes[self.Opcode()].name
}

// OpMode 获取编码模式
func (self Instruction) OpMode() byte {
	return opcodes[self.Opcode()].opMode
}

// BMode 获取操作数B的使用模式
func (self Instruction) BMode() byte {
	return opcodes[self.Opcode()].argBMode
}

// CMode 获取操作数C的使用模式
func (self Instruction) CMode() byte {
	return opcodes[self.Opcode()].argCMode
}
