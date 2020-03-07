package resources

import (
	"arkanoid/lib/ecs"

	"github.com/ByteArena/box2d"
)

// B2PixelRatio is the number of pixels representing 1 meter in box2D world
const B2PixelRatio = 50

// BlockCollisionEvent is triggered when a block collision occurs
type BlockCollisionEvent struct {
	Entity ecs.Entity
}

// Events contains game events for communication between game systems
type Events struct {
	BlockCollisionEvents []BlockCollisionEvent
}

// Game contains game resources
type Game struct {
	CollisionWorld *box2d.B2World
	Events         Events
}
