package states

import (
	"fmt"

	"github.com/x-hgg-x/arkanoid-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"
)

// LevelCompleteState is the level complete menu state
type LevelCompleteState struct {
	Score             int
	levelCompleteMenu []ecs.Entity
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

func (st *LevelCompleteState) confirmSelection() states.Transition {
	switch st.selection {
	case 0:
		// Main Menu
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}
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

// OnPause method
func (st *LevelCompleteState) OnPause(world w.World) {}

// OnResume method
func (st *LevelCompleteState) OnResume(world w.World) {}

// OnStart method
func (st *LevelCompleteState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.levelCompleteMenu = append(st.levelCompleteMenu, loader.AddEntities(world, prefabs.Menu.LevelCompleteMenu)...)

	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		if text.ID == "score" {
			text.Text = fmt.Sprintf("SCORE: %d", st.Score)
		}
	}))
}

// OnStop method
func (st *LevelCompleteState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.levelCompleteMenu...)
}

// Update method
func (st *LevelCompleteState) Update(world w.World) states.Transition {
	return updateMenu(st, world)
}
