package states

import (
	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	"arkanoid/lib/math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type menu interface {
	getSelection() int
	setSelection(selection int)
	confirmSelection() transition
	getMenuIDs() []string
	getCursorMenuIDs() []string
}

func updateMenu(menu menu, world ecs.World) transition {
	selection := menu.getSelection()
	numItems := len(menu.getCursorMenuIDs())

	// Handle keyboard events
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyDown):
		menu.setSelection(math.Mod(selection+1, numItems))
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		menu.setSelection(math.Mod(selection-1, numItems))
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace):
		return menu.confirmSelection()
	}

	// Handle mouse events
	for iElem, id := range menu.getMenuIDs() {
		for _, result := range world.Views.MouseReactive.Get() {
			mouseReactive := result.Components[world.Components.MouseReactive].(*c.MouseReactive)
			if mouseReactive.ID == id && mouseReactive.Hovered {
				menu.setSelection(iElem)
				if mouseReactive.JustClicked {
					return menu.confirmSelection()
				}
			}

		}
	}

	// Set cursor color
	newSelection := menu.getSelection()
	for iCursor, id := range menu.getCursorMenuIDs() {
		for _, result := range world.Views.TextView.Get() {
			text := result.Components[world.Components.Text].(*c.Text)
			if text.ID == id {
				text.Color.A = 0
				if iCursor == newSelection {
					text.Color.A = 255
				}
			}
		}
	}
	return transition{}
}
