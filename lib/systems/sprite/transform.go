package spritesystem

import (
	c "arkanoid/lib/components"
	e "arkanoid/lib/ecs"
)

// TransformSystem updates geometry matrix.
// Geometry matrix is first rotated, then scaled, and finally translated.
func TransformSystem(ecs e.Ecs) {
	for _, result := range ecs.Views.SpriteView.Get() {
		sprite := result.Components[ecs.Components.SpriteRender].(*c.SpriteRender)
		transform := result.Components[ecs.Components.Transform].(*c.Transform)

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
		screenHeight := float64(ecs.Resources.ScreenDimensions.Height)
		sprite.Options.GeoM.Translate(transform.Translation.X, screenHeight-transform.Translation.Y)
	}
}
