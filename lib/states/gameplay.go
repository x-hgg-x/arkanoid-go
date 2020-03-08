package states

import (
	"fmt"

	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"
	"arkanoid/lib/loader"
	"arkanoid/lib/resources"
	g "arkanoid/lib/systems/game"
	i "arkanoid/lib/systems/input"
	s "arkanoid/lib/systems/sprite"
	u "arkanoid/lib/systems/ui"
	"arkanoid/lib/utils"

	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// GameplayState is the main game state
type GameplayState struct {
	game []ecs.Entity
}

func (st *GameplayState) onPause(world w.World)  {}
func (st *GameplayState) onResume(world w.World) {}

func (st *GameplayState) onStart(world w.World) {
	// Load game and ui entities
	st.game = append(st.game, loader.LoadEntities("assets/metadata/entities/background.toml", world)...)
	st.game = append(st.game, loader.LoadEntities("assets/metadata/entities/game.toml", world)...)
	st.game = append(st.game, loader.LoadEntities("assets/metadata/entities/ui/score.toml", world)...)
	st.game = append(st.game, loader.LoadEntities("assets/metadata/entities/ui/life.toml", world)...)

	world.Resources.Game = resources.NewGame()
	initializeCollisionWorld(world)
}

func (st *GameplayState) onStop(world w.World) {
	destroyCollisionWorld(world)
	world.Resources.Game = nil
	world.Manager.DeleteEntities(st.game...)
}

func (st *GameplayState) update(world w.World, screen *ebiten.Image) transition {
	i.InputSystem(world)
	u.UISystem(world)

	g.MovePaddleSystem(world)
	g.StickyBallSystem(world)
	g.BallAttractionSystem(world)
	g.BallAttractionVfxSystem(world)
	g.MoveBallSystem(world)
	g.CollisionSystem(world)
	g.BlockHealthSystem(world)
	g.LifeSystem(world)
	g.ScoreSystem(world)

	s.TransformSystem(world)
	s.RenderSpriteSystem(world, screen)
	u.RenderUISystem(world, screen)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return transition{transType: transPush, newStates: []state{&PauseMenuState{}}}
	}

	switch world.Resources.Game.StateEvent {
	case resources.StateEventGameOver:
		world.Resources.Game.StateEvent = resources.StateEventNone
		return transition{transType: transSwitch, newStates: []state{&GameOverState{Score: world.Resources.Game.Score}}}
	case resources.StateEventLevelComplete:
		world.Resources.Game.StateEvent = resources.StateEventNone
		return transition{transType: transSwitch, newStates: []state{&LevelCompleteState{Score: world.Resources.Game.Score}}}
	}

	return transition{}
}

func initializeCollisionWorld(world w.World) {
	// Init Box2D world
	collisionWorld := box2d.MakeB2World(box2d.MakeB2Vec2(0, 0))

	// Create paddle body
	paddles := ecs.Join(world.Components.Paddle, world.Components.Transform)
	if paddles.Empty() {
		utils.LogError(fmt.Errorf("unable to find paddle"))
	}
	firstPaddle := ecs.Entity(paddles.Next(-1))
	paddle := world.Components.Paddle.Get(firstPaddle).(*c.Paddle)

	paddleDef := box2d.MakeB2BodyDef()
	paddleBody := collisionWorld.CreateBody(&paddleDef)
	paddleShape := box2d.MakeB2PolygonShape()
	paddleShape.SetAsBox(paddle.Width/2/resources.B2PixelRatio, paddle.Height/2/resources.B2PixelRatio)
	paddleBody.CreateFixtureFromDef(&box2d.B2FixtureDef{Shape: &paddleShape})
	paddleBody.SetUserData(firstPaddle)
	paddle.Body = paddleBody

	// Create blocks bodies
	ecs.Join(world.Components.Block, world.Components.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		block := world.Components.Block.Get(entity).(*c.Block)
		blockTranslation := world.Components.Transform.Get(entity).(*c.Transform).Translation

		blockDef := box2d.MakeB2BodyDef()
		blockDef.Position.Set(blockTranslation.X/resources.B2PixelRatio, blockTranslation.Y/resources.B2PixelRatio)
		blockBody := collisionWorld.CreateBody(&blockDef)
		blockShape := box2d.MakeB2PolygonShape()
		blockShape.SetAsBox(block.Width/2/resources.B2PixelRatio, block.Height/2/resources.B2PixelRatio)
		blockBody.CreateFixtureFromDef(&box2d.B2FixtureDef{Shape: &blockShape})
		blockBody.SetUserData(entity)
		block.Body = blockBody
	}))

	// Create balls bodies
	ecs.Join(world.Components.Ball, world.Components.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		ball := world.Components.Ball.Get(entity).(*c.Ball)

		ballDef := box2d.MakeB2BodyDef()
		ballDef.Type = box2d.B2BodyType.B2_dynamicBody
		ballBody := collisionWorld.CreateBody(&ballDef)
		ballShape := box2d.MakeB2CircleShape()
		ballShape.M_radius = ball.Radius / resources.B2PixelRatio
		ballBody.CreateFixtureFromDef(&box2d.B2FixtureDef{Shape: &ballShape})
		ballBody.SetUserData(entity)
		ball.Body = ballBody
	}))

	world.Resources.Game.CollisionWorld = &collisionWorld
}

func destroyCollisionWorld(world w.World) {
	world.Resources.Game.CollisionWorld = nil
}
