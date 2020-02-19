package components

import (
	"arkanoid/math"

	"github.com/hajimehoshi/ebiten"
)

// Sprite component
type Sprite struct {
	Image   *ebiten.Image
	Options *ebiten.DrawImageOptions
}

// Transform component.
// The origin (0, 0) is the lower left part of screen.
// Image is first rotated, then scaled, and finally translated.
type Transform struct {
	// Scale vector defines image scaling. Identity is (1, 1).
	Scale math.Vector2
	// Rotation angle is measured counterclockwise.
	Rotation float64
	// Translation defines the position of the image center relative to the origin.
	Translation math.Vector2
	// Depth determines the drawing order on the screen. Images with higher depth are drawn above others.
	Depth float64
}
