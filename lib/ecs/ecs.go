package ecs

import "github.com/yourbasic/bit"

//
// Manager
//

// Manager manages components and entities
type Manager struct {
	currentEntityIndex int
	components         []*Component
}

// NewComponent creates a new component
func (manager *Manager) NewComponent() *Component {
	component := &Component{
		data: make(map[Entity]interface{}),
	}

	manager.components = append(manager.components, component)
	return component
}

// NewEntity creates a new entity
func (manager *Manager) NewEntity() Entity {
	manager.currentEntityIndex++
	return Entity(manager.currentEntityIndex - 1)
}

// DeleteEntity removes entity for all associated components
func (manager *Manager) DeleteEntity(entity Entity) {
	for _, component := range manager.components {
		entity.RemoveComponent(component)
	}
}

// DeleteEntities removes entities for all associated components
func (manager *Manager) DeleteEntities(entities ...Entity) {
	for _, entity := range entities {
		manager.DeleteEntity(entity)
	}
}

//
// Entity
//

// Entity is an index
type Entity int

// AddComponent adds entity for component
func (entity Entity) AddComponent(component *Component, data interface{}) Entity {
	component.tag.Add(int(entity))
	component.data[entity] = data
	return entity
}

// RemoveComponent removes entity for component
func (entity Entity) RemoveComponent(component *Component) Entity {
	component.tag.Delete(int(entity))
	delete(component.data, entity)
	return entity
}

//
// Component
//

type component interface {
	_Tag() *bit.Set
	_Join(*bit.Set) *bit.Set
}

// Join returns tag describing intersection of components
func Join(components ...component) *bit.Set {
	tag := &bit.Set{}
	tag.Set(components[0]._Tag())
	for _, component := range components[1:] {
		tag = component._Join(tag)
	}
	return tag
}

// Component is a data storage
type Component struct {
	tag  bit.Set
	data map[Entity]interface{}
}

// Get returns data corresponding to entity
func (c *Component) Get(index int) interface{} {
	return c.data[Entity(index)]
}

// Not returns an inverted component used for filtering entities that don't have the component
func (c *Component) Not() *AntiComponent {
	return &AntiComponent{tag: c.tag}
}

func (c *Component) _Tag() *bit.Set {
	return &c.tag
}

func (c *Component) _Join(tag *bit.Set) *bit.Set {
	return tag.SetAnd(tag, &c.tag)
}

// AntiComponent is an inverted component used for filtering entities that don't have a component
type AntiComponent struct {
	tag bit.Set
}

func (a *AntiComponent) _Tag() *bit.Set {
	return &a.tag
}

func (a *AntiComponent) _Join(tag *bit.Set) *bit.Set {
	return tag.SetAndNot(tag, &a.tag)
}

//
// Others
//

// Visit is a decorator function for bit.Set.Visit() method
func Visit(f func(index int)) func(index int) bool {
	return func(index int) bool {
		f(index)
		return false
	}
}
