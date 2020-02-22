package main

import (
	e "arkanoid/ecs"
	"arkanoid/loader"
	"arkanoid/systems/sprite"
	"arkanoid/utils"

	"github.com/hajimehoshi/ebiten"
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
	ecsData := e.InitEcs()

	// Load sprite sheets
	spriteSheets := loader.LoadSpriteSheet("assets/metadata/spritesheets/spritesheets.toml")

	// Load game entities
	loader.LoadEntities("assets/metadata/entities/background.toml", ecsData, spriteSheets)
	loader.LoadEntities("assets/metadata/entities/game.toml", ecsData, spriteSheets)

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Arkanoid")

	utils.LogError(ebiten.RunGame(game{ecsData}))
}
