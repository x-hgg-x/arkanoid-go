package gamesystem

import (
	"fmt"

	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"
	"arkanoid/lib/resources"
)

// LifeSystem manages lifes
func LifeSystem(world w.World) {
	for range world.Resources.Game.Events.LifeEvents {
		world.Resources.Game.Lifes--

		ecs.Join(world.Components.Text, world.Components.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
			text := world.Components.Text.Get(entity).(*c.Text)
			if text.ID == "life" {
				text.Text = fmt.Sprintf("LIFES: %d", world.Resources.Game.Lifes)
			}
		}))

		if world.Resources.Game.Lifes <= 0 {
			world.Resources.Game.StateEvent = resources.StateEventGameOver
		}
	}
	world.Resources.Game.Events.LifeEvents = nil
}
