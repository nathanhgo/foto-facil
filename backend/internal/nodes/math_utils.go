package nodes

// imax returns the larger of two ints
func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// imin returns the smaller of two ints
func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
