package ecs

import (
	c "arkanoid/lib/components"
	"arkanoid/lib/resources"
	"arkanoid/lib/systems"

	"github.com/ByteArena/ecs"
)

// Ecs is the main ECS structure
type Ecs struct {
	Manager    *ecs.Manager
	Components *c.Components
	Views      *systems.Views
	Resources  *resources.Resources
}

// InitEcs initializes the main ECS structure
func InitEcs() Ecs {
	manager := ecs.NewManager()
	components := c.InitComponents(manager)
	views := systems.InitViews(manager, components)
	resources := resources.InitResources()

	return Ecs{
		Manager:    manager,
		Components: components,
		Views:      views,
		Resources:  resources,
	}
}
