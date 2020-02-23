package loader

import (
	c "arkanoid/components"
	"arkanoid/utils"

	"github.com/BurntSushi/toml"
)

type spriteSheetMetadata struct {
	SpriteSheets map[string]c.SpriteSheet `toml:"sprite_sheet"`
}

// LoadSpriteSheet loads sprite sheets from a TOML file
func LoadSpriteSheet(spriteSheetMetadataPath string) map[string]c.SpriteSheet {
	var spriteSheetMetadata spriteSheetMetadata
	_, err := toml.DecodeFile("assets/metadata/spritesheets/spritesheets.toml", &spriteSheetMetadata)
	utils.LogError(err)
	return spriteSheetMetadata.SpriteSheets
}
