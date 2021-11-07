package state

import . "luago/api"

// CreateTable 创建 Lua Table 并压入栈，两个参数用于指定创建的数组、哈希表大小
func (lState *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	lState.stack.push(t)
}

// NewTable 创建 Table
func (lState *luaState) NewTable() {
	lState.CreateTable(0, 0)
}

// GetTable 根据栈顶的键来从表（索引指定）里获取值，最后把值推入栈顶，并返回取到的值的类型
func (lState *luaState) GetTable(idx int) LuaType {
	t := lState.stack.get(idx)
	k := lState.stack.pop()
	return lState.getTable(t, k)
}

// getTable 取值并返回值的类型
func (lState *luaState) getTable(t, k luaValue) LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		lState.stack.push(v)
		return typeOf(v)
	}
	panic("not a table!")
}

// GetField 传入的字符串作为键，获取保存的数据（推入栈顶）
func (lState *luaState) GetField(idx int, k string) LuaType {
	t := lState.stack.get(idx)
	return lState.getTable(t, k)
}

// GetI 根据传入的索引参数获取数组数据，并压入栈顶
func (lState *luaState) GetI(idx int, i int64) LuaType {
	t := lState.stack.get(idx)
	return lState.getTable(t, i)
}
