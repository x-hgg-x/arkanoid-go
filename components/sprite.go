package components

import (
	"arkanoid/math"

	"github.com/hajimehoshi/ebiten"
)

// Sprite structure
type Sprite struct {
	// Horizontal position of the sprite in the sprite sheet
	X int
	// Vertical position of the sprite in the sprite sheet
	Y int
	// Width of the sprite
	Width int
	// Height of the sprite
	Height int
}

// SpriteSheet structure
type SpriteSheet struct {
	// Texture image
	Texture *ebiten.Image
	// List of sprites
	Sprites []Sprite
}

// SpriteRender component
type SpriteRender struct {
	// Reference sprite sheet
	SpriteSheet *SpriteSheet
	// Index of the sprite on the sprite sheet
	SpriteNumber int
	// Draw options
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

// NewTransform creates a new default transform, corresponding to identity.
func NewTransform() *Transform {
	return &Transform{
		Scale:       math.Vector2{X: 1, Y: 1},
		Rotation:    0,
		Translation: math.Vector2{X: 0, Y: 0},
		Depth:       0,
	}
}

// SetScale sets transform scale.
func (t *Transform) SetScale(sx, sy float64) *Transform {
	t.Scale.X = sx
	t.Scale.Y = sy
	return t
}

// SetRotation sets transform rotation.
func (t *Transform) SetRotation(angle float64) *Transform {
	t.Rotation = angle
	return t
}

// SetTranslation sets transform translation.
func (t *Transform) SetTranslation(tx, ty float64) *Transform {
	t.Translation.X = tx
	t.Translation.Y = ty
	return t
}

// SetDepth sets transform depth.
func (t *Transform) SetDepth(depth float64) *Transform {
	t.Depth = depth
	return t
}
