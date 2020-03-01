package gamesystem

import (
	"math"

	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	m "arkanoid/lib/math"
	"arkanoid/lib/resources"

	"github.com/hajimehoshi/ebiten"
)

var stickyBallFrame = 0

// StickyBallSystem moves sticky ball with paddle
func StickyBallSystem(world ecs.World) {
	stickyBallFrame++

	if len(world.Views.PaddleView.Get()) < 1 {
		return
	}
	paddleComponents := world.Views.PaddleView.Get()[0].Components
	paddleWidth := paddleComponents[world.Components.Paddle].(*c.Paddle).Width
	paddleX := paddleComponents[world.Components.Transform].(*c.Transform).Translation.X

	for _, result := range world.Views.StickyBallView.Get() {
		ball := result.Components[world.Components.Ball].(*c.Ball)
		stickyBall := result.Components[world.Components.StickyBall].(*c.StickyBall)
		ballTransform := result.Components[world.Components.Transform].(*c.Transform)

		// Follow paddle
		translationMinValue := ball.Radius / 2
		translationMaxValue := float64(world.Resources.ScreenDimensions.Width) - ball.Radius/2
		ballTransform.Translation.X = math.Min(math.Max(paddleX, translationMinValue), translationMaxValue)

		// Add oscillation
		ballTransform.Translation.X += stickyBall.WidthExtent / 2 * math.Sin(2*math.Pi*float64(stickyBallFrame)/ebiten.DefaultTPS/stickyBall.Period)

		// Set ball direction
		angleMinValue := -math.Pi / 3
		angleMaxValue := math.Pi / 3
		angle := math.Min(math.Max((paddleX-ballTransform.Translation.X)/paddleWidth*math.Pi, angleMinValue), angleMaxValue)
		ball.Direction = m.Vector2{X: math.Sin(-angle), Y: math.Cos(angle)}
	}

	if world.Resources.InputHandler.Actions[resources.ReleaseBallAction] {
		for _, result := range world.Views.StickyBallView.Get() {
			result.Entity.RemoveComponent(world.Components.StickyBall)
		}
	}
}
