package state

import "math"
import "luago/number"

type luaTable struct {
	arr		[]luaValue
	_map	map[luaValue]luaValue
}

// newLuaTable 创建一个空的Lua Table
func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0 {
		t.arr = make([]luaValue, 0, nArr)
	}
	if nRec > 0 {
		t._map = make(map[luaValue]luaValue, nRec)
	}
	return t
}

// get 根据键从表中查找值，如果键是整数（或者可以被转化成整数）且在数组索引范围内，那么就会按照索引访问数组；否则会从哈希表查找值
func (lT *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(lT.arr)) {
			return lT.arr[idx - 1]
		}
	}
	return lT._map[key]
}

// 尝试把浮点数转换成整数
func _floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

// put 往表里存入键值对
func (lT *luaTable) put(key, val luaValue)  {
	if key == nil {
		panic("table index is nil!")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN!")
	}
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(lT.arr))
		if idx <= arrLen {
			lT.arr[idx - 1] = val
			if idx == arrLen && val == nil {
				// 新添加的值是nil且在数组最末端，那么可以调用函数删除掉尾部的洞
				lT._shrinkArray()
			}
			return
		}
		if idx == arrLen + 1 {
			// 避免哈希表中存在同样的键
			delete(lT._map, key)
			if val != nil {
				lT.arr = append(lT.arr, val)
				// 动态扩展数组
				lT._expandArray()
			}
			return
		}
	}
	if val != nil {
		// 确保写入的时候哈希表已经被创建
		if lT._map == nil {
			lT._map = make(map[luaValue]luaValue, 8)
		}
		lT._map[key] = val
	} else {
		// 如果要插入的值是nil那么就会删除掉这个键
		delete(lT._map, key)
	}
}

// _shrinkArray 将数组中值为 nil 的部分移除，删除掉数组中的洞
func (lT *luaTable) _shrinkArray() {
	for i := len(lT.arr) - 1; i >= 0; i-- {
		if lT.arr[i] == nil {
			lT.arr = lT.arr[0:i]
		}
	}
}

// _expandArray 数组部分动态拓展后，要将原本存在哈希表中的值移动到数组中
func (lT *luaTable) _expandArray() {
	for idx := int64(len(lT.arr)) + 1; true; idx++ {
		if val, found := lT._map[idx]; found {
			delete(lT._map, idx)
			lT.arr = append(lT.arr, val)
		} else {
			break
		}
	}
}

// len 返回数组长度
func (lT *luaTable) len() int {
	return len(lT.arr)
}
