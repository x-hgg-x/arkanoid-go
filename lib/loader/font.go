package loader

import (
	"arkanoid/lib/resources"
	"arkanoid/lib/utils"

	"github.com/BurntSushi/toml"
)

type fontMetadata struct {
	Fonts map[string]resources.Font `toml:"font"`
}

// LoadFonts loads fonts from a TOML file
func LoadFonts(fontPath string) map[string]resources.Font {
	var fontMetadata fontMetadata
	_, err := toml.DecodeFile("assets/metadata/fonts/fonts.toml", &fontMetadata)
	utils.LogError(err)
	return fontMetadata.Fonts
}
