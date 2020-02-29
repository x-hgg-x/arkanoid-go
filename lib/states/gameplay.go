package states

import (
	"arkanoid/lib/ecs"
	"arkanoid/lib/loader"
	g "arkanoid/lib/systems/game"
	i "arkanoid/lib/systems/input"
	s "arkanoid/lib/systems/sprite"
	u "arkanoid/lib/systems/ui"

	"github.com/hajimehoshi/ebiten"
)

// GameplayState is the main game state
type GameplayState struct{}

func (state *GameplayState) onPause(world ecs.World)  {}
func (state *GameplayState) onResume(world ecs.World) {}
func (state *GameplayState) onStop(world ecs.World)   {}

func (state *GameplayState) onStart(world ecs.World) {
	// Load game entities
	loader.LoadEntities("assets/metadata/entities/background.toml", world)
	loader.LoadEntities("assets/metadata/entities/game.toml", world)

	// Load ui entities
	loader.LoadEntities("assets/metadata/entities/ui/score.toml", world)
	loader.LoadEntities("assets/metadata/entities/ui/life.toml", world)
}

func (state *GameplayState) update(world ecs.World, screen *ebiten.Image) transition {
	i.InputSystem(world)
	u.UISystem(world)

	g.MovePaddleSystem(world)

	s.TransformSystem(world)
	s.RenderSpriteSystem(world, screen)
	u.RenderUISystem(world, screen)

	return transition{}
}
