package state

func (lState *luaState) PC() int {
	return lState.pc
}

func (lState *luaState) AddPC(n int) {
	lState.pc += n
}

// Fetch 从函数原型的指令表中提取出当前指令，并把PC加1
func (lState *luaState) Fetch() uint32 {
	i := lState.proto.Code[lState.pc]
	lState.pc++
	return i
}

// GetConst 从函数原型的常量表中提取出一个常量值，然后推入栈
func (lState *luaState) GetConst(idx int) {
	c := lState.proto.Constants[idx]
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
