package main

import (
	gc "github.com/x-hgg-x/arkanoid-go/lib/components"
	gr "github.com/x-hgg-x/arkanoid-go/lib/resources"
	gs "github.com/x-hgg-x/arkanoid-go/lib/states"

	"github.com/x-hgg-x/goecsengine/loader"
	er "github.com/x-hgg-x/goecsengine/resources"
	es "github.com/x-hgg-x/goecsengine/states"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

const (
	windowWidth  = 720
	windowHeight = 600
)

type mainGame struct {
	world        w.World
	stateMachine es.StateMachine
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
	world := w.InitWorld(&gc.Components{}, nil)

	// Init screen dimensions
	world.Resources.ScreenDimensions = &er.ScreenDimensions{Width: windowWidth, Height: windowHeight}

	// Load controls
	axes := []string{gr.PaddleAxis}
	actions := []string{gr.ReleaseBallAction, gr.BallAttractionAction}
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

	utils.LogError(ebiten.RunGame(&mainGame{world, es.Init(&gs.MainMenuState{}, world)}))
}
