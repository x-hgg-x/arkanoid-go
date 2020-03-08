package spritesystem

import (
	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"
)

// TransformSystem updates geometry matrix.
// Geometry matrix is first rotated, then scaled, and finally translated.
func TransformSystem(world w.World) {
	ecs.Join(world.Components.SpriteRender, world.Components.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		sprite := world.Components.SpriteRender.Get(entity).(*c.SpriteRender)
		transform := world.Components.Transform.Get(entity).(*c.Transform)

		spriteWidth := float64(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Width)
		spriteHeight := float64(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Height)

		// Reset geometry matrix
		sprite.Options.GeoM.Reset()

		// Center sprite on top left pixel
		sprite.Options.GeoM.Translate(-spriteWidth/2, -spriteHeight/2)

		// Perform scale
		sprite.Options.GeoM.Scale(transform.Scale1.X+1, transform.Scale1.Y+1)

		// Perform rotation
		sprite.Options.GeoM.Rotate(-transform.Rotation)

		// Perform translation
		screenHeight := float64(world.Resources.ScreenDimensions.Height)
		sprite.Options.GeoM.Translate(transform.Translation.X, screenHeight-transform.Translation.Y)
	}))
}
