package gamesystem

import (
	c "arkanoid/lib/components"
	w "arkanoid/lib/ecs/world"
)

// BlockHealthSystem manages block health
func BlockHealthSystem(world w.World) {
	for _, blockCollisionEvent := range world.Resources.Game.Events.BlockCollisionEvents {
		block := world.Components.Block.Get(blockCollisionEvent.Entity).(*c.Block)
		sprite := world.Components.SpriteRender.Get(blockCollisionEvent.Entity).(*c.SpriteRender)

		block.Health--
		if block.Health > 0 {
			sprite.SpriteNumber += 6
		} else {
			world.Resources.Game.CollisionWorld.DestroyBody(block.Body)
			world.Manager.DeleteEntity(blockCollisionEvent.Entity)
		}
	}
	world.Resources.Game.Events.BlockCollisionEvents = nil
}
