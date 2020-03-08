package gamesystem

import (
	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"

	"github.com/hajimehoshi/ebiten"
)

// MoveBallSystem moves ball
func MoveBallSystem(world w.World) {
	ecs.Join(world.Components.Ball, world.Components.StickyBall.Not(), world.Components.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		ball := world.Components.Ball.Get(entity).(*c.Ball)
		ballTranslation := &world.Components.Transform.Get(entity).(*c.Transform).Translation

		ballTranslation.X += ball.Velocity * ball.VelocityMult * ball.Direction.X / ebiten.DefaultTPS
		ballTranslation.Y += ball.Velocity * ball.VelocityMult * ball.Direction.Y / ebiten.DefaultTPS
	}))
}
