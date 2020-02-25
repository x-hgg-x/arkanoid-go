package inputsystem

import (
	"math"

	e "arkanoid/lib/ecs"
	"arkanoid/lib/resources"

	"github.com/hajimehoshi/ebiten"
)

// InputSystem updates input axis values and actions
func InputSystem(ecs e.Ecs) {
	for k, v := range ecs.Resources.Controls.Axes {
		ecs.Resources.InputHandler.Axes[k] = getAxisValue(ecs, v)
	}

	for k, v := range ecs.Resources.Controls.Actions {
		ecs.Resources.InputHandler.Actions[k] = isActionDone(v)
	}
}

func getAxisValue(ecs e.Ecs, axis resources.Axis) float64 {
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
		screenWidth := float64(ecs.Resources.ScreenDimensions.Width)
		screenHeight := float64(ecs.Resources.ScreenDimensions.Height)

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
