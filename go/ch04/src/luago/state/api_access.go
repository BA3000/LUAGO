package state

import "fmt"
import . "luago/api"

// TypeName 将Lua类型转换成对应的字符串表示
func (lState *luaState) TypeName(tp LuaType) string {
	switch tp {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

// Type 根据索引返回值数据类型（go类型）返回对应的Lua类型
func (lState *luaState) Type(idx int) LuaType {
	if lState.stack.isValid(idx) {
		val := lState.stack.get(idx)
		return typeOf(val)
	}
	return LUA_TNONE
}

// IsNone 判断指定索引处的值是不是None类型
func (lState *luaState) IsNone(idx int) bool {
	return lState.Type(idx) == LUA_TNONE
}

// IsNil 判断指定索引处的值是不是Nil类型
func (lState *luaState) IsNil(idx int) bool {
	return lState.Type(idx) == LUA_TNIL
}

// IsNoneOrNil 判断是不是None或者Nil类型
func (lState *luaState) IsNoneOrNil(idx int) bool {
	return lState.Type(idx) <= LUA_TNIL
}

// IsBoolean 判断是不是Bool类型
func (lState *luaState) IsBoolean(idx int) bool {
	return lState.Type(idx) == LUA_TBOOLEAN
}

func (lState *luaState) IsTable(idx int) bool {
	return lState.Type(idx) == LUA_TTABLE
}

func (lState *luaState) IsFunction(idx int) bool {
	return lState.Type(idx) == LUA_TFUNCTION
}

func (lState *luaState) IsThread(idx int) bool {
	return lState.Type(idx) == LUA_TTHREAD
}

// IsString 判断索引处是不是string或者数字类型
func (lState *luaState) IsString(idx int) bool {
	t := lState.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

// IsNumber 判定索引处的值是不是数字或者能不能够被转换成数字
func (lState *luaState) IsNumber(idx int) bool {
	_, ok := lState.ToNumberX(idx)
	return ok
}

// IsInteger 判断索引处的值是否是整数类型
func (lState *luaState) IsInteger(idx int) bool {
	val := lState.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

// ToBoolean 将指定索引上的值转换为bool
func (lState *luaState) ToBoolean(idx int) bool {
	val := lState.stack.get(idx)
	return convertToBoolean(val)
}

// ToNumber 取出指定索引上的数字，如果索引上值不是数字，那么将指定索引上的值转换为数字。如果无法转换那么返回0
func (lState *luaState) ToNumber(idx int) float64 {
	n, _ := lState.ToNumberX(idx)
	return n
}

// ToNumberX 取出索引上的数字，如果不是数字则尝试转换，转换成功会返回转换结果和true，转换失败会返回0和false
func (lState *luaState) ToNumberX(idx int) (float64, bool) {
	val := lState.stack.get(idx)
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

// ToInteger 取出索引处值并转换为int64
func (lState *luaState) ToInteger(idx int) int64 {
	i, _ := lState.ToIntegerX(idx)
	return i
}

// ToIntegerX 取出索引处值，并转换成int64，会返回象征转换失败与否的bool
func (lState *luaState) ToIntegerX(idx int) (int64, bool) {
	val := lState.stack.get(idx)
	i, ok := val.(int64)
	return i, ok
}

// ToString 获取指定索引上的string，如果不是string会试图转换成string，会修改栈上内容
func (lState *luaState) ToString(idx int) string {
	s, _ := lState.ToStringX(idx)
	return s
}

// ToStringX 如果是string类型那么会直接返回，如果是数字那么会转换成字符串然后返回，如果不能转换那么会返回空字符串，并返回false代表转换失败
// 注意会修改栈上的内容
func (lState *luaState) ToStringX(idx int) (string, bool) {
	val := lState.stack.get(idx)
	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		// 这里会修改栈上的内容
		lState.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}
