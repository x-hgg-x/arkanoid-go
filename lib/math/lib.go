package math

// Vector2 type
type Vector2 struct {
	X float64
	Y float64
}

// VectorInt2 type
type VectorInt2 struct {
	X int
	Y int
}

// Min returns the minimum between 2 integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum between 2 integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
