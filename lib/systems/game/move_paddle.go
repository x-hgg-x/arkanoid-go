package gamesystem

import (
	"math"

	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	"arkanoid/lib/resources"

	"github.com/hajimehoshi/ebiten"
)

// MovePaddleSystem moves paddle
func MovePaddleSystem(world ecs.World) {
	for _, result := range world.Views.PaddleView.Get() {
		paddle := result.Components[world.Components.Paddle].(*c.Paddle)
		transform := result.Components[world.Components.Transform].(*c.Transform)

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
		transform.Translation.X = math.Min(math.Max(float64(paddleX), minValue), maxValue)
	}
}
