package loader

import (
	gc "arkanoid/lib/components"

	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/BurntSushi/toml"
	ecs "github.com/x-hgg-x/goecs"
)

type gameComponentList struct {
	Paddle         *gc.Paddle
	Ball           *gc.Ball
	StickyBall     *gc.StickyBall
	AttractionLine *gc.AttractionLine
	Block          *gc.Block
}

type gameComponents struct {
	Game gameComponentList
}

type entity struct {
	Components gameComponents
}

type entityGameMetadata struct {
	Entities []entity `toml:"entity"`
}

func loadGameComponents(entityMetadataPath string, world w.World) []gameComponentList {
	var entityGameMetadata entityGameMetadata
	_, err := toml.DecodeFile(entityMetadataPath, &entityGameMetadata)
	utils.LogError(err)

	gameComponentList := make([]gameComponentList, len(entityGameMetadata.Entities))
	for iEntity, entity := range entityGameMetadata.Entities {
		gameComponentList[iEntity] = entity.Components.Game
	}
	return gameComponentList
}

// LoadEntities creates entities with components from a TOML file
func LoadEntities(entityMetadataPath string, world w.World) []ecs.Entity {
	engineComponentList := loader.LoadEngineComponents(entityMetadataPath, world)
	gameComponentList := loadGameComponents(entityMetadataPath, world)

	entities := make([]ecs.Entity, len(engineComponentList))
	for iEntity := range engineComponentList {
		// Add components to a new entity
		entities[iEntity] = world.Manager.NewEntity()
		loader.AddEntityComponents(entities[iEntity], world.Components.Engine, engineComponentList[iEntity])
		loader.AddEntityComponents(entities[iEntity], world.Components.Game, gameComponentList[iEntity])
	}
	return entities
}
