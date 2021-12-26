package state

type luaStack struct {
	// 存放值
	slots	[]luaValue
	// 栈顶索引
	top		int
	prev	*luaStack
	closure	*closure
	varargs	[]luaValue
	pc		int
}

// newLuaStack 创建指定容量的栈
func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots:	make([]luaValue, size),
		top:	0,
	}
}

// check 检查栈的空闲空间是否还可以推入n个值，如果不能就扩容
func (lStack *luaStack) check(n int) {
	free := len(lStack.slots) - lStack.top
	for i := free; i < n; i++ {
		lStack.slots = append(lStack.slots, nil)
	}
}

// push 入栈
func (lStack *luaStack) push(val luaValue) {
	if lStack.top == len(lStack.slots) {
		panic("stack overflow!")
	}
	lStack.slots[lStack.top] = val
	lStack.top++
}

// pop 弹出栈
func (lStack *luaStack) pop() luaValue {
	if lStack.top < 1 {
		panic("stack underflow!")
	}
	lStack.top--
	val := lStack.slots[lStack.top]
	lStack.slots[lStack.top] = nil
	return val
}

// absIndex 将相对索引转换成绝对索引
func (lStack *luaStack) absIndex(idx int) int {
	if idx >= 0 {
		return idx
	}
	return idx + lStack.top + 1
}

// isValid 检查索引是否有效
func (lStack *luaStack) isValid(idx int) bool {
	absIdx := lStack.absIndex(idx)
	return absIdx > 0 && absIdx <= lStack.top
}

// get 根据索引从栈取值，如果无效就返回nil
func (lStack *luaStack) get(idx int) luaValue {
	absIdx := lStack.absIndex(idx)
	if absIdx > 0 && absIdx <= lStack.top {
		return lStack.slots[absIdx-1]
	}
	return nil
}

// set 往栈写入值，索引无效会panic
func (lStack *luaStack) set(idx int, val luaValue) {
	absIdx := lStack.absIndex(idx)
	if absIdx > 0 && absIdx <= lStack.top {
		lStack.slots[absIdx-1] = val
		return
	}
	panic("invalid index!")
}

// reverse 反转
func (lStack *luaStack) reverse(from, to int) {
	slots := lStack.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}