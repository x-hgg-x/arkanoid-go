package systems

import (
	"math"
	"time"

	gc "github.com/x-hgg-x/arkanoid-go/lib/components"
	gm "github.com/x-hgg-x/arkanoid-go/lib/math"
	"github.com/x-hgg-x/arkanoid-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	em "github.com/x-hgg-x/goecsengine/math"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/ByteArena/box2d"
)

// CollisionSystem manages collisions
func CollisionSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)
	gameEvents := &gameResources.Events

	firstPaddle := ecs.GetFirst(world.Manager.Join(gameComponents.Paddle, world.Components.Engine.Transform))
	if firstPaddle == nil {
		return
	}
	paddle := gameComponents.Paddle.Get(ecs.Entity(*firstPaddle)).(*gc.Paddle)
	paddleTranslation := world.Components.Engine.Transform.Get(ecs.Entity(*firstPaddle)).(*ec.Transform).Translation

	// Set paddle body transform
	paddle.Body.SetTransform(box2d.MakeB2Vec2(paddleTranslation.X/resources.B2PixelRatio, paddleTranslation.Y/resources.B2PixelRatio), 0)

	// Set balls body transform
	world.Manager.Join(gameComponents.Ball, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		ball := gameComponents.Ball.Get(entity).(*gc.Ball)
		ballTranslation := world.Components.Engine.Transform.Get(entity).(*ec.Transform).Translation
		ball.Body.SetTransform(box2d.MakeB2Vec2(ballTranslation.X/resources.B2PixelRatio, ballTranslation.Y/resources.B2PixelRatio), 0)
	}))

	// Find contacts
	collisionWorld := gameResources.CollisionWorld
	collisionWorld.M_contactManager.FindNewContacts()
	collisionWorld.M_contactManager.Collide()

	// Get list of contacts with normals and bodies
	contactsNormal := []box2d.B2Vec2{}
	contactsBodies := [][2]*box2d.B2Body{}
	for contactList := collisionWorld.GetContactList(); contactList != nil; contactList = contactList.GetNext() {
		wm := box2d.MakeB2WorldManifold()
		contactList.GetWorldManifold(&wm)
		// Test if normal is defined
		if (wm.Normal != box2d.B2Vec2{}) {
			contactsNormal = append(contactsNormal, wm.Normal)
			contactsBodies = append(contactsBodies, [2]*box2d.B2Body{contactList.GetFixtureA().GetBody(), contactList.GetFixtureB().GetBody()})
		}
	}

	attractionLines := []ecs.Entity{}
	world.Manager.Join(gameComponents.AttractionLine, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		attractionLines = append(attractionLines, entity)
	}))

	// Loop on balls
	attractionLineIndex := 0
	world.Manager.Join(gameComponents.Ball, gameComponents.StickyBall.Not(), world.Components.Engine.Transform).Visit(ecs.Visit(func(ballEntity ecs.Entity) {
		ball := gameComponents.Ball.Get(ballEntity).(*gc.Ball)
		ballTranslation := &world.Components.Engine.Transform.Get(ballEntity).(*ec.Transform).Translation

		// Bounce at the top, left and right of the arena
		if ballTranslation.X <= ball.Radius {
			ball.Direction.X = math.Abs(ball.Direction.X)
		}
		if ballTranslation.X >= float64(world.Resources.ScreenDimensions.Width)-ball.Radius {
			ball.Direction.X = -math.Abs(ball.Direction.X)
		}
		if ballTranslation.Y >= float64(world.Resources.ScreenDimensions.Height)-ball.Radius {
			ball.Direction.Y = -math.Abs(ball.Direction.Y)
		}

		// Bounce at the paddle
		bounced := false
		for iContact := range contactsBodies {
			if contactsBodies[iContact] == [2]*box2d.B2Body{paddle.Body, ball.Body} || contactsBodies[iContact] == [2]*box2d.B2Body{ball.Body, paddle.Body} {
				bounced = true
				minValue := -math.Pi / 3
				maxValue := math.Pi / 3
				angle := math.Min(math.Max((paddleTranslation.X-ballTranslation.X)/paddle.Width*math.Pi, minValue), maxValue)
				ball.Direction = gm.Vector2{X: math.Sin(-angle), Y: math.Cos(angle)}

				gameEvents.StopBallAttractionEvents = append(gameEvents.StopBallAttractionEvents, resources.StopBallAttractionEvent{CollisionTime: time.Now()})
				break
			}
		}

		// Lose a life when ball reach the bottom of the arena
		if ballTranslation.Y <= ball.Radius && !bounced {
			ballEntity.AddComponent(gameComponents.StickyBall, &gc.StickyBall{Period: 2})
			*ballTranslation = em.Vector2{X: paddleTranslation.X, Y: paddle.Height + ball.Radius}

			gameEvents.LifeEvents = append(gameEvents.LifeEvents, resources.LifeEvent{})
			gameEvents.ScoreEvents = append(gameEvents.ScoreEvents, resources.ScoreEvent{Score: -1000})

			if attractionLineIndex < len(attractionLines) {
				gameEvents.BallAttractionVfxEvents = append(gameEvents.BallAttractionVfxEvents, resources.BallAttractionVfxEvent{
					BallEntity:               ballEntity,
					BallColorScale:           [4]float64{1, 1, 1, 1},
					AttractionLineEntity:     attractionLines[attractionLineIndex],
					AttractionLineColorScale: [4]float64{1, 1, 1, 0},
				})
			}
		}

		// Bounce at the blocks
		blockNormals := []gm.Vector2{}
		blockbodies := []*box2d.B2Body{}
		for iContact := range contactsNormal {
			// Normal is pointing towards block exterior
			var blockBody *box2d.B2Body
			if contactsBodies[iContact][0].GetUserData().(ecs.Entity).HasComponent(gameComponents.Block) && contactsBodies[iContact][1] == ball.Body {
				blockBody = contactsBodies[iContact][0]
				blockNormals = append(blockNormals, gm.Vector2{X: contactsNormal[iContact].X, Y: contactsNormal[iContact].Y})
			} else if contactsBodies[iContact][1].GetUserData().(ecs.Entity).HasComponent(gameComponents.Block) && contactsBodies[iContact][0] == ball.Body {
				blockBody = contactsBodies[iContact][1]
				blockNormals = append(blockNormals, gm.Vector2{X: -contactsNormal[iContact].X, Y: -contactsNormal[iContact].Y})
			}

			if blockBody != nil {
				blockbodies = append(blockbodies, blockBody)

				blockCollisionEvent := resources.BlockCollisionEvent{Entity: blockBody.GetUserData().(ecs.Entity)}
				gameEvents.BlockCollisionEvents = append(gameEvents.BlockCollisionEvents, blockCollisionEvent)
				gameEvents.ScoreEvents = append(gameEvents.ScoreEvents, resources.ScoreEvent{Score: 50})
			}
		}

		if len(blockNormals) == 0 {
			// No colliding blocks
			return
		}

		gameEvents.StopBallAttractionEvents = append(gameEvents.StopBallAttractionEvents, resources.StopBallAttractionEvent{CollisionTime: time.Now()})

		if len(blockNormals) >= 3 {
			// 3 or more colliding blocks: reverse ball direction
			ball.Direction.X *= -1
			ball.Direction.Y *= -1
			return
		}

		// 1 or 2 colliding blocks: ball is reflected wrt the contact normal
		var incidenceAngle float64
		if len(blockNormals) == 1 {
			// 1 colliding block: use computed normal
			incidenceAngle = math.Atan2(-ball.Direction.Perp(blockNormals[0]), -ball.Direction.Dot(blockNormals[0]))
		} else if len(blockNormals) == 2 {
			// 2 colliding blocks: define normal as the perpendicular of the line between blocks center (towards ball)
			positions := []box2d.B2Vec2{blockbodies[0].GetPosition(), blockbodies[1].GetPosition()}
			positionDiff := gm.Vector2{X: positions[1].X - positions[0].X, Y: positions[1].Y - positions[0].Y}
			positionDiffPerp := gm.Vector2{X: -positionDiff.Y, Y: positionDiff.X}
			ballLocalWorldTranslation := gm.Vector2{
				X: ballTranslation.X/resources.B2PixelRatio - positions[0].X,
				Y: ballTranslation.Y/resources.B2PixelRatio - positions[0].Y,
			}

			var normal gm.Vector2
			if positionDiffPerp.Dot(ballLocalWorldTranslation) > 0 {
				normal = gm.Vector2{X: positionDiffPerp.X, Y: positionDiffPerp.Y}
			} else {
				normal = gm.Vector2{X: -positionDiffPerp.X, Y: -positionDiffPerp.Y}
			}
			normal.Normalize()
			incidenceAngle = math.Atan2(-ball.Direction.Perp(normal), -ball.Direction.Dot(normal))
		}

		// Compute ball reflection
		sin, cos := math.Sincos(2 * incidenceAngle)
		ball.Direction = gm.Vector2{
			X: -ball.Direction.X*cos + ball.Direction.Y*sin,
			Y: -ball.Direction.X*sin - ball.Direction.Y*cos,
		}
		ball.Direction.Normalize()

		attractionLineIndex++
	}))
}
