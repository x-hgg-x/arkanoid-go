package systems

import (
	c "arkanoid/lib/components"

	"github.com/ByteArena/ecs"
)

// Views contains references to all views
type Views struct {
	SpriteView *ecs.View
	PaddleView *ecs.View
	BallView   *ecs.View
}

// InitViews initializes views
func InitViews(manager *ecs.Manager, components *c.Components) *Views {
	return &Views{
		SpriteView: manager.CreateView(ecs.BuildTag(components.SpriteRender, components.Transform)),
		PaddleView: manager.CreateView(ecs.BuildTag(components.Paddle, components.Transform)),
		BallView:   manager.CreateView(ecs.BuildTag(components.Ball, components.Transform)),
	}
}
