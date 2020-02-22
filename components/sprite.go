package components

import (
	"arkanoid/math"
	"arkanoid/utils"

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

// SpriteSheet structure
type SpriteSheet struct {
	// Texture image
	Texture *ebiten.Image `toml:"texture_image"`
	// List of sprites
	Sprites []Sprite
}

// UnmarshalTOML fills structure fields with TOML data
func (s *SpriteSheet) UnmarshalTOML(i interface{}) error {
	subSection := i.(map[string]interface{})

	textureImage, _, err := ebitenutil.NewImageFromFile(subSection["texture_image"].(string), ebiten.FilterNearest)
	utils.LogError(err)
	s.Texture = textureImage

	sprites := subSection["sprites"].([]interface{})
	s.Sprites = make([]Sprite, len(sprites))
	for iSprite, v := range sprites {
		sprite := v.(map[string]interface{})

		s.Sprites[iSprite].X = int(sprite["x"].(int64))
		s.Sprites[iSprite].Y = int(sprite["y"].(int64))
		s.Sprites[iSprite].Width = int(sprite["width"].(int64))
		s.Sprites[iSprite].Height = int(sprite["height"].(int64))
	}
	return nil
}

// SpriteRender component
type SpriteRender struct {
	// Reference sprite sheet name
	SpriteSheetName string `toml:"sprite_sheet_name"`
	// Reference sprite sheet
	SpriteSheet *SpriteSheet
	// Index of the sprite on the sprite sheet
	SpriteNumber int `toml:"sprite_number"`
	// Draw options
	Options *ebiten.DrawImageOptions
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
