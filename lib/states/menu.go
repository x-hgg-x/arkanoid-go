package states

import (
	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"
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

func updateMenu(menu menu, world w.World) transition {
	var transition transition
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
		if ecs.Join(world.Components.SpriteRender, world.Components.Transform, world.Components.MouseReactive).Visit(
			func(index int) (skip bool) {
				mouseReactive := world.Components.MouseReactive.Get(ecs.Entity(index)).(*c.MouseReactive)
				if mouseReactive.ID == id && mouseReactive.Hovered {
					menu.setSelection(iElem)
					if mouseReactive.JustClicked {
						transition = menu.confirmSelection()
						return true
					}
				}
				return false
			}) {
			return transition
		}
	}

	// Set cursor color
	newSelection := menu.getSelection()
	for iCursor, id := range menu.getCursorMenuIDs() {
		ecs.Join(world.Components.Text, world.Components.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
			text := world.Components.Text.Get(entity).(*c.Text)
			if text.ID == id {
				text.Color.A = 0
				if iCursor == newSelection {
					text.Color.A = 255
				}
			}
		}))
	}
	return transition
}
