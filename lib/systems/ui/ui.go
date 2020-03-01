package uisystem

import (
	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// UISystem sets mouse reactive components
func UISystem(world ecs.World) {
	for _, result := range world.Views.MouseReactive.Get() {
		sprite := result.Components[world.Components.SpriteRender].(*c.SpriteRender)
		transform := result.Components[world.Components.Transform].(*c.Transform)
		mouseReactive := result.Components[world.Components.MouseReactive].(*c.MouseReactive)

		screenHeight := float64(world.Resources.ScreenDimensions.Height)

		spriteWidth := float64(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Width)
		spriteHeight := float64(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Height)

		minX := transform.Translation.X - spriteWidth/2
		maxX := transform.Translation.X + spriteWidth/2
		minY := screenHeight - transform.Translation.Y - spriteHeight/2
		maxY := screenHeight - transform.Translation.Y + spriteHeight/2

		x, y := ebiten.CursorPosition()

		mouseReactive.Hovered = minX <= float64(x) && float64(x) <= maxX && minY <= float64(y) && float64(y) <= maxY
		mouseReactive.JustClicked = mouseReactive.Hovered && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	}
}
