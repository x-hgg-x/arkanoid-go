package main

import (
	"log"

	c "arkanoid/components"
	e "arkanoid/ecs"
	"arkanoid/systems/sprite"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	windowWidth  = 720
	windowHeight = 600
)

type game struct {
	ecs e.Ecs
}

func (g game) Layout(outsideWidth, outsideHeight int) (int, int) {
	ebiten.SetWindowSize(outsideWidth, outsideHeight)
	return windowWidth, windowHeight
}

func (g game) Update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	sprite.TransformSystem(g.ecs, screen)
	sprite.RenderSystem(g.ecs, screen)

	return nil
}

func main() {
	backgroundTextureImage, _, _ := ebitenutil.NewImageFromFile("assets/textures/background.png", ebiten.FilterNearest)
	gameTextureImage, _, _ := ebitenutil.NewImageFromFile("assets/textures/spritesheet.png", ebiten.FilterNearest)

	backgroundSpriteSheet := c.SpriteSheet{
		Texture: backgroundTextureImage,
		Sprites: []c.Sprite{
			c.Sprite{
				X:      0,
				Y:      0,
				Width:  windowWidth,
				Height: windowHeight,
			},
		},
	}

	gameSpriteSheet := c.SpriteSheet{
		Texture: gameTextureImage,
		Sprites: []c.Sprite{
			c.Sprite{
				X:      0,
				Y:      96,
				Width:  144,
				Height: 24,
			},
			c.Sprite{
				X:      144,
				Y:      96,
				Width:  24,
				Height: 24,
			},
		},
	}

	ecs := e.InitEcs()

	ecs.Manager.NewEntity().
		AddComponent(ecs.Components.SpriteRender,
			&c.SpriteRender{
				SpriteSheet:  &backgroundSpriteSheet,
				SpriteNumber: 0,
				Options:      &ebiten.DrawImageOptions{},
			}).
		AddComponent(ecs.Components.Transform, c.NewTransform().SetTranslation(360, 300).SetDepth(-1))

	ecs.Manager.NewEntity().
		AddComponent(ecs.Components.SpriteRender,
			&c.SpriteRender{
				SpriteSheet:  &gameSpriteSheet,
				SpriteNumber: 0,
				Options:      &ebiten.DrawImageOptions{},
			}).
		AddComponent(ecs.Components.Transform, c.NewTransform().SetTranslation(360, 12))

	ecs.Manager.NewEntity().
		AddComponent(ecs.Components.SpriteRender,
			&c.SpriteRender{
				SpriteSheet:  &gameSpriteSheet,
				SpriteNumber: 1,
				Options:      &ebiten.DrawImageOptions{},
			}).
		AddComponent(ecs.Components.Transform, c.NewTransform().SetTranslation(360, 35).SetDepth(0.2))

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Arkanoid")

	if err := ebiten.RunGame(game{ecs}); err != nil {
		log.Fatal(err)
	}
}
