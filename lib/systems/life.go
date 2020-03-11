package systems

import (
	"fmt"

	"github.com/x-hgg-x/arkanoid-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"
)

// LifeSystem manages lifes
func LifeSystem(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)

	for range gameResources.Events.LifeEvents {
		gameResources.Lifes--

		world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
			text := world.Components.Engine.Text.Get(entity).(*ec.Text)
			if text.ID == "life" {
				text.Text = fmt.Sprintf("LIFES: %d", gameResources.Lifes)
			}
		}))

		if gameResources.Lifes <= 0 {
			gameResources.StateEvent = resources.StateEventGameOver
		}
	}
	gameResources.Events.LifeEvents = nil
}
