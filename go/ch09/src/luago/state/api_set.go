package state

// SetTable 把键值对写入表，参数输入的是表在栈中的索引，键值对会通过 pop 从栈获取
func (lState *luaState) SetTable(idx int) {
	t := lState.stack.get(idx)
	v := lState.stack.pop()
	k := lState.stack.pop()
	lState.setTable(t, k, v)
}

func (lState *luaState) setTable(t, k, v luaValue) {
	if tbl, ok := t.(*luaTable); ok {
		tbl.put(k, v)
		return
	}
	panic("not a table!")
}

// SetField 传入字符串键值，用来给记录的字段赋值，要赋的值在栈顶
func (lState *luaState) SetField(idx int, k string) {
	t := lState.stack.get(idx)
	v := lState.stack.pop()
	lState.setTable(t, k, v)
}

// SetI 设置表传入参数 i 索引上的值，要设的值要在栈顶
func (lState *luaState) SetI(idx int, i int64) {
	t := lState.stack.get(idx)
	v := lState.stack.pop()
	lState.setTable(t, i, v)
}
