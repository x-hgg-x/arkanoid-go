package states

import (
	e "arkanoid/lib/ecs"
	"arkanoid/lib/loader"
	g "arkanoid/lib/systems/game"
	i "arkanoid/lib/systems/input"
	s "arkanoid/lib/systems/sprite"
	u "arkanoid/lib/systems/ui"

	"github.com/hajimehoshi/ebiten"
)

// GameplayState is the main game state
type GameplayState struct{}

func (state *GameplayState) onPause(ecs e.Ecs)  {}
func (state *GameplayState) onResume(ecs e.Ecs) {}
func (state *GameplayState) onStop(ecs e.Ecs)   {}

func (state *GameplayState) onStart(ecs e.Ecs) {
	// Load game entities
	loader.LoadEntities("assets/metadata/entities/background.toml", ecs)
	loader.LoadEntities("assets/metadata/entities/game.toml", ecs)

	// Load ui entities
	loader.LoadEntities("assets/metadata/entities/ui/score.toml", ecs)
	loader.LoadEntities("assets/metadata/entities/ui/life.toml", ecs)
}

func (state *GameplayState) update(ecs e.Ecs, screen *ebiten.Image) transition {
	i.InputSystem(ecs)
	u.UISystem(ecs)

	g.MovePaddleSystem(ecs)

	s.TransformSystem(ecs)
	s.RenderSpriteSystem(ecs, screen)
	u.RenderUISystem(ecs, screen)

	return transition{}
}
