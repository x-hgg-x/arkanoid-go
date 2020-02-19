package main

import (
	"log"
	"math"

	c "arkanoid/components"
	e "arkanoid/ecs"
	m "arkanoid/math"
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
	sprite.TransformSystem(g.ecs, screen)
	sprite.RenderSystem(g.ecs, screen)

	return nil
}

func main() {
	image, _, _ := ebitenutil.NewImageFromFile("gopher.png", ebiten.FilterNearest)

	ecs := e.InitEcs()

	ecs.Manager.NewEntity().
		AddComponent(ecs.Components.Sprite,
			&c.Sprite{
				Image:   image,
				Options: &ebiten.DrawImageOptions{},
			}).
		AddComponent(ecs.Components.Transform,
			&c.Transform{
				Scale:       m.Vector2{X: 1, Y: 1},
				Rotation:    math.Pi / 4,
				Translation: m.Vector2{X: 300, Y: 300},
				Depth:       1,
			})

	ecs.Manager.NewEntity().
		AddComponent(ecs.Components.Sprite,
			&c.Sprite{
				Image:   image,
				Options: &ebiten.DrawImageOptions{},
			}).
		AddComponent(ecs.Components.Transform,
			&c.Transform{
				Scale:       m.Vector2{X: 5, Y: 5},
				Rotation:    math.Pi / 2,
				Translation: m.Vector2{X: 360, Y: 300},
				Depth:       0,
			})

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Arkanoid")

	if err := ebiten.RunGame(game{ecs}); err != nil {
		log.Fatal(err)
	}
}
