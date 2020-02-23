package resources

import c "arkanoid/lib/components"

// ScreenDimensions contains current screen dimensions
type ScreenDimensions struct {
	Width  int
	Height int
}

// Resources contains reference to data not related to any entity
type Resources struct {
	ScreenDimensions *ScreenDimensions
	SpriteSheets     *map[string]c.SpriteSheet
}

// InitResources initializes resources
func InitResources() *Resources {
	return &Resources{}
}
