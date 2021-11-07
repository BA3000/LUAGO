package state

import "math"
import . "luago/api"
import "luago/number"

type operator struct {
	integerFunc	func(int64, int64) int64
	floatFunc	func(float64, float64) float64
}
// 基本运算符
var (
	iadd  = func(a, b int64) int64 { return a + b }
	fadd  = func(a, b float64) float64 { return a + b }
	isub  = func(a, b int64) int64 { return a - b }
	fsub  = func(a, b float64) float64 { return a - b }
	imul  = func(a, b int64) int64 { return a * b }
	fmul  = func(a, b float64) float64 { return a * b }
	imod  = number.IMod
	fmod  = number.FMod
	pow   = math.Pow
	div   = func(a, b float64) float64 { return a / b }
	iidiv = number.IFloorDiv
	fidiv = number.FFloorDiv
	band  = func(a, b int64) int64 { return a & b }
	bor   = func(a, b int64) int64 { return a | b }
	bxor  = func(a, b int64) int64 { return a ^ b }
	shl   = number.ShiftLeft
	shr   = number.ShiftRight
	iunm  = func(a, _ int64) int64 { return -a }
	funm  = func(a, _ float64) float64 { return -a }
	bnot  = func(a, _ int64) int64 { return ^a }
)

var operators = []operator{
	{iadd, fadd},
	{isub, fsub},
	{imul, fmul},
	{imod, fmod},
	{nil, pow},
	{nil, div},
	{iidiv, fidiv},
	{band, nil},
	{bor, nil},
	{bxor, nil},
	{shl, nil},
	{shr, nil},
	{iunm, funm},
	{bnot, nil},
}

// Arith 根据情况从栈弹出1～2个操作数，然后按照索引取出对应的operator实例，最后调用_arith()执行计算
func (lState *luaState) Arith(op ArithOp) {
	// operands
	var a, b luaValue
	b = lState.stack.pop()
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		a = lState.stack.pop()
	} else {
		a = b
	}

	operator := operators[op]
	if result := _arith(a, b, operator); result != nil {
		lState.stack.push(result)
	} else {
		panic("arithmetic error!")
	}
}

// _arith 根据传入的a、b、操作数来执行运算，如果是按位操作那么会先转换为整数然后计算，结果也是整数；如果是加减乘除、整除、取反，如果都是整数那么
// 就会执行整数运算，结果会是整数；其他情况则会尝试转换为浮点数然后进行运算，运算结果也是浮点数
func _arith(a, b luaValue, op operator) luaValue {
	if op.floatFunc == nil { // bitwise运算，要转换成整数才能执行
		if x, ok := convertToInteger(a); ok {
			if y, ok := convertToInteger(b); ok {
				return op.integerFunc(x, y)
			}
		}
	} else { // arith
		if op.integerFunc != nil { // add, sub, mul, mod, idiv, unm
			if x, ok := a.(int64); ok {
				if y, ok := b.(int64); ok {
					return op.integerFunc(x, y)
				}
			}
		}
		if x, ok := convertToFloat(a); ok {
			if y, ok := convertToFloat(b); ok {
				return op.floatFunc(x, y)
			}
		}
	}
	return nil
}