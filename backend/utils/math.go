package utils

func Min(x, y uint32) uint32 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y uint32) uint32 {
	if x > y {
		return x
	}
	return y
}
