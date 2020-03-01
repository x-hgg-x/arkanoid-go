package states

import (
	"fmt"

	e "arkanoid/lib/ecs"
	"arkanoid/lib/loader"
	s "arkanoid/lib/systems/sprite"
	u "arkanoid/lib/systems/ui"

	"github.com/ByteArena/ecs"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// PauseMenuState is the pause menu state
type PauseMenuState struct {
	pauseMenu []*ecs.Entity
	selection int
}

//
// Menu interface
//
func (st *PauseMenuState) getSelection() int {
	return st.selection
}

func (st *PauseMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *PauseMenuState) confirmSelection() transition {
	switch st.selection {
	case 0:
		// Resume
		return transition{transType: transPop}
	case 1:
		// Main Menu
		return transition{transType: transReplace, newStates: []state{&MainMenuState{}}}
	case 2:
		// Exit
		return transition{transType: transQuit}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *PauseMenuState) getMenuIDs() []string {
	return []string{"resume", "main_menu", "exit"}
}

func (st *PauseMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_resume", "cursor_main_menu", "cursor_exit"}
}

//
// State interface
//
func (st *PauseMenuState) onPause(world e.World)  {}
func (st *PauseMenuState) onResume(world e.World) {}

func (st *PauseMenuState) onStart(world e.World) {
	st.pauseMenu = loader.LoadEntities("assets/metadata/entities/ui/pause_menu.toml", world)
}

func (st *PauseMenuState) onStop(world e.World) {
	world.Manager.DisposeEntities(st.pauseMenu...)
}

func (st *PauseMenuState) update(world e.World, screen *ebiten.Image) transition {
	u.UISystem(world)
	s.TransformSystem(world)
	s.RenderSpriteSystem(world, screen)
	u.RenderUISystem(world, screen)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return transition{transType: transPop}
	}
	return updateMenu(st, world)
}
