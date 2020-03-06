package gamesystem

import (
	"math"

	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"
	m "arkanoid/lib/math"
	"arkanoid/lib/resources"

	"github.com/hajimehoshi/ebiten"
)

var stickyBallFrame = 0

// StickyBallSystem moves sticky ball with paddle
func StickyBallSystem(world w.World) {
	stickyBallFrame++

	paddles := ecs.Join(world.Components.Paddle, world.Components.Transform)
	if paddles.Empty() {
		return
	}
	firstPaddle := ecs.Entity(paddles.Next(-1))
	paddleWidth := world.Components.Paddle.Get(firstPaddle).(*c.Paddle).Width
	paddleX := world.Components.Transform.Get(firstPaddle).(*c.Transform).Translation.X

	stickyBalls := ecs.Join(world.Components.Ball, world.Components.StickyBall, world.Components.Transform)

	stickyBalls.Visit(ecs.Visit(func(entity ecs.Entity) {
		ball := world.Components.Ball.Get(entity).(*c.Ball)
		stickyBall := world.Components.StickyBall.Get(entity).(*c.StickyBall)
		ballTransform := world.Components.Transform.Get(entity).(*c.Transform)

		// Follow paddle
		translationMinValue := ball.Radius / 2
		translationMaxValue := float64(world.Resources.ScreenDimensions.Width) - ball.Radius/2
		ballTransform.Translation.X = math.Min(math.Max(paddleX, translationMinValue), translationMaxValue)

		// Add oscillation
		ballTransform.Translation.X += paddleWidth / 4 * math.Sin(2*math.Pi*float64(stickyBallFrame)/ebiten.DefaultTPS/stickyBall.Period)

		// Set ball direction
		angle := (paddleX - ballTransform.Translation.X) / paddleWidth * math.Pi
		ball.Direction = m.Vector2{X: math.Sin(-angle), Y: math.Cos(angle)}
	}))

	if world.Resources.InputHandler.Actions[resources.ReleaseBallAction] {
		stickyBalls.Visit(ecs.Visit(func(entity ecs.Entity) {
			ecs.Entity(entity).RemoveComponent(world.Components.StickyBall)
		}))
	}
}
