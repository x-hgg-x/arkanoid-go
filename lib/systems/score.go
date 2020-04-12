package systems

import (
	"fmt"

	"github.com/x-hgg-x/arkanoid-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"
)

// ScoreSystem manages score
func ScoreSystem(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)

	for _, scoreEvent := range gameResources.Events.ScoreEvents {
		gameResources.Score += scoreEvent.Score

		world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
			text := world.Components.Engine.Text.Get(entity).(*ec.Text)
			if text.ID == "score" {
				text.Text = fmt.Sprintf("SCORE: %d", gameResources.Score)
			}
		}))
	}
	gameResources.Events.ScoreEvents = nil
}
