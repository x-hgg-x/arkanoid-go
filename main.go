package main

import (
	"os"

	e "arkanoid/lib/ecs"
	"arkanoid/lib/loader"
	"arkanoid/lib/resources"
	g "arkanoid/lib/systems/game"
	i "arkanoid/lib/systems/input"
	s "arkanoid/lib/systems/sprite"
	u "arkanoid/lib/systems/ui"
	"arkanoid/lib/utils"

	"github.com/hajimehoshi/ebiten"
)

const (
	windowWidth  = 720
	windowHeight = 600
)

type mainGame struct {
	ecs e.Ecs
}

func (game mainGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	ebiten.SetWindowSize(outsideWidth, outsideHeight)
	return windowWidth, windowHeight
}

func (game mainGame) Update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	i.InputSystem(game.ecs)
	u.UISystem(game.ecs)

	g.MovePaddleSystem(game.ecs)

	s.TransformSystem(game.ecs)
	s.RenderSpriteSystem(game.ecs, screen)
	u.RenderUISystem(game.ecs, screen)

	return nil
}

func main() {
	ecsData := e.InitEcs()

	// Init screen dimensions
	ecsData.Resources.ScreenDimensions = &resources.ScreenDimensions{Width: windowWidth, Height: windowHeight}

	// Load controls
	axes := []string{resources.PaddleAxis}
	actions := []string{resources.ReleaseBallAction}
	controls, inputHandler := loader.LoadControls("config/controls.toml", axes, actions)
	ecsData.Resources.Controls = &controls
	ecsData.Resources.InputHandler = &inputHandler

	// Load sprite sheets
	spriteSheets := loader.LoadSpriteSheets("assets/metadata/spritesheets/spritesheets.toml")
	ecsData.Resources.SpriteSheets = &spriteSheets

	// Load fonts
	fonts := loader.LoadFonts("assets/metadata/fonts/fonts.toml")
	ecsData.Resources.Fonts = &fonts

	// Load game entities
	loader.LoadEntities("assets/metadata/entities/background.toml", ecsData)
	loader.LoadEntities("assets/metadata/entities/game.toml", ecsData)

	// Load score and life entities
	loader.LoadEntities("assets/metadata/entities/ui/score.toml", ecsData)
	loader.LoadEntities("assets/metadata/entities/ui/life.toml", ecsData)

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Arkanoid")

	utils.LogError(ebiten.RunGame(mainGame{ecsData}))
}
