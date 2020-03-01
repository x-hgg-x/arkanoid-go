package states

import (
	e "arkanoid/lib/ecs"
	"arkanoid/lib/loader"
	g "arkanoid/lib/systems/game"
	i "arkanoid/lib/systems/input"
	s "arkanoid/lib/systems/sprite"
	u "arkanoid/lib/systems/ui"

	"github.com/ByteArena/ecs"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// GameplayState is the main game state
type GameplayState struct {
	game []*ecs.Entity
}

func (st *GameplayState) onPause(world e.World)  {}
func (st *GameplayState) onResume(world e.World) {}

func (st *GameplayState) onStart(world e.World) {
	// Load game entities
	st.game = append(st.game, loader.LoadEntities("assets/metadata/entities/background.toml", world)...)
	st.game = append(st.game, loader.LoadEntities("assets/metadata/entities/game.toml", world)...)

	// Load ui entities
	st.game = append(st.game, loader.LoadEntities("assets/metadata/entities/ui/score.toml", world)...)
	st.game = append(st.game, loader.LoadEntities("assets/metadata/entities/ui/life.toml", world)...)
}

func (st *GameplayState) onStop(world e.World) {
	world.Manager.DisposeEntities(st.game...)
}

func (st *GameplayState) update(world e.World, screen *ebiten.Image) transition {
	i.InputSystem(world)
	u.UISystem(world)

	g.MovePaddleSystem(world)
	g.StickyBallSystem(world)

	s.TransformSystem(world)
	s.RenderSpriteSystem(world, screen)
	u.RenderUISystem(world, screen)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return transition{transType: transPush, newStates: []state{&PauseMenuState{}}}
	}
	return transition{}
}
