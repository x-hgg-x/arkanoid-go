package resources

import (
	"io/ioutil"

	"arkanoid/lib/utils"

	"github.com/golang/freetype/truetype"
)

// Font structure
type Font struct {
	Font *truetype.Font
}

// UnmarshalTOML fills structure fields from TOML data
func (f *Font) UnmarshalTOML(i interface{}) error {
	fontFile, err := ioutil.ReadFile(i.(map[string]interface{})["font"].(string))
	utils.LogError(err)
	f.Font, err = truetype.Parse(fontFile)
	utils.LogError(err)
	return nil
}
