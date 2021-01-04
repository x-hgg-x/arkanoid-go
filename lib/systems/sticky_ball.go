package systems

import (
	"math"

	gc "github.com/x-hgg-x/arkanoid-go/lib/components"
	m "github.com/x-hgg-x/arkanoid-go/lib/math"
	"github.com/x-hgg-x/arkanoid-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
)

var stickyBallFrame = 0

// StickyBallSystem moves sticky ball with paddle
func StickyBallSystem(world w.World) {
	stickyBallFrame++

	gameComponents := world.Components.Game.(*gc.Components)

	firstPaddle := ecs.GetFirst(world.Manager.Join(gameComponents.Paddle, world.Components.Engine.Transform))
	if firstPaddle == nil {
		return
	}
	paddleWidth := gameComponents.Paddle.Get(ecs.Entity(*firstPaddle)).(*gc.Paddle).Width
	paddleX := world.Components.Engine.Transform.Get(ecs.Entity(*firstPaddle)).(*ec.Transform).Translation.X

	stickyBalls := world.Manager.Join(gameComponents.Ball, gameComponents.StickyBall, world.Components.Engine.Transform)

	stickyBalls.Visit(ecs.Visit(func(entity ecs.Entity) {
		ball := gameComponents.Ball.Get(entity).(*gc.Ball)
		stickyBall := gameComponents.StickyBall.Get(entity).(*gc.StickyBall)
		ballTranslation := &world.Components.Engine.Transform.Get(entity).(*ec.Transform).Translation

		// Follow paddle
		translationMinValue := ball.Radius / 2
		translationMaxValue := float64(world.Resources.ScreenDimensions.Width) - ball.Radius/2
		ballTranslation.X = math.Min(math.Max(paddleX, translationMinValue), translationMaxValue)

		// Add oscillation
		ballTranslation.X += paddleWidth / 4 * math.Sin(2*math.Pi*float64(stickyBallFrame)/ebiten.DefaultTPS/stickyBall.Period)

		// Set ball direction
		angle := (paddleX - ballTranslation.X) / paddleWidth * math.Pi
		ball.Direction = m.Vector2{X: math.Sin(-angle), Y: math.Cos(angle)}

		// Reset ball velocity
		ball.VelocityMult = 1
	}))

	if world.Resources.InputHandler.Actions[resources.ReleaseBallAction] {
		stickyBalls.Visit(ecs.Visit(func(entity ecs.Entity) {
			entity.RemoveComponent(gameComponents.StickyBall)
		}))
	}
}
