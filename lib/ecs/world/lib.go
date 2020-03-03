package world

import (
	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	"arkanoid/lib/resources"
)

// World is the main ECS structure
type World struct {
	Manager    *ecs.Manager
	Components *c.Components
	Resources  *resources.Resources
}

// InitWorld initializes the world
func InitWorld() World {
	manager := &ecs.Manager{}
	components := c.InitComponents(manager)
	resources := resources.InitResources()

	return World{
		Manager:    manager,
		Components: components,
		Resources:  resources,
	}
}
