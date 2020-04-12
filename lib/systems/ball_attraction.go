package systems

import (
	"time"

	gc "github.com/x-hgg-x/arkanoid-go/lib/components"
	m "github.com/x-hgg-x/arkanoid-go/lib/math"
	"github.com/x-hgg-x/arkanoid-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"
)

const timeoutMillis = 250

var (
	ballAttractionTimeAccelerated   = time.Now()
	ballAttractionLastCollisionTime = time.Now()
)

// BallAttractionSystem attracts ball towards paddle
func BallAttractionSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)

	attractionLines := []ecs.Entity{}
	world.Manager.Join(gameComponents.AttractionLine, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		attractionLines = append(attractionLines, entity)
	}))

	// Test if a ball is accelerated
	if world.Manager.Join(gameComponents.Ball, gameComponents.StickyBall.Not(), world.Components.Engine.Transform).Visit(
		func(index int) (skip bool) {
			return gameComponents.Ball.Get(ecs.Entity(index)).(*gc.Ball).VelocityMult > 1
		}) {
		// Get last collision time
		for _, stopBallAttractionEvent := range gameResources.Events.StopBallAttractionEvents {
			if ballAttractionLastCollisionTime.Sub(ballAttractionTimeAccelerated) < 0 {
				ballAttractionLastCollisionTime = stopBallAttractionEvent.CollisionTime
			}
		}

		// Decelerate ball after timeout
		if ballAttractionLastCollisionTime.Sub(ballAttractionTimeAccelerated) >= 0 {
			if time.Now().Sub(ballAttractionLastCollisionTime.Add(timeoutMillis*time.Millisecond)) < 0 {
				return
			}

			attractionLineIndex := 0
			world.Manager.Join(gameComponents.Ball, gameComponents.StickyBall.Not(), world.Components.Engine.Transform).Visit(ecs.Visit(func(ballEntity ecs.Entity) {
				gameComponents.Ball.Get(ballEntity).(*gc.Ball).VelocityMult = 1

				if attractionLineIndex < len(attractionLines) {
					gameResources.Events.BallAttractionVfxEvents = append(gameResources.Events.BallAttractionVfxEvents, resources.BallAttractionVfxEvent{
						BallEntity:               ballEntity,
						BallColorScale:           [4]float64{1, 1, 1, 1},
						AttractionLineEntity:     attractionLines[attractionLineIndex],
						AttractionLineColorScale: [4]float64{1, 1, 1, 0},
					})
				}
				attractionLineIndex++
			}))
		}
	}
	gameResources.Events.StopBallAttractionEvents = nil

	firstPaddle := ecs.GetFirst(world.Manager.Join(gameComponents.Paddle, world.Components.Engine.Transform))
	if firstPaddle == nil {
		return
	}
	paddle := gameComponents.Paddle.Get(ecs.Entity(*firstPaddle)).(*gc.Paddle)
	paddleTranslation := world.Components.Engine.Transform.Get(ecs.Entity(*firstPaddle)).(*ec.Transform).Translation

	attractionLineIndex := 0
	world.Manager.Join(gameComponents.Ball, gameComponents.StickyBall.Not(), world.Components.Engine.Transform).Visit(ecs.Visit(func(ballEntity ecs.Entity) {
		ball := gameComponents.Ball.Get(ballEntity).(*gc.Ball)
		ballTranslation := &world.Components.Engine.Transform.Get(ballEntity).(*ec.Transform).Translation

		// Attract and accelerate ball with user action
		if world.Resources.InputHandler.Actions[resources.BallAttractionAction] {
			ballAttractionTimeAccelerated = time.Now()
			ball.VelocityMult = 3
			ball.Direction = m.Vector2{
				X: paddleTranslation.X - ballTranslation.X,
				Y: paddleTranslation.Y + paddle.Height/2 + ball.Radius - ballTranslation.Y,
			}
			ball.Direction.Normalize()

			if attractionLineIndex < len(attractionLines) {
				gameResources.Events.BallAttractionVfxEvents = append(gameResources.Events.BallAttractionVfxEvents, resources.BallAttractionVfxEvent{
					BallEntity:               ballEntity,
					BallColorScale:           [4]float64{0.9, 0.6, 0.6, 1},
					AttractionLineEntity:     attractionLines[attractionLineIndex],
					AttractionLineColorScale: [4]float64{1, 1, 1, 1},
				})
			}
		} else if ball.VelocityMult > 1 && ballAttractionLastCollisionTime.Sub(ballAttractionTimeAccelerated) < 0 {
			ball.VelocityMult = 1

			if attractionLineIndex < len(attractionLines) {
				gameResources.Events.BallAttractionVfxEvents = append(gameResources.Events.BallAttractionVfxEvents, resources.BallAttractionVfxEvent{
					BallEntity:               ballEntity,
					BallColorScale:           [4]float64{1, 1, 1, 1},
					AttractionLineEntity:     attractionLines[attractionLineIndex],
					AttractionLineColorScale: [4]float64{1, 1, 1, 0},
				})
			}
		}

		attractionLineIndex++
	}))
}
