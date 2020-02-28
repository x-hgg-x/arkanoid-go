package uisystem

import (
	c "arkanoid/lib/components"
	e "arkanoid/lib/ecs"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// UISystem sets mouse reactive components
func UISystem(ecs e.Ecs) {
	for _, result := range ecs.Views.MouseReactive.Get() {
		sprite := result.Components[ecs.Components.SpriteRender].(*c.SpriteRender)
		transform := result.Components[ecs.Components.Transform].(*c.Transform)
		mouseReactive := result.Components[ecs.Components.MouseReactive].(*c.MouseReactive)

		screenHeight := float64(ecs.Resources.ScreenDimensions.Height)

		spriteWidth := float64(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Width)
		spriteHeight := float64(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Height)

		minX := transform.Translation.X - spriteWidth/2
		maxX := transform.Translation.X + spriteWidth/2
		minY := screenHeight - transform.Translation.Y - spriteHeight/2
		maxY := screenHeight - transform.Translation.Y + spriteHeight/2

		x, y := ebiten.CursorPosition()

		mouseReactive.Hovered = minX <= float64(x) && float64(x) <= maxX && minY <= float64(y) && float64(y) <= maxY
		if mouseReactive.Hovered && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mouseReactive.JustClicked = true
		}
	}
}
