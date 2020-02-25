package gamesystem

import (
	"math"

	c "arkanoid/lib/components"
	e "arkanoid/lib/ecs"
	"arkanoid/lib/resources"

	"github.com/hajimehoshi/ebiten"
)

// MovePaddleSystem moves paddle
func MovePaddleSystem(ecs e.Ecs) {
	for _, result := range ecs.Views.PaddleView.Get() {
		paddle := result.Components[ecs.Components.Paddle].(*c.Paddle)
		transform := result.Components[ecs.Components.Transform].(*c.Transform)

		screenWidth := float64(ecs.Resources.ScreenDimensions.Width)
		paddleX := transform.Translation.X
		axisValue := ecs.Resources.InputHandler.Axes[resources.PaddleAxis]

		if ecs.Resources.Controls.Axes[resources.PaddleAxis].Type == "MouseAxis" {
			paddleX = axisValue * screenWidth
		} else {
			paddleX += axisValue * screenWidth / ebiten.DefaultTPS
		}

		minValue := paddle.Width / 2
		maxValue := float64(ecs.Resources.ScreenDimensions.Width) - paddle.Width/2
		transform.Translation.X = math.Min(math.Max(float64(paddleX), minValue), maxValue)
	}
}
