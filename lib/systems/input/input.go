package inputsystem

import (
	"math"

	w "arkanoid/lib/ecs/world"
	"arkanoid/lib/resources"

	"github.com/hajimehoshi/ebiten"
)

// InputSystem updates input axis values and actions
func InputSystem(world w.World) {
	for k, v := range world.Resources.Controls.Axes {
		world.Resources.InputHandler.Axes[k] = getAxisValue(world, v)
	}

	for k, v := range world.Resources.Controls.Actions {
		world.Resources.InputHandler.Actions[k] = isActionDone(v)
	}
}

func getAxisValue(world w.World, axis resources.Axis) float64 {
	axisValue := 0.0

	switch axis.Type {
	case "Emulated":
		if isPressed(axis.Emulated.Pos) {
			axisValue++
		}
		if isPressed(axis.Emulated.Neg) {
			axisValue--
		}
	case "ControllerAxis":
		deadZone := math.Abs(axis.ControllerAxis.DeadZone)
		axisValue = ebiten.GamepadAxis(axis.ControllerAxis.ID, axis.ControllerAxis.Axis)

		if axisValue < -deadZone {
			axisValue = (axisValue + deadZone) / (1.0 - deadZone)
		} else if axisValue > deadZone {
			axisValue = (axisValue - deadZone) / (1.0 - deadZone)
		} else {
			axisValue = 0
		}

		if axis.ControllerAxis.Invert {
			axisValue *= -1
		}
	case "MouseAxis":
		screenWidth := float64(world.Resources.ScreenDimensions.Width)
		screenHeight := float64(world.Resources.ScreenDimensions.Height)

		x, y := ebiten.CursorPosition()
		switch axis.MouseAxis.Axis {
		case 0:
			axisValue = float64(x) / screenWidth
		case 1:
			axisValue = (screenHeight - float64(y)) / screenHeight
		}
	}
	return axisValue
}

func isActionDone(action resources.Action) bool {
	for _, combination := range action {
		actionDone := true
		for _, button := range combination {
			actionDone = actionDone && isPressed(button)
		}
		if actionDone {
			return true
		}
	}
	return false
}

func isPressed(b resources.Button) bool {
	switch b.Type {
	case "Key":
		return ebiten.IsKeyPressed(b.Key.Key)
	case "MouseButton":
		return ebiten.IsMouseButtonPressed(b.MouseButton.MouseButton)
	case "ControllerButton":
		return ebiten.IsGamepadButtonPressed(b.ControllerButton.ID, b.ControllerButton.GamepadButton)
	}
	return false
}
