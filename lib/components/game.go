package components

import "arkanoid/lib/math"

// Paddle component
type Paddle struct {
	Width  float64
	Height float64
}

// Ball component
type Ball struct {
	Radius    float64
	Velocity  float64
	Direction math.Vector2
}

// StickyBall component
type StickyBall struct {
	Period float64
}

// Block component
type Block struct{}
