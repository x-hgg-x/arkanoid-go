package gamesystem

import (
	"fmt"

	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"
)

// ScoreSystem manages score
func ScoreSystem(world w.World) {
	for _, scoreEvent := range world.Resources.Game.Events.ScoreEvents {
		world.Resources.Game.Score += scoreEvent.Score

		ecs.Join(world.Components.Text, world.Components.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
			text := world.Components.Text.Get(entity).(*c.Text)
			if text.ID == "score" {
				text.Text = fmt.Sprintf("SCORE: %d", world.Resources.Game.Score)
			}
		}))
	}
	world.Resources.Game.Events.ScoreEvents = nil
}
