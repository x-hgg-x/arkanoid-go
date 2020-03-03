package spritesystem

import (
	"image"
	"math"
	"sort"

	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"
	m "arkanoid/lib/math"

	"github.com/hajimehoshi/ebiten"
)

type spriteTransform struct {
	sprite    *c.SpriteRender
	transform *c.Transform
}

// RenderSpriteSystem draws images.
// Images are drawn in ascending order of depth.
// Images with higher depth are thus drawn above images with lower depth.
func RenderSpriteSystem(world w.World, screen *ebiten.Image) {
	sprites := ecs.Join(world.Components.SpriteRender, world.Components.Transform)

	// Copy query slice into a struct slice for sorting
	iSprite := 0
	spritesTransforms := make([]spriteTransform, sprites.Size())
	sprites.Visit(ecs.Visit(func(index int) {
		spritesTransforms[iSprite] = spriteTransform{
			sprite:    world.Components.SpriteRender.Get(index).(*c.SpriteRender),
			transform: world.Components.Transform.Get(index).(*c.Transform),
		}
		iSprite++
	}))

	// Sort by increasing values of depth
	sort.Slice(spritesTransforms, func(i, j int) bool {
		return spritesTransforms[i].transform.Depth < spritesTransforms[j].transform.Depth
	})

	// Sprites with higher values of depth are drawn later so they are on top
	for _, st := range spritesTransforms {
		drawImageWithWrap(screen, st.sprite)
	}
}

// Draw sprite with texture wrapping.
// Image is tiled when texture coordinates are greater than image size.
func drawImageWithWrap(screen *ebiten.Image, spriteRender *c.SpriteRender) {
	sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]
	texture := spriteRender.SpriteSheet.Texture
	textureWidth, textureHeight := texture.Image.Size()

	startX := int(math.Floor(float64(sprite.X) / float64(textureWidth)))
	startY := int(math.Floor(float64(sprite.Y) / float64(textureHeight)))

	stopX := int(math.Ceil(float64(sprite.X+sprite.Width) / float64(textureWidth)))
	stopY := int(math.Ceil(float64(sprite.Y+sprite.Height) / float64(textureHeight)))

	currentX := 0
	for indX := startX; indX < stopX; indX++ {
		left := m.Max(0, sprite.X-indX*textureWidth)
		right := m.Min(textureWidth, sprite.X+sprite.Width-indX*textureWidth)

		currentY := 0
		for indY := startY; indY < stopY; indY++ {
			top := m.Max(0, sprite.Y-indY*textureHeight)
			bottom := m.Min(textureHeight, sprite.Y+sprite.Height-indY*textureHeight)

			op := spriteRender.Options
			op.GeoM.Translate(float64(currentX), float64(currentY))
			screen.DrawImage(texture.Image.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image), &op)

			currentY += bottom - top
		}
		currentX += right - left
	}
}
