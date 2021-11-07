package state

import "luago/binchunk"

type luaState struct {
	stack	*luaStack
	proto	*binchunk.Prototype
	// 程序计数器
	pc		int
}

// New 创建luaState实例
func New(stackSize int, proto *binchunk.Prototype) *luaState {
	return &luaState{
		stack:	newLuaStack(stackSize),
		proto:	proto,
		pc:		0,
	}
}
