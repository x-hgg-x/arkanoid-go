package sprite

import (
	c "arkanoid/components"
	e "arkanoid/ecs"
	"sort"

	"github.com/hajimehoshi/ebiten"
)

type spriteTransform struct {
	sprite    *c.Sprite
	transform *c.Transform
}

// RenderSystem draw images.
// Images are drawn in ascending order of depth.
// Images with higher depth are thus drawn above images with lower depth.
func RenderSystem(ecs e.Ecs, screen *ebiten.Image) {
	spriteQuery := ecs.Views.SpriteView.Get()

	// Copy query slice into a struct slice for sorting
	spritesTransforms := make([]spriteTransform, len(spriteQuery))
	for i, result := range spriteQuery {
		spritesTransforms[i] = spriteTransform{
			sprite:    result.Components[ecs.Components.Sprite].(*c.Sprite),
			transform: result.Components[ecs.Components.Transform].(*c.Transform),
		}
	}

	// Sort by increasing values of depth
	sort.Slice(spritesTransforms, func(i, j int) bool {
		return spritesTransforms[i].transform.Depth < spritesTransforms[j].transform.Depth
	})

	// Sprites with higher values of depth are drawn later so they are on top
	for _, st := range spritesTransforms {
		screen.DrawImage(st.sprite.Image, st.sprite.Options)
	}
}
