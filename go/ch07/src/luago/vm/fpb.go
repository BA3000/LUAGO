package vm

// Int2fb 将 int 转换为浮点字节（eeeeexxx），eeeee为0的时候字节表示的整数就是xxx，否则是(1xxx) * 2^(eeeee - 1)
func Int2fb(x int) int {
	e := 0 /* exponent */
	if x < 8 {
		return x
	}
	for x >= (8 << 4) {
		/* coarse steps */
		x = (x + 0xf) >> 4 /* x = ceil(x / 16) */
		e += 4
	}
	for x >= (8 << 1) {
		/* fine steps */
		x = (x + 1) >> 1 /* x = ceil(x / 2) */
		e++
	}
	return ((e + 1) << 3) | (x - 8)
}

// Fb2int converts back
func Fb2int(x int) int {
	if x < 8 {
		return x
	} else {
		return ((x & 7) + 8 ) << uint((x>>3)-1)
	}
}
