package states

import (
	"fmt"

	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"
	"arkanoid/lib/loader"
	s "arkanoid/lib/systems/sprite"
	u "arkanoid/lib/systems/ui"

	"github.com/hajimehoshi/ebiten"
)

// MainMenuState is the main menu state
type MainMenuState struct {
	mainMenu  []ecs.Entity
	selection int
}

//
// Menu interface
//

func (st *MainMenuState) getSelection() int {
	return st.selection
}

func (st *MainMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *MainMenuState) confirmSelection() transition {
	switch st.selection {
	case 0:
		// New game
		return transition{transType: transSwitch, newStates: []state{&GameplayState{}}}
	case 1:
		// Exit
		return transition{transType: transQuit}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *MainMenuState) getMenuIDs() []string {
	return []string{"new_game", "exit"}
}

func (st *MainMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_new_game", "cursor_exit"}
}

//
// State interface
//

func (st *MainMenuState) onPause(world w.World)  {}
func (st *MainMenuState) onResume(world w.World) {}

func (st *MainMenuState) onStart(world w.World) {
	st.mainMenu = loader.LoadEntities("assets/metadata/entities/ui/main_menu.toml", world)
}

func (st *MainMenuState) onStop(world w.World) {
	world.Manager.DeleteEntities(st.mainMenu...)
}

func (st *MainMenuState) update(world w.World, screen *ebiten.Image) transition {
	u.UISystem(world)
	s.TransformSystem(world)
	s.RenderSpriteSystem(world, screen)
	u.RenderUISystem(world, screen)

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return transition{transType: transQuit}
	}
	return updateMenu(st, world)
}
