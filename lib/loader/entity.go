package loader

import (
	"fmt"
	"image/color"
	"reflect"

	c "arkanoid/lib/components"
	e "arkanoid/lib/ecs"
	"arkanoid/lib/utils"

	"github.com/BurntSushi/toml"
	"github.com/ByteArena/ecs"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

type componentList struct {
	SpriteRender  *c.SpriteRender
	Transform     *c.Transform
	Text          *c.Text
	UITransform   *c.UITransform
	MouseReactive *c.MouseReactive
	Paddle        *c.Paddle
	Ball          *c.Ball
	Block         *c.Block
}

type componentListData struct {
	SpriteRender  *spriteRenderData
	Transform     *c.Transform
	Text          *textData
	UITransform   *c.UITransform
	MouseReactive *c.MouseReactive
	Paddle        *c.Paddle
	Ball          *c.Ball
	Block         *c.Block
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
	return componentList{
		SpriteRender:  processSpriteRenderData(ecsData, data.SpriteRender),
		Transform:     data.Transform,
		Text:          processTextData(ecsData, data.Text),
		UITransform:   data.UITransform,
		MouseReactive: data.MouseReactive,
		Paddle:        data.Paddle,
		Ball:          data.Ball,
		Block:         data.Block,
	}
}

type fillData struct {
	Width  int
	Height int
	Color  [4]uint8
}

type spriteRenderData struct {
	Fill            *fillData
	SpriteSheetName string `toml:"sprite_sheet_name"`
	SpriteNumber    int    `toml:"sprite_number"`
}

func processSpriteRenderData(ecsData e.Ecs, spriteRenderData *spriteRenderData) *c.SpriteRender {
	if spriteRenderData == nil {
		return nil
	}
	if spriteRenderData.Fill != nil && spriteRenderData.SpriteSheetName != "" {
		utils.LogError(fmt.Errorf("fill and sprite_sheet_name fields are exclusive"))
	}

	// Sprite is included in sprite sheet
	if spriteRenderData.SpriteSheetName != "" {
		// Add reference to sprite sheet from its name
		spriteSheet, ok := (*ecsData.Resources.SpriteSheets)[spriteRenderData.SpriteSheetName]
		if !ok {
			utils.LogError(fmt.Errorf("unable to find sprite sheet with name '%s'", spriteRenderData.SpriteSheetName))
		}
		return &c.SpriteRender{
			SpriteSheet:  &spriteSheet,
			SpriteNumber: spriteRenderData.SpriteNumber,
		}
	}

	// Sprite is a colored rectangle
	textureImage, err := ebiten.NewImage(spriteRenderData.Fill.Width, spriteRenderData.Fill.Height, ebiten.FilterNearest)
	utils.LogError(err)

	textureImage.Fill(color.RGBA{
		R: spriteRenderData.Fill.Color[0],
		G: spriteRenderData.Fill.Color[1],
		B: spriteRenderData.Fill.Color[2],
		A: spriteRenderData.Fill.Color[3],
	})

	return &c.SpriteRender{
		SpriteSheet: &c.SpriteSheet{
			Texture: c.Texture{Image: textureImage},
			Sprites: []c.Sprite{c.Sprite{X: 0, Y: 0, Width: spriteRenderData.Fill.Width, Height: spriteRenderData.Fill.Height}},
		},
		SpriteNumber: 0,
	}
}

type fontFaceOptions struct {
	Size              float64
	DPI               float64
	Hinting           string
	GlyphCacheEntries int `toml:"glyph_cache_entries"`
	SubPixelsX        int `toml:"sub_pixels_x"`
	SubPixelsY        int `toml:"sub_pixels_y"`
}

var hintingMap = map[string]font.Hinting{
	"":         font.HintingNone,
	"None":     font.HintingNone,
	"Vertical": font.HintingVertical,
	"Full":     font.HintingFull,
}

type fontFaceData struct {
	Font    string
	Options fontFaceOptions
}

type textData struct {
	ID       string
	Text     string
	FontFace fontFaceData `toml:"font_face"`
	Color    [4]uint8
}

func processTextData(ecsData e.Ecs, textData *textData) *c.Text {
	if textData == nil {
		return nil
	}

	// Search font from its name
	textFont, ok := (*ecsData.Resources.Fonts)[textData.FontFace.Font]
	if !ok {
		utils.LogError(fmt.Errorf("unable to find font with name '%s'", textData.FontFace.Font))
	}

	options := &truetype.Options{
		Size:              textData.FontFace.Options.Size,
		DPI:               textData.FontFace.Options.DPI,
		Hinting:           hintingMap[textData.FontFace.Options.Hinting],
		GlyphCacheEntries: textData.FontFace.Options.GlyphCacheEntries,
		SubPixelsX:        textData.FontFace.Options.SubPixelsX,
		SubPixelsY:        textData.FontFace.Options.SubPixelsY,
	}

	return &c.Text{
		ID:       textData.ID,
		Text:     textData.Text,
		FontFace: truetype.NewFace(textFont.Font, options),
		Color:    color.RGBA{R: textData.Color[0], G: textData.Color[1], B: textData.Color[2], A: textData.Color[3]},
	}
}
