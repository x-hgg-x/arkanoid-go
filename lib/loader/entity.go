package loader

import (
	"fmt"
	"reflect"

	c "arkanoid/lib/components"
	e "arkanoid/lib/ecs"
	"arkanoid/lib/utils"

	"github.com/BurntSushi/toml"
	"github.com/ByteArena/ecs"
)

type componentList struct {
	SpriteRender *c.SpriteRender
	Transform    *c.Transform
	Paddle       *c.Paddle
	Ball         *c.Ball
	Block        *c.Block
}

type spriteRenderData struct {
	SpriteSheetName string `toml:"sprite_sheet_name"`
	SpriteNumber    int    `toml:"sprite_number"`
}

type componentListData struct {
	SpriteRender *spriteRenderData
	Transform    *c.Transform
	Paddle       *c.Paddle
	Ball         *c.Ball
	Block        *c.Block
}

type entity struct {
	Components componentListData
}

type entityMetadata struct {
	Entities []entity `toml:"entity"`
}

// LoadEntities creates entities with components from a TOML file
func LoadEntities(entityMetadataPath string, ecsData e.Ecs) []*ecs.Entity {
	var entityMetadata entityMetadata
	_, err := toml.DecodeFile(entityMetadataPath, &entityMetadata)
	utils.LogError(err)

	entities := make([]*ecs.Entity, len(entityMetadata.Entities))
	for iEntity, entity := range entityMetadata.Entities {
		// Add components to a new entity
		entities[iEntity] = addEntityComponents(ecsData.Manager.NewEntity(), ecsData.Components, processComponentsListData(ecsData, entity.Components))
	}
	return entities
}

func addEntityComponents(entity *ecs.Entity, ecsComponentList *c.Components, components componentList) *ecs.Entity {
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

func processComponentsListData(ecsData e.Ecs, data componentListData) componentList {
	// SpriteRender
	var spriteRender *c.SpriteRender
	if data.SpriteRender != nil {
		// Add reference to sprite sheet from its name
		if spriteSheet, ok := (*ecsData.Resources.SpriteSheets)[data.SpriteRender.SpriteSheetName]; ok {
			spriteRender = &c.SpriteRender{
				SpriteSheet:  &spriteSheet,
				SpriteNumber: data.SpriteRender.SpriteNumber,
			}
		} else {
			utils.LogError(fmt.Errorf("unable to find sprite sheet with name '%s'", data.SpriteRender.SpriteSheetName))
		}
	}

	return componentList{
		SpriteRender: spriteRender,
		Transform:    data.Transform,
		Paddle:       data.Paddle,
		Ball:         data.Ball,
		Block:        data.Block,
	}
}
