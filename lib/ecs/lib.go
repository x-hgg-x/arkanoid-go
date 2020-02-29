package ecs

import (
	c "arkanoid/lib/components"
	"arkanoid/lib/resources"
	"arkanoid/lib/systems"

	"github.com/ByteArena/ecs"
)

// World is the main ECS structure
type World struct {
	Manager    *ecs.Manager
	Components *c.Components
	Views      *systems.Views
	Resources  *resources.Resources
}

// InitWorld initializes the world
func InitWorld() World {
	manager := ecs.NewManager()
	components := c.InitComponents(manager)
	views := systems.InitViews(manager, components)
	resources := resources.InitResources()

	return World{
		Manager:    manager,
		Components: components,
		Views:      views,
		Resources:  resources,
	}
}
