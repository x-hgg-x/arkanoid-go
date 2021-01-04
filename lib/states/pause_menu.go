package states

import (
	"fmt"

	"github.com/x-hgg-x/arkanoid-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// PauseMenuState is the pause menu state
type PauseMenuState struct {
	pauseMenu []ecs.Entity
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

func (st *PauseMenuState) confirmSelection() states.Transition {
	switch st.selection {
	case 0:
		// Resume
		return states.Transition{Type: states.TransPop}
	case 1:
		// Main Menu
		return states.Transition{Type: states.TransReplace, NewStates: []states.State{&MainMenuState{}}}
	case 2:
		// Exit
		return states.Transition{Type: states.TransQuit}
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

// OnPause method
func (st *PauseMenuState) OnPause(world w.World) {}

// OnResume method
func (st *PauseMenuState) OnResume(world w.World) {}

// OnStart method
func (st *PauseMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.pauseMenu = append(st.pauseMenu, loader.AddEntities(world, prefabs.Menu.PauseMenu)...)
}

// OnStop method
func (st *PauseMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.pauseMenu...)
}

// Update method
func (st *PauseMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPop}
	}
	return updateMenu(st, world)
}
