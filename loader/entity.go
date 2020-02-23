package loader

import (
	"log"
	"reflect"

	c "arkanoid/components"
	e "arkanoid/ecs"
	"arkanoid/utils"

	"github.com/BurntSushi/toml"
	"github.com/ByteArena/ecs"
)

type componentList struct {
	SpriteRender *c.SpriteRender
	Transform    *c.Transform
	Paddle       *c.Paddle
	Ball         *c.Ball
}

type entity struct {
	Components componentList
}

type entityMetadata struct {
	Entities []entity `toml:"entity"`
}

// LoadEntities create entities with components from a TOML file
func LoadEntities(entityMetadataPath string, ecsData e.Ecs, spriteSheets map[string]c.SpriteSheet) []*ecs.Entity {
	var entityMetadata entityMetadata
	_, err := toml.DecodeFile(entityMetadataPath, &entityMetadata)
	utils.LogError(err)

	entities := make([]*ecs.Entity, len(entityMetadata.Entities))
	for iEntity, entity := range entityMetadata.Entities {
		// Add reference to sprite sheet from its name
		if entity.Components.SpriteRender != nil {
			if spriteSheet, ok := spriteSheets[entity.Components.SpriteRender.SpriteSheetName]; ok {
				entity.Components.SpriteRender.SpriteSheet = &spriteSheet
			} else {
				log.Fatalf("Unable to find sprite sheet with name '%s'", entity.Components.SpriteRender.SpriteSheetName)
			}
		}

		// Add components to a new entity
		entities[iEntity] = addEntity(ecsData.Manager.NewEntity(), ecsData.Components, entity.Components)
	}
	return entities
}

func addEntity(entity *ecs.Entity, ecsComponentList *c.Components, components componentList) *ecs.Entity {
	v := reflect.ValueOf(components)
	for iField := 0; iField < v.NumField(); iField++ {
		if !v.Field(iField).IsNil() {
			component := v.Field(iField)
			componentName := component.Elem().Type().Name()
			ecsComponent := reflect.ValueOf(ecsComponentList).Elem().FieldByName(componentName).Interface().(*ecs.Component)
			entity.AddComponent(ecsComponent, component.Interface())
		}
	}
	return entity
}
