package gamesystem

import (
	"time"

	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"
	m "arkanoid/lib/math"
	"arkanoid/lib/resources"
)

const timeoutMillis = 250

var (
	ballAttractionTimeAccelerated   = time.Now()
	ballAttractionLastCollisionTime = time.Now()
)

// BallAttractionSystem attracts ball towards paddle
func BallAttractionSystem(world w.World) {
	attractionLines := []ecs.Entity{}
	ecs.Join(world.Components.AttractionLine, world.Components.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		attractionLines = append(attractionLines, entity)
	}))

	// Test if a ball is accelerated
	if ecs.Join(world.Components.Ball, world.Components.StickyBall.Not(), world.Components.Transform).Visit(
		func(index int) (skip bool) {
			return world.Components.Ball.Get(ecs.Entity(index)).(*c.Ball).VelocityMult > 1
		}) {
		// Get last collision time
		for _, stopBallAttractionEvent := range world.Resources.Game.Events.StopBallAttractionEvents {
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
			ecs.Join(world.Components.Ball, world.Components.StickyBall.Not(), world.Components.Transform).Visit(ecs.Visit(func(ballEntity ecs.Entity) {
				world.Components.Ball.Get(ballEntity).(*c.Ball).VelocityMult = 1

				if attractionLineIndex < len(attractionLines) {
					world.Resources.Game.Events.BallAttractionVfxEvents = append(world.Resources.Game.Events.BallAttractionVfxEvents, resources.BallAttractionVfxEvent{
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
	world.Resources.Game.Events.StopBallAttractionEvents = nil

	paddles := ecs.Join(world.Components.Paddle, world.Components.Transform)
	if paddles.Empty() {
		return
	}
	firstPaddle := ecs.Entity(paddles.Next(-1))
	paddle := world.Components.Paddle.Get(firstPaddle).(*c.Paddle)
	paddleTranslation := world.Components.Transform.Get(firstPaddle).(*c.Transform).Translation

	attractionLineIndex := 0
	ecs.Join(world.Components.Ball, world.Components.StickyBall.Not(), world.Components.Transform).Visit(ecs.Visit(func(ballEntity ecs.Entity) {
		ball := world.Components.Ball.Get(ballEntity).(*c.Ball)
		ballTranslation := &world.Components.Transform.Get(ballEntity).(*c.Transform).Translation

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
				world.Resources.Game.Events.BallAttractionVfxEvents = append(world.Resources.Game.Events.BallAttractionVfxEvents, resources.BallAttractionVfxEvent{
					BallEntity:               ballEntity,
					BallColorScale:           [4]float64{0.9, 0.6, 0.6, 1},
					AttractionLineEntity:     attractionLines[attractionLineIndex],
					AttractionLineColorScale: [4]float64{1, 1, 1, 1},
				})
			}
		} else if ball.VelocityMult > 1 && ballAttractionLastCollisionTime.Sub(ballAttractionTimeAccelerated) < 0 {
			ball.VelocityMult = 1

			if attractionLineIndex < len(attractionLines) {
				world.Resources.Game.Events.BallAttractionVfxEvents = append(world.Resources.Game.Events.BallAttractionVfxEvents, resources.BallAttractionVfxEvent{
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
