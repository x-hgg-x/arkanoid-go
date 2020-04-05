package loader

import (
	gc "github.com/x-hgg-x/arkanoid-go/lib/components"

	ecs "github.com/x-hgg-x/goecs"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/BurntSushi/toml"
)

type gameComponentList struct {
	Paddle         *gc.Paddle
	Ball           *gc.Ball
	StickyBall     *gc.StickyBall
	AttractionLine *gc.AttractionLine
	Block          *gc.Block
}

type entity struct {
	Components gameComponentList
}

type entityGameMetadata struct {
	Entities []entity `toml:"entity"`
}

func loadGameComponents(entityMetadataPath string, world w.World) []interface{} {
	var entityGameMetadata entityGameMetadata
	_, err := toml.DecodeFile(entityMetadataPath, &entityGameMetadata)
	utils.LogError(err)

	gameComponentList := make([]interface{}, len(entityGameMetadata.Entities))
	for iEntity, entity := range entityGameMetadata.Entities {
		gameComponentList[iEntity] = entity.Components
	}
	return gameComponentList
}

// LoadEntities creates entities with components from a TOML file
func LoadEntities(entityMetadataPath string, world w.World) []ecs.Entity {
	gameComponentList := loadGameComponents(entityMetadataPath, world)
	return loader.LoadEntities(entityMetadataPath, world, gameComponentList)
}
