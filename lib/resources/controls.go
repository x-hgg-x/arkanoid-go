package resources

import (
	"arkanoid/lib/utils"

	"github.com/hajimehoshi/ebiten"
)

const (
	// PaddleAxis is the axis for moving paddle
	PaddleAxis = "Paddle"
	// ReleaseBallAction is the action for releasing ball
	ReleaseBallAction = "ReleaseBall"
)

// Key is a US keyboard key
type Key struct {
	Key ebiten.Key
}

// UnmarshalText fills structure fields from text data
func (k *Key) UnmarshalText(text []byte) error {
	k.Key = utils.KeyMap[string(text)]
	return nil
}

// MouseButton is a mouse button
type MouseButton struct {
	MouseButton ebiten.MouseButton
}

// UnmarshalText fills structure fields from text data
func (b *MouseButton) UnmarshalText(text []byte) error {
	b.MouseButton = utils.MouseButtonMap[string(text)]
	return nil
}

// ControllerButton is a gamepad button
type ControllerButton struct {
	ID            int
	GamepadButton ebiten.GamepadButton
}

// UnmarshalTOML fills structure fields from TOML data
func (b *ControllerButton) UnmarshalTOML(i interface{}) error {
	data := i.(map[string]interface{})
	b.ID = int(data["id"].(int64))
	b.GamepadButton = utils.GamepadButtonMap[data["button"].(string)]
	return nil
}

// Button can be a US keyboard key, a mouse button or a gamepad button
type Button struct {
	Type             string
	Key              *Key
	MouseButton      *MouseButton      `toml:"mouse_button"`
	ControllerButton *ControllerButton `toml:"controller"`
}

// Emulated is an emulated axis
type Emulated struct {
	Pos Button
	Neg Button
}

// ControllerAxis is a gamepad axis
type ControllerAxis struct {
	ID       int
	Axis     int
	Invert   bool
	DeadZone float64
}

// UnmarshalTOML fills structure fields from TOML data
func (a *ControllerAxis) UnmarshalTOML(i interface{}) error {
	data := i.(map[string]interface{})
	a.ID = int(data["id"].(int64))
	a.Axis = int(data["axis"].(int64))
	a.Invert = data["invert"].(bool)
	a.DeadZone = data["dead_zone"].(float64)
	return nil
}

// MouseAxis is a mouse axis
type MouseAxis struct {
	Axis int
}

// Axis can be an emulated axis, a gamepad axis or a mouse axis
type Axis struct {
	Type           string
	Emulated       *Emulated
	ControllerAxis *ControllerAxis `toml:"controller_axis"`
	MouseAxis      *MouseAxis      `toml:"mouse_axis"`
}

// Action contains buttons combinations
type Action = [][]Button

// Controls contains input controls
type Controls struct {
	// Axes contains axis controls, used for inputs represented by a float value from -1 to 1
	Axes map[string]Axis
	// Actions contains buttons combinations, used for general inputs
	Actions map[string]Action
}

// InputHandler contains input axis values and actions corresponding to specified controls
type InputHandler struct {
	// Axes contains input axis values
	Axes map[string]float64
	// Actions contains input actions
	Actions map[string]bool
}
