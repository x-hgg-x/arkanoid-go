package components

import (
	"reflect"

	"github.com/ByteArena/ecs"
)

// Components contains references to all components
type Components struct {
	SpriteRender  *ecs.Component
	Transform     *ecs.Component
	Text          *ecs.Component
	UITransform   *ecs.Component
	MouseReactive *ecs.Component
	Paddle        *ecs.Component
	Ball          *ecs.Component
	Block         *ecs.Component
}

// InitComponents initializes components
func InitComponents(manager *ecs.Manager) *Components {
	components := &Components{}

	v := reflect.ValueOf(components).Elem()
	for iField := 0; iField < v.NumField(); iField++ {
		v.Field(iField).Set(reflect.ValueOf(manager.NewComponent()))
	}

	return components
}
