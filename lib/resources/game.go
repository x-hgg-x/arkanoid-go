package resources

import "github.com/ByteArena/box2d"

// B2PixelRatio is the number of pixels representing 1 meter in box2D world
const B2PixelRatio = 50

// Game contains game resources
type Game struct {
	CollisionWorld *box2d.B2World
}
