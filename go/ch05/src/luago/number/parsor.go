package number

import "strconv"

// ParseInteger 将字符串转换为整数，2个返回值，转换结果以及转换是否成功
func ParseInteger(str string) (int64, bool) {
	i, err := strconv.ParseInt(str, 10, 64)
	return i, err == nil
}

// ParseFloat 将字符串转换为浮点数，2个返回值，转换结果以及转换是否成功
func ParseFloat(str string) (float64, bool) {
	f, err := strconv.ParseFloat(str, 64)
	return f, err == nil
}
