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

// Abs returns the absolute value of an integer
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Mod returns the euclidian division between 2 integers
func Mod(a, b int) int {
	m := a % b
	if m < 0 {
		m += Abs(b)
	}
	return m
}
