package state

type luaState struct {
	stack	*luaStack
}

// New 创建luaState实例
func New() *luaState {
	return &luaState{
		stack:	newLuaStack(20),
	}
}

// pushLuaStack 相当于单向链表的 push
func (lState *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = lState.stack
	lState.stack = stack
}

// popLuaStack 相当于单向链表的 pop
func (lState *luaState) popLuaStack() {
	stack := lState.stack
	lState.stack = stack.prev
	stack.prev = nil
}
