package state

import (
	. "luago/api"
)

// Compare 按索引得到两个操作数，然后根据操作码调用对应的函数进行运算
func (lState *luaState) Compare(idx1, idx2 int, op CompareOp) bool {
	a := lState.stack.get(idx1)
	b := lState.stack.get(idx2)
	switch op {
	case LUA_OPEQ:
		return _eq(a, b)
	case LUA_OPLT:
		return _lt(a, b)
	case LUA_OPLE:
		return _le(a, b)
	default:
		panic("invalid compare op!")
	}
}

// _eq 用来判断两个操作数是否相等的函数
func _eq(a, b luaValue) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		default:
			return false
		}
	default:
		return a == b

	}
}

// _lt 用来比较两个操作数大小，判断 a 是否小于 b 的函数
func _lt(a, b luaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x < y
		case int64:
			return x < float64(y)
		}
	}
	panic("comparison error!")
}

// _le 用来比较两个操作数大小，判断 a 是否小于等于 b 的函数
func _le(a, b luaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x <= y
		case int64:
			return x <= float64(y)
		}
	}
	panic("comparison error!")
}