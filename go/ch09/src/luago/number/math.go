package number

import "math"

// IFloorDiv int64整数整除，会向下取整
func IFloorDiv(a, b int64) int64 {
	if a > 0 && b > 0 || a < 0 && b < 0 || a%b ==0 {
		return a / b
	} else {
		return a/b - 1
	}
}

// FFloorDiv 浮点数整除
func FFloorDiv(a, b float64) float64 {
	return math.Floor(a / b)
}

// IMod 整数取模
func IMod(a, b int64) int64 {
	return a - IFloorDiv(a, b)*b
}

// FMod 浮点数取模
func FMod(a, b float64) float64 {
	return a - math.Floor(a/b)*b
}

// ShiftLeft 左移函数，如果输入位移数量为负，那么会自动转换为右移
func ShiftLeft(a, n int64) int64 {
	if n >= 0 {
		return a << uint64(n)
	} else {
		return ShiftRight(a, n)
	}
}

// ShiftRight 右移，如果输入数量为负，那么会自动转换为左移
func ShiftRight(a, n int64) int64 {
	if n >= 0 {
		// Go里面右移只能够是无符号数，所以这里进行转换
		return int64(uint64(a) >> uint64(n))
	} else {
		return ShiftLeft(a, n)
	}
}

// FloatToInteger 将浮点数转换为整数
func FloatToInteger(f float64) (int64, bool) {
	i := int64(f)
	// 检测浮点数有没有超出整数表示范围
	return i, float64(i) == f
}
