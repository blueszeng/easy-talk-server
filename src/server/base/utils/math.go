package utils

//
// Abs
//
func AbsInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func AbsInt32(v int32) int32 {
	if v < int32(0) {
		return -v
	}
	return v
}

//
// Clamp
//
func ClampInt(v int, min int, max int) int {
	if min > max {
		min, max = max, min
	}

	if v < min {
		return min
	}

	if v > max {
		return max
	}

	return v
}

func ClampInt32(v int32, min int32, max int32) int32 {
	if min > max {
		min, max = max, min
	}

	if v < min {
		return min
	}

	if v > max {
		return max
	}

	return v
}
