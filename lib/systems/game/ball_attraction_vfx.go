package gamesystem

import (
	"math"

	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"
	m "arkanoid/lib/math"
)

// BallAttractionVfxSystem add visual effets when ball is accelerated
func BallAttractionVfxSystem(world w.World) {
	for _, ballAttractionVfxEvent := range world.Resources.Game.Events.BallAttractionVfxEvents {
		ballColorScale := ballAttractionVfxEvent.BallColorScale
		ballColorM := &world.Components.SpriteRender.Get(ballAttractionVfxEvent.BallEntity).(*c.SpriteRender).Options.ColorM
		ballColorM.Reset()
		ballColorM.Scale(ballColorScale[0], ballColorScale[1], ballColorScale[2], ballColorScale[3])

		attractionLineColorScale := ballAttractionVfxEvent.AttractionLineColorScale
		world.Components.Transform.Get(ballAttractionVfxEvent.AttractionLineEntity).(*c.Transform).Depth = 0.1
		attractionLineColorM := &world.Components.SpriteRender.Get(ballAttractionVfxEvent.AttractionLineEntity).(*c.SpriteRender).Options.ColorM
		attractionLineColorM.Reset()
		attractionLineColorM.Scale(attractionLineColorScale[0], attractionLineColorScale[1], attractionLineColorScale[2], attractionLineColorScale[3])
	}
	world.Resources.Game.Events.BallAttractionVfxEvents = nil

	paddles := ecs.Join(world.Components.Paddle, world.Components.Transform)
	if paddles.Empty() {
		return
	}
	firstPaddle := ecs.Entity(paddles.Next(-1))
	paddle := world.Components.Paddle.Get(firstPaddle).(*c.Paddle)
	paddleTranslation := world.Components.Transform.Get(firstPaddle).(*c.Transform).Translation

	attractionLines := []ecs.Entity{}
	ecs.Join(world.Components.AttractionLine, world.Components.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		attractionLines = append(attractionLines, entity)
	}))

	balls := ecs.Join(world.Components.Ball, world.Components.StickyBall.Not(), world.Components.Transform)
	if len(attractionLines) != balls.Size() {
		return
	}

	attractionLineIndex := 0
	balls.Visit(ecs.Visit(func(entity ecs.Entity) {
		ball := world.Components.Ball.Get(entity).(*c.Ball)
		ballTranslation := &world.Components.Transform.Get(entity).(*c.Transform).Translation

		attractionLineTransform := world.Components.Transform.Get(attractionLines[attractionLineIndex]).(*c.Transform)

		attractionLineTransform.Translation = m.Vector2{
			X: (paddleTranslation.X + ballTranslation.X) / 2,
			Y: (paddleTranslation.Y + paddle.Height/2 + ball.Radius + ballTranslation.Y) / 2,
		}

		attractionLineVect := m.Vector2{
			X: paddleTranslation.X - ballTranslation.X,
			Y: paddleTranslation.Y + paddle.Height/2 + ball.Radius - ballTranslation.Y,
		}

		attractionLineTransform.SetRotation(math.Atan2(m.Vector2{Y: -1}.Perp(attractionLineVect), m.Vector2{Y: -1}.Dot(attractionLineVect)))
		attractionLineTransform.SetScale(1, attractionLineVect.Norm())

		attractionLineIndex++
	}))
}
