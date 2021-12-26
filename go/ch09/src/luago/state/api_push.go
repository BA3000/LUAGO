package state

import . "luago/api"

func (lState *luaState) PushNil()             { lState.stack.push(nil) }
func (lState *luaState) PushBoolean(b bool)   { lState.stack.push(b) }
func (lState *luaState) PushInteger(n int64)  { lState.stack.push(n) }
func (lState *luaState) PushNumber(n float64) { lState.stack.push(n) }
func (lState *luaState) PushString(s string)  { lState.stack.push(s) }

// PushGoFunction 接受一个Go函数参数，转变成Go闭包之后压入栈
func (lState *luaState) PushGoFunction(f GoFunction) {
	lState.stack.push(newGoClosure(f))
}
