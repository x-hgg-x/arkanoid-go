package uisystem

import (
	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// UISystem sets mouse reactive components
func UISystem(world w.World) {
	ecs.Join(world.Components.SpriteRender, world.Components.Transform, world.Components.MouseReactive).Visit(ecs.Visit(func(index int) {
		sprite := world.Components.SpriteRender.Get(index).(*c.SpriteRender)
		transform := world.Components.Transform.Get(index).(*c.Transform)
		mouseReactive := world.Components.MouseReactive.Get(index).(*c.MouseReactive)

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
	}))
}
