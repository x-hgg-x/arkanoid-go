package components

import (
	"image/color"

	"arkanoid/lib/math"

	"golang.org/x/image/font"
)

// Text component
type Text struct {
	ID       string
	Text     string
	FontFace font.Face
	Color    color.RGBA
}

// Pivot variants
const (
	TopLeft      = "TopLeft"
	TopMiddle    = "TopMiddle"
	TopRight     = "TopRight"
	MiddleLeft   = "MiddleLeft"
	Middle       = "Middle"
	MiddleRight  = "MiddleRight"
	BottomLeft   = "BottomLeft"
	BottomMiddle = "BottomMiddle"
	BottomRight  = "BottomRight"
)

// UITransform component
type UITransform struct {
	// Translation defines the position of the pivot relative to the origin.
	Translation math.VectorInt2
	// Pivot defines the position of the element relative to its translation (default is Middle).
	Pivot string
}
