package loader

import (
	"os"

	gc "github.com/x-hgg-x/arkanoid-go/lib/components"

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

func loadGameComponents(entityMetadataContent []byte, world w.World) []interface{} {
	var entityGameMetadata entityGameMetadata
	utils.Try(toml.Decode(string(entityMetadataContent), &entityGameMetadata))

	gameComponentList := make([]interface{}, len(entityGameMetadata.Entities))
	for iEntity, entity := range entityGameMetadata.Entities {
		gameComponentList[iEntity] = entity.Components
	}
	return gameComponentList
}

// PreloadEntities preloads entities with components
func PreloadEntities(entityMetadataPath string, world w.World) loader.EntityComponentList {
	entityMetadataContent := utils.Try(os.ReadFile(entityMetadataPath))

	return loader.EntityComponentList{
		Engine: loader.LoadEngineComponents(entityMetadataContent, world),
		Game:   loadGameComponents(entityMetadataContent, world),
	}
}
