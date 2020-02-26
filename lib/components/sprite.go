package components

import (
	"arkanoid/lib/math"
	"arkanoid/lib/utils"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
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

// Texture structure
type Texture struct {
	// Texture image
	Image *ebiten.Image
}

// UnmarshalText fills structure fields from text data
func (t *Texture) UnmarshalText(text []byte) error {
	textureImage, _, err := ebitenutil.NewImageFromFile(string(text), ebiten.FilterNearest)
	utils.LogError(err)
	t.Image = textureImage
	return nil
}

// SpriteSheet structure
type SpriteSheet struct {
	// Texture image
	Texture Texture `toml:"texture_image"`
	// List of sprites
	Sprites []Sprite
}

// SpriteRender component
type SpriteRender struct {
	// Reference sprite sheet
	SpriteSheet *SpriteSheet
	// Index of the sprite on the sprite sheet
	SpriteNumber int `toml:"sprite_number"`
	// Draw options
	Options ebiten.DrawImageOptions
}

// Transform component.
// The origin (0, 0) is the lower left part of screen.
// Image is first rotated, then scaled, and finally translated.
type Transform struct {
	// Scale1 vector defines image scaling. Contains scale value minus 1 so that zero value is identity.
	Scale1 math.Vector2 `toml:"scale_minus_1"`
	// Rotation angle is measured counterclockwise.
	Rotation float64
	// Translation defines the position of the image center relative to the origin.
	Translation math.Vector2
	// Depth determines the drawing order on the screen. Images with higher depth are drawn above others.
	Depth float64
}

// NewTransform creates a new default transform, corresponding to identity.
func NewTransform() *Transform {
	return &Transform{}
}

// SetScale sets transform scale.
func (t *Transform) SetScale(sx, sy float64) *Transform {
	t.Scale1.X = sx - 1
	t.Scale1.Y = sy - 1
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
