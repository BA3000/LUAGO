package state

// Len 获取指定索引上内容的长度，并将结果入栈
func (lState *luaState) Len(idx int) {
	val := lState.stack.get(idx)
	if s, ok := val.(string); ok {
		lState.stack.push(int64(len(s)))
	} else {
		panic("length error!")
	}
}

// Concat 从栈顶弹出 n 个值，拼接这些值之后将结果入栈
func (lState *luaState) Concat(n int) {
	if n == 0 {
		lState.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if lState.IsString(-1) && lState.IsString(-2) {
				s2 := lState.ToString(-1)
				s1 := lState.ToString(-2)
				lState.stack.pop()
				lState.stack.pop()
				lState.stack.push(s1 + s2)
				continue
			}
			panic("concatenation error!")
		}
	}
	// n == 1, do nothing
}