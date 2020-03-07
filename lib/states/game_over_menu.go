package states

import (
	"fmt"

	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"
	"arkanoid/lib/loader"
	s "arkanoid/lib/systems/sprite"
	u "arkanoid/lib/systems/ui"

	"github.com/hajimehoshi/ebiten"
)

// GameOverState is the game over menu state
type GameOverState struct {
	Score        int
	gameOverMenu []ecs.Entity
	selection    int
}

//
// Menu interface
//

func (st *GameOverState) getSelection() int {
	return st.selection
}

func (st *GameOverState) setSelection(selection int) {
	st.selection = selection
}

func (st *GameOverState) confirmSelection() transition {
	switch st.selection {
	case 0:
		// Restart
		return transition{transType: transSwitch, newStates: []state{&GameplayState{}}}
	case 1:
		// Main Menu
		return transition{transType: transSwitch, newStates: []state{&MainMenuState{}}}
	case 2:
		// Exit
		return transition{transType: transQuit}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *GameOverState) getMenuIDs() []string {
	return []string{"restart", "main_menu", "exit"}
}

func (st *GameOverState) getCursorMenuIDs() []string {
	return []string{"cursor_restart", "cursor_main_menu", "cursor_exit"}
}

//
// State interface
//

func (st *GameOverState) onPause(world w.World)  {}
func (st *GameOverState) onResume(world w.World) {}

func (st *GameOverState) onStart(world w.World) {
	st.gameOverMenu = loader.LoadEntities("assets/metadata/entities/ui/game_over_menu.toml", world)

	ecs.Join(world.Components.Text, world.Components.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Text.Get(entity).(*c.Text)
		if text.ID == "score" {
			text.Text = fmt.Sprintf("SCORE: %d", st.Score)
		}
	}))
}

func (st *GameOverState) onStop(world w.World) {
	world.Manager.DeleteEntities(st.gameOverMenu...)
}

func (st *GameOverState) update(world w.World, screen *ebiten.Image) transition {
	u.UISystem(world)
	s.TransformSystem(world)
	s.RenderSpriteSystem(world, screen)
	u.RenderUISystem(world, screen)

	return updateMenu(st, world)
}
