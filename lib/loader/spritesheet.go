package loader

import (
	c "arkanoid/lib/components"
	"arkanoid/lib/utils"

	"github.com/BurntSushi/toml"
)

type spriteSheetMetadata struct {
	SpriteSheets map[string]c.SpriteSheet `toml:"sprite_sheet"`
}

// LoadSpriteSheets loads sprite sheets from a TOML file
func LoadSpriteSheets(spriteSheetMetadataPath string) map[string]c.SpriteSheet {
	var spriteSheetMetadata spriteSheetMetadata
	_, err := toml.DecodeFile("assets/metadata/spritesheets/spritesheets.toml", &spriteSheetMetadata)
	utils.LogError(err)
	return spriteSheetMetadata.SpriteSheets
}
