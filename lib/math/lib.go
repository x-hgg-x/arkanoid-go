package math

import "math"

// Vector2 type
type Vector2 struct {
	X float64
	Y float64
}

// Dot returns the dot product between two vectors
func (v Vector2) Dot(other Vector2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Perp returns the perpendicular product between two vectors
func (v Vector2) Perp(other Vector2) float64 {
	return v.X*other.Y - v.Y*other.X
}

// Norm returns the vector norm
func (v Vector2) Norm() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize normalizes the vector
func (v *Vector2) Normalize() {
	norm := v.Norm()
	v.X /= norm
	v.Y /= norm
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
