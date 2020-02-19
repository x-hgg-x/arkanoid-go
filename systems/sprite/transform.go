package sprite

import (
	c "arkanoid/components"
	e "arkanoid/ecs"

	"github.com/hajimehoshi/ebiten"
)

// TransformSystem updates geometry matrix.
// Geometry matrix is first rotated, then scaled, and finally translated.
func TransformSystem(ecs e.Ecs, screen *ebiten.Image) {
	_, h := screen.Size()
	screenHeight := float64(h)

	for _, result := range ecs.Views.SpriteView.Get() {
		sprite := result.Components[ecs.Components.Sprite].(*c.Sprite)
		transform := result.Components[ecs.Components.Transform].(*c.Transform)

		w, h := sprite.Image.Size()
		spriteWidth, spriteHeight := float64(w), float64(h)

		// Reset geometry matrix
		sprite.Options.GeoM.Reset()

		// Perform rotation
		sprite.Options.GeoM.Translate(-spriteWidth/2, -spriteHeight/2)
		sprite.Options.GeoM.Rotate(-transform.Rotation)

		// Perform scale
		sprite.Options.GeoM.Scale(transform.Scale.X, transform.Scale.Y)

		// Perform translation
		sprite.Options.GeoM.Translate(transform.Translation.X, screenHeight-transform.Translation.Y)
	}
}
