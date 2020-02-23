package ecs

import (
	c "arkanoid/lib/components"
	"arkanoid/lib/systems"

	"github.com/ByteArena/ecs"
)

// Ecs is the main ECS structure
type Ecs struct {
	Manager    *ecs.Manager
	Components *c.Components
	Views      *systems.Views
}

// InitEcs initializes the main ECS structure
func InitEcs() Ecs {
	manager := ecs.NewManager()
	components := c.InitComponents(manager)
	views := systems.InitViews(manager, components)

	return Ecs{
		Manager:    manager,
		Components: components,
		Views:      views,
	}
}
