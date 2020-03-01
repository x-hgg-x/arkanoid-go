package spritesystem

import (
	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
)

// TransformSystem updates geometry matrix.
// Geometry matrix is first rotated, then scaled, and finally translated.
func TransformSystem(world ecs.World) {
	for _, result := range world.Views.SpriteView.Get() {
		sprite := result.Components[world.Components.SpriteRender].(*c.SpriteRender)
		transform := result.Components[world.Components.Transform].(*c.Transform)

		spriteWidth := float64(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Width)
		spriteHeight := float64(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Height)

		// Reset geometry matrix
		sprite.Options.GeoM.Reset()

		// Perform rotation
		sprite.Options.GeoM.Translate(-spriteWidth/2, -spriteHeight/2)
		sprite.Options.GeoM.Rotate(-transform.Rotation)

		// Perform scale
		sprite.Options.GeoM.Scale(transform.Scale1.X+1, transform.Scale1.Y+1)

		// Perform translation
		screenHeight := float64(world.Resources.ScreenDimensions.Height)
		sprite.Options.GeoM.Translate(transform.Translation.X, screenHeight-transform.Translation.Y)
	}
}