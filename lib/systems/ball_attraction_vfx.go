package systems

import (
	"math"

	gc "github.com/x-hgg-x/arkanoid-go/lib/components"
	gm "github.com/x-hgg-x/arkanoid-go/lib/math"
	"github.com/x-hgg-x/arkanoid-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	em "github.com/x-hgg-x/goecsengine/math"
	w "github.com/x-hgg-x/goecsengine/world"
)

// BallAttractionVfxSystem add visual effets when ball is accelerated
func BallAttractionVfxSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)

	for _, ballAttractionVfxEvent := range gameResources.Events.BallAttractionVfxEvents {
		ballColorScale := ballAttractionVfxEvent.BallColorScale
		ballColorM := &world.Components.Engine.SpriteRender.Get(ballAttractionVfxEvent.BallEntity).(*ec.SpriteRender).Options.ColorM
		ballColorM.Reset()
		ballColorM.Scale(ballColorScale[0], ballColorScale[1], ballColorScale[2], ballColorScale[3])

		attractionLineColorScale := ballAttractionVfxEvent.AttractionLineColorScale
		world.Components.Engine.Transform.Get(ballAttractionVfxEvent.AttractionLineEntity).(*ec.Transform).Depth = 0.1
		attractionLineColorM := &world.Components.Engine.SpriteRender.Get(ballAttractionVfxEvent.AttractionLineEntity).(*ec.SpriteRender).Options.ColorM
		attractionLineColorM.Reset()
		attractionLineColorM.Scale(attractionLineColorScale[0], attractionLineColorScale[1], attractionLineColorScale[2], attractionLineColorScale[3])
	}
	gameResources.Events.BallAttractionVfxEvents = nil

	firstPaddle := ecs.GetFirst(world.Manager.Join(gameComponents.Paddle, world.Components.Engine.Transform))
	if firstPaddle == nil {
		return
	}
	paddle := gameComponents.Paddle.Get(ecs.Entity(*firstPaddle)).(*gc.Paddle)
	paddleTranslation := world.Components.Engine.Transform.Get(ecs.Entity(*firstPaddle)).(*ec.Transform).Translation

	attractionLines := []ecs.Entity{}
	world.Manager.Join(gameComponents.AttractionLine, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		attractionLines = append(attractionLines, entity)
	}))

	balls := world.Manager.Join(gameComponents.Ball, gameComponents.StickyBall.Not(), world.Components.Engine.Transform)
	if len(attractionLines) != balls.Size() {
		return
	}

	attractionLineIndex := 0
	balls.Visit(ecs.Visit(func(entity ecs.Entity) {
		ball := gameComponents.Ball.Get(entity).(*gc.Ball)
		ballTranslation := &world.Components.Engine.Transform.Get(entity).(*ec.Transform).Translation

		attractionLineTransform := world.Components.Engine.Transform.Get(attractionLines[attractionLineIndex]).(*ec.Transform)

		attractionLineTransform.Translation = em.Vector2{
			X: (paddleTranslation.X + ballTranslation.X) / 2,
			Y: (paddleTranslation.Y + paddle.Height/2 + ball.Radius + ballTranslation.Y) / 2,
		}

		attractionLineVect := gm.Vector2{
			X: paddleTranslation.X - ballTranslation.X,
			Y: paddleTranslation.Y + paddle.Height/2 + ball.Radius - ballTranslation.Y,
		}

		attractionLineTransform.SetRotation(math.Atan2(gm.Vector2{Y: -1}.Perp(attractionLineVect), gm.Vector2{Y: -1}.Dot(attractionLineVect)))
		attractionLineTransform.SetScale(1, attractionLineVect.Norm())

		attractionLineIndex++
	}))
}
