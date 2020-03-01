package states

import (
	"fmt"

	e "arkanoid/lib/ecs"
	"arkanoid/lib/loader"
	s "arkanoid/lib/systems/sprite"
	u "arkanoid/lib/systems/ui"

	"github.com/ByteArena/ecs"
	"github.com/hajimehoshi/ebiten"
)

// LevelCompleteState is the level complete menu state
type LevelCompleteState struct {
	levelCompleteMenu []*ecs.Entity
	selection         int
}

//
// Menu interface
//
func (st *LevelCompleteState) getSelection() int {
	return st.selection
}

func (st *LevelCompleteState) setSelection(selection int) {
	st.selection = selection
}

func (st *LevelCompleteState) confirmSelection() transition {
	switch st.selection {
	case 0:
		// Main Menu
		return transition{transType: transSwitch, newStates: []state{&MainMenuState{}}}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *LevelCompleteState) getMenuIDs() []string {
	return []string{"main_menu"}
}

func (st *LevelCompleteState) getCursorMenuIDs() []string {
	return []string{"cursor_main_menu"}
}

//
// State interface
//
func (st *LevelCompleteState) onPause(world e.World)  {}
func (st *LevelCompleteState) onResume(world e.World) {}

func (st *LevelCompleteState) onStart(world e.World) {
	st.levelCompleteMenu = loader.LoadEntities("assets/metadata/entities/ui/level_complete_menu.toml", world)
}

func (st *LevelCompleteState) onStop(world e.World) {
	world.Manager.DisposeEntities(st.levelCompleteMenu...)
}

func (st *LevelCompleteState) update(world e.World, screen *ebiten.Image) transition {
	u.UISystem(world)
	s.TransformSystem(world)
	s.RenderSpriteSystem(world, screen)
	u.RenderUISystem(world, screen)

	return updateMenu(st, world)
}
