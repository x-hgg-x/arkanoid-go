package loader

import (
	"fmt"
	"reflect"

	"arkanoid/lib/resources"
	"arkanoid/lib/utils"

	"github.com/BurntSushi/toml"
)

type controlsConfig struct {
	Controls resources.Controls `toml:"controls"`
}

// LoadControls loads controls from a TOML file
func LoadControls(controlsConfigPath string, axes []string, actions []string) (resources.Controls, resources.InputHandler) {
	var controlsConfig controlsConfig
	_, err := toml.DecodeFile(controlsConfigPath, &controlsConfig)
	utils.LogError(err)

	var inputHandler resources.InputHandler
	inputHandler.Axes = make(map[string]float64)
	inputHandler.Actions = make(map[string]bool)

	// Check axes
	for _, axis := range axes {
		if _, ok := controlsConfig.Controls.Axes[axis]; !ok {
			utils.LogError(fmt.Errorf("unable to find controls for axis '%s'", axis))
		}
		inputHandler.Axes[axis] = 0
	}

	// Check actions
	for _, action := range actions {
		if _, ok := controlsConfig.Controls.Actions[action]; !ok {
			utils.LogError(fmt.Errorf("unable to find controls for action '%s'", action))
		}
		inputHandler.Actions[action] = false
	}

	// Set "Type" field for axes
	for k, v := range controlsConfig.Controls.Axes {
		controlsConfig.Controls.Axes[k] = setTypeFields(reflect.ValueOf(&v)).Interface().(resources.Axis)
	}

	// Set "Type" field for actions
	for k, v := range controlsConfig.Controls.Actions {
		for i := range v {
			for j := range v[i] {
				setTypeFields(reflect.ValueOf(&v[i][j]))
			}
		}
		controlsConfig.Controls.Actions[k] = v
	}

	return controlsConfig.Controls, inputHandler
}

// Set "Type" field to the name of the non null field (only one non null field is allowed)
func setTypeFields(v reflect.Value) reflect.Value {
	v = reflect.Indirect(v)
	for iField := 0; iField < v.NumField(); iField++ {
		if v.Type().Field(iField).Name == "Type" {
			for iSubField := 0; iSubField < v.NumField(); iSubField++ {
				subField := v.Field(iSubField)
				if subField.Kind() == reflect.Ptr && !subField.IsNil() {
					typeName := subField.Elem().Type().Name()
					if v.FieldByName("Type").String() != "" {
						utils.LogError(fmt.Errorf("duplicate fields found: %s, %s", v.FieldByName("Type").String(), typeName))
					}
					v.FieldByName("Type").SetString(typeName)
				}
			}
		}

		field := reflect.Indirect(v.Field(iField))
		if field.Kind() == reflect.Struct {
			setTypeFields(field)
		}
	}
	return v
}
