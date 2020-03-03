package main

import (
	w "arkanoid/lib/ecs/world"
	"arkanoid/lib/loader"
	"arkanoid/lib/resources"
	"arkanoid/lib/states"
	"arkanoid/lib/utils"

	"github.com/hajimehoshi/ebiten"
)

const (
	windowWidth  = 720
	windowHeight = 600
)

type mainGame struct {
	world        w.World
	stateMachine states.StateMachine
}

func (game *mainGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	ebiten.SetWindowSize(outsideWidth, outsideHeight)
	return windowWidth, windowHeight
}

func (game *mainGame) Update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	game.stateMachine.Update(game.world, screen)
	return nil
}

func main() {
	world := w.InitWorld()

	// Init screen dimensions
	world.Resources.ScreenDimensions = &resources.ScreenDimensions{Width: windowWidth, Height: windowHeight}

	// Load controls
	axes := []string{resources.PaddleAxis}
	actions := []string{resources.ReleaseBallAction}
	controls, inputHandler := loader.LoadControls("config/controls.toml", axes, actions)
	world.Resources.Controls = &controls
	world.Resources.InputHandler = &inputHandler

	// Load sprite sheets
	spriteSheets := loader.LoadSpriteSheets("assets/metadata/spritesheets/spritesheets.toml")
	world.Resources.SpriteSheets = &spriteSheets

	// Load fonts
	fonts := loader.LoadFonts("assets/metadata/fonts/fonts.toml")
	world.Resources.Fonts = &fonts

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Arkanoid")

	utils.LogError(ebiten.RunGame(&mainGame{world, states.Init(&states.MainMenuState{}, world)}))
}
