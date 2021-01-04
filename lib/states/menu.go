package states

import (
	"github.com/x-hgg-x/arkanoid-go/lib/math"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	m "github.com/x-hgg-x/goecsengine/math"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type menu interface {
	getSelection() int
	setSelection(selection int)
	confirmSelection() states.Transition
	getMenuIDs() []string
	getCursorMenuIDs() []string
}

var menuLastCursorPosition = m.VectorInt2{}

func updateMenu(menu menu, world w.World) states.Transition {
	var transition states.Transition
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

	// Handle mouse events only if mouse is moved or clicked
	x, y := ebiten.CursorPosition()
	if x != menuLastCursorPosition.X || y != menuLastCursorPosition.Y || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		menuLastCursorPosition = m.VectorInt2{X: x, Y: y}

		for iElem, id := range menu.getMenuIDs() {
			if world.Manager.Join(world.Components.Engine.SpriteRender, world.Components.Engine.Transform, world.Components.Engine.MouseReactive).Visit(
				func(index int) (skip bool) {
					mouseReactive := world.Components.Engine.MouseReactive.Get(ecs.Entity(index)).(*ec.MouseReactive)
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
	}

	// Set cursor color
	newSelection := menu.getSelection()
	for iCursor, id := range menu.getCursorMenuIDs() {
		world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
			text := world.Components.Engine.Text.Get(entity).(*ec.Text)
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
