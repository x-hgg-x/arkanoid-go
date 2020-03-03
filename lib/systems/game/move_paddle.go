package gamesystem

import (
	"math"

	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"
	"arkanoid/lib/resources"

	"github.com/hajimehoshi/ebiten"
)

// MovePaddleSystem moves paddle
func MovePaddleSystem(world w.World) {
	ecs.Join(world.Components.Paddle, world.Components.Transform).Visit(ecs.Visit(func(index int) {
		paddle := world.Components.Paddle.Get(index).(*c.Paddle)
		transform := world.Components.Transform.Get(index).(*c.Transform)

		screenWidth := float64(world.Resources.ScreenDimensions.Width)
		paddleX := transform.Translation.X
		axisValue := world.Resources.InputHandler.Axes[resources.PaddleAxis]

		if world.Resources.Controls.Axes[resources.PaddleAxis].Type == "MouseAxis" {
			paddleX = axisValue * screenWidth
		} else {
			paddleX += axisValue * screenWidth / ebiten.DefaultTPS
		}

		minValue := paddle.Width / 2
		maxValue := float64(world.Resources.ScreenDimensions.Width) - paddle.Width/2
		transform.Translation.X = math.Min(math.Max(paddleX, minValue), maxValue)
	}))
}
