package systems

import (
	"math"

	gc "github.com/x-hgg-x/arkanoid-go/lib/components"
	"github.com/x-hgg-x/arkanoid-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	er "github.com/x-hgg-x/goecsengine/resources"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
)

// MovePaddleSystem moves paddle
func MovePaddleSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(gameComponents.Paddle, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		paddle := gameComponents.Paddle.Get(entity).(*gc.Paddle)
		paddleTransform := world.Components.Engine.Transform.Get(entity).(*ec.Transform)

		screenWidth := float64(world.Resources.ScreenDimensions.Width)
		paddleX := paddleTransform.Translation.X
		axisValue := world.Resources.InputHandler.Axes[resources.PaddleAxis]

		if _, ok := world.Resources.Controls.Axes[resources.PaddleAxis].Value.(*er.MouseAxis); ok {
			paddleX = (axisValue + 1) / 2 * screenWidth
		} else {
			paddleX += axisValue * screenWidth / ebiten.DefaultTPS
		}

		minValue := paddle.Width / 2
		maxValue := float64(world.Resources.ScreenDimensions.Width) - paddle.Width/2
		paddleTransform.Translation.X = math.Min(math.Max(paddleX, minValue), maxValue)
	}))
}
