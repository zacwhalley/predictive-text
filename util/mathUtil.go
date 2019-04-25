package util

// MaxInt returns the maximum of two integers, a and b
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// MinInt returns the minimum of two integers, a and b
func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}
