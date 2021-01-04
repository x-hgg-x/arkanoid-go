package states

import (
	"fmt"

	gc "github.com/x-hgg-x/arkanoid-go/lib/components"
	"github.com/x-hgg-x/arkanoid-go/lib/resources"
	g "github.com/x-hgg-x/arkanoid-go/lib/systems"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// GameplayState is the main game state
type GameplayState struct{}

// OnStart method
func (st *GameplayState) OnStart(world w.World) {
	// Load game and ui entities
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Game.Background)
	loader.AddEntities(world, prefabs.Game.Game)
	loader.AddEntities(world, prefabs.Game.Score)
	loader.AddEntities(world, prefabs.Game.Life)

	world.Resources.Game = resources.NewGame()
	initializeCollisionWorld(world)

	ebiten.SetCursorMode(ebiten.CursorModeHidden)
}

// OnPause method
func (st *GameplayState) OnPause(world w.World) {
	ebiten.SetCursorMode(ebiten.CursorModeVisible)
}

// OnResume method
func (st *GameplayState) OnResume(world w.World) {
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
}

// OnStop method
func (st *GameplayState) OnStop(world w.World) {
	destroyCollisionWorld(world)
	world.Resources.Game = nil
	world.Manager.DeleteAllEntities()

	ebiten.SetCursorMode(ebiten.CursorModeVisible)
}

// Update method
func (st *GameplayState) Update(world w.World) states.Transition {
	g.MovePaddleSystem(world)
	g.StickyBallSystem(world)
	g.BallAttractionSystem(world)
	g.BallAttractionVfxSystem(world)
	g.MoveBallSystem(world)
	g.CollisionSystem(world)
	g.BlockHealthSystem(world)
	g.LifeSystem(world)
	g.ScoreSystem(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&PauseMenuState{}}}
	}

	gameResources := world.Resources.Game.(*resources.Game)
	switch gameResources.StateEvent {
	case resources.StateEventGameOver:
		gameResources.StateEvent = resources.StateEventNone
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&GameOverState{Score: gameResources.Score}}}
	case resources.StateEventLevelComplete:
		gameResources.StateEvent = resources.StateEventNone
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&LevelCompleteState{Score: gameResources.Score}}}
	}

	return states.Transition{}
}

func initializeCollisionWorld(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	// Init Box2D world
	collisionWorld := box2d.MakeB2World(box2d.MakeB2Vec2(0, 0))

	// Create paddle body
	firstPaddle := ecs.GetFirst(world.Manager.Join(gameComponents.Paddle, world.Components.Engine.Transform))
	if firstPaddle == nil {
		utils.LogError(fmt.Errorf("unable to find paddle"))
	}
	paddle := gameComponents.Paddle.Get(ecs.Entity(*firstPaddle)).(*gc.Paddle)

	paddleDef := box2d.MakeB2BodyDef()
	paddleBody := collisionWorld.CreateBody(&paddleDef)
	paddleShape := box2d.MakeB2PolygonShape()
	paddleShape.SetAsBox(paddle.Width/2/resources.B2PixelRatio, paddle.Height/2/resources.B2PixelRatio)
	paddleBody.CreateFixtureFromDef(&box2d.B2FixtureDef{Shape: &paddleShape, Filter: box2d.MakeB2Filter()})
	paddleBody.SetUserData(*firstPaddle)
	paddle.Body = paddleBody

	// Create blocks bodies
	world.Manager.Join(gameComponents.Block, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		block := gameComponents.Block.Get(entity).(*gc.Block)
		blockTranslation := world.Components.Engine.Transform.Get(entity).(*ec.Transform).Translation

		blockDef := box2d.MakeB2BodyDef()
		blockDef.Position.Set(blockTranslation.X/resources.B2PixelRatio, blockTranslation.Y/resources.B2PixelRatio)
		blockBody := collisionWorld.CreateBody(&blockDef)
		blockShape := box2d.MakeB2PolygonShape()
		blockShape.SetAsBox(block.Width/2/resources.B2PixelRatio, block.Height/2/resources.B2PixelRatio)
		blockBody.CreateFixtureFromDef(&box2d.B2FixtureDef{Shape: &blockShape, Filter: box2d.MakeB2Filter()})
		blockBody.SetUserData(entity)
		block.Body = blockBody
	}))

	// Create balls bodies
	world.Manager.Join(gameComponents.Ball, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		ball := gameComponents.Ball.Get(entity).(*gc.Ball)

		ballDef := box2d.MakeB2BodyDef()
		ballDef.Type = box2d.B2BodyType.B2_dynamicBody
		ballBody := collisionWorld.CreateBody(&ballDef)
		ballShape := box2d.MakeB2CircleShape()
		ballShape.M_radius = ball.Radius / resources.B2PixelRatio
		ballBody.CreateFixtureFromDef(&box2d.B2FixtureDef{Shape: &ballShape, Filter: box2d.MakeB2Filter()})
		ballBody.SetUserData(entity)
		ball.Body = ballBody
	}))

	world.Resources.Game.(*resources.Game).CollisionWorld = &collisionWorld
}

func destroyCollisionWorld(world w.World) {
	world.Resources.Game.(*resources.Game).CollisionWorld = nil
}
