package state

// GetTop 返回栈顶索引
func (self *luaState) GetTop() int {
	return self.stack.top
}

// AbsIndex 将索引转化为绝对索引
func (self *luaState) AbsIndex(idx int) int {
	return self.stack.absIndex(idx)
}

// CheckStack 检查是否还能再推n个元素入栈，这里忽略了扩容失败的情况，不进行处理
func (self *luaState) CheckStack(n int) bool {
	self.stack.check(n)
	// TODO 处理扩容失败的情况
	return true
}

// Pop 弹出n个元素
func (self *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		self.stack.pop()
	}
}

// Copy 把值从一个位置复制到另一个位置
func (self *luaState) Copy(fromIdx, toIdx int) {
	val := self.stack.get(fromIdx)
	self.stack.set(toIdx, val)
}

// PushValue 把指定索引处的值推入栈顶，会把指定索引处的值复制并推入
func (self *luaState) PushValue(idx int) {
	val := self.stack.get(idx)
	self.stack.push(val)
}

// Replace 弹出栈顶值并将值写入到指定的位置
func (self *luaState) Replace(idx int) {
	val := self.stack.pop()
	self.stack.set(idx, val)
}

// Insert 将栈顶值弹出，插入到指定位置，实际上就是旋转的特例
func (self *luaState) Insert(idx int) {
	self.Rotate(idx, 1)
}

// Remove 删除指定索引处的值，然后该值上面的所有元素全部向下移动一个位置
func (self *luaState) Remove(idx int) {
	self.Rotate(idx, -1)
	self.Pop(1)
}

// Rotate 旋转
func (self *luaState) Rotate(idx, n int) {
	t := self.stack.top - 1
	p := self.stack.absIndex(idx) - 1
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	self.stack.reverse(p, m)
	self.stack.reverse(m+1, t)
	self.stack.reverse(p, t)
}

// SetTop 把栈顶索引设定为指定值，如果值小于当前栈索引，那么会不断Pop，如果大于，那么会Push nil直到栈索引等于给定值
func (self *luaState) SetTop(idx int) {
	newTop := self.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}

	n := self.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			self.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			self.stack.push(nil)
		}
	}
}
