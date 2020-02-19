package components

import (
	"reflect"

	"github.com/ByteArena/ecs"
)

// Components contains references to all components
type Components struct {
	Sprite    *ecs.Component
	Transform *ecs.Component
}

// InitComponents initializes components
func InitComponents(manager *ecs.Manager) *Components {
	components := &Components{}

	v := reflect.ValueOf(components).Elem()
	for i := 0; i < v.NumField(); i++ {
		v.Field(i).Set(reflect.ValueOf(manager.NewComponent()))
	}

	return components
}
