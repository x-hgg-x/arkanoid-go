package main

import (
	_ "image/png"

	gc "github.com/x-hgg-x/arkanoid-go/lib/components"
	gloader "github.com/x-hgg-x/arkanoid-go/lib/loader"
	gr "github.com/x-hgg-x/arkanoid-go/lib/resources"
	gs "github.com/x-hgg-x/arkanoid-go/lib/states"

	"github.com/x-hgg-x/goecsengine/loader"
	er "github.com/x-hgg-x/goecsengine/resources"
	es "github.com/x-hgg-x/goecsengine/states"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
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
	return windowWidth, windowHeight
}

func (game *mainGame) Update() error {
	game.stateMachine.Update(game.world)
	return nil
}

func (game *mainGame) Draw(screen *ebiten.Image) {
	game.stateMachine.Draw(game.world, screen)
}

func main() {
	world := w.InitWorld(&gc.Components{})

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

	// Load prefabs
	world.Resources.Prefabs = &gr.Prefabs{
		Menu: gr.MenuPrefabs{
			MainMenu:          gloader.PreloadEntities("assets/metadata/entities/ui/main_menu.toml", world),
			PauseMenu:         gloader.PreloadEntities("assets/metadata/entities/ui/pause_menu.toml", world),
			GameOverMenu:      gloader.PreloadEntities("assets/metadata/entities/ui/game_over_menu.toml", world),
			LevelCompleteMenu: gloader.PreloadEntities("assets/metadata/entities/ui/level_complete_menu.toml", world),
		},
		Game: gr.GamePrefabs{
			Background: gloader.PreloadEntities("assets/metadata/entities/background.toml", world),
			Game:       gloader.PreloadEntities("assets/metadata/entities/game.toml", world),
			Score:      gloader.PreloadEntities("assets/metadata/entities/ui/score.toml", world),
			Life:       gloader.PreloadEntities("assets/metadata/entities/ui/life.toml", world),
		},
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Arkanoid")

	utils.LogError(ebiten.RunGame(&mainGame{world, es.Init(&gs.MainMenuState{}, world)}))
}
