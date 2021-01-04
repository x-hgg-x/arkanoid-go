package systems

import (
	gc "github.com/x-hgg-x/arkanoid-go/lib/components"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
)

// MoveBallSystem moves ball
func MoveBallSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(gameComponents.Ball, gameComponents.StickyBall.Not(), world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		ball := gameComponents.Ball.Get(entity).(*gc.Ball)
		ballTranslation := &world.Components.Engine.Transform.Get(entity).(*ec.Transform).Translation

		ballTranslation.X += ball.Velocity * ball.VelocityMult * ball.Direction.X / ebiten.DefaultTPS
		ballTranslation.Y += ball.Velocity * ball.VelocityMult * ball.Direction.Y / ebiten.DefaultTPS
	}))
}
