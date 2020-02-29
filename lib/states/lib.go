package states

import (
	"fmt"
	"os"

	e "arkanoid/lib/ecs"
	"arkanoid/lib/utils"

	"github.com/hajimehoshi/ebiten"
)

type trans int

const (
	// No transition
	transNone trans = iota
	// Remove the active state and resume the next state
	transPop
	// Pause the active state and add new states to the stack
	transPush
	// Remove the active state and replace it by a new one
	transSwitch
	// Remove all states and insert a new stack
	transReplace
	// Remove all states and quit
	transQuit
)

type transition struct {
	transType trans
	newStates []state
}

type state interface {
	// Executed when the state begins
	onStart(ecs e.Ecs)
	// Executed when the state exits
	onStop(ecs e.Ecs)
	// Executed when a new state is pushed over this one
	onPause(ecs e.Ecs)
	// Executed when the state become active again (states pushed over this one have been popped)
	onResume(ecs e.Ecs)
	// Executed on every frame when the state is active
	update(ecs e.Ecs, screen *ebiten.Image) transition
}

// StateMachine contains a stack of states.
// Only the top state is active.
type StateMachine struct {
	states []state
}

// Init creates a new state machine with an initial state
func Init(s state, ecs e.Ecs) StateMachine {
	s.onStart(ecs)
	return StateMachine{[]state{s}}
}

// Update updates the state machine
func (sm *StateMachine) Update(ecs e.Ecs, screen *ebiten.Image) {
	if len(sm.states) < 1 {
		os.Exit(0)
	}

	switch t := sm.states[len(sm.states)-1].update(ecs, screen); t.transType {
	case transPop:
		sm._Pop(ecs)
	case transPush:
		sm._Push(ecs, t.newStates)
	case transSwitch:
		sm._Switch(ecs, t.newStates)
	case transReplace:
		sm._Replace(ecs, t.newStates)
	case transQuit:
		sm._Quit(ecs)
	}
}

// Remove the active state and resume the next state
func (sm *StateMachine) _Pop(ecs e.Ecs) {
	sm.states[len(sm.states)-1].onStop(ecs)
	sm.states = sm.states[:len(sm.states)-1]

	if len(sm.states) > 0 {
		sm.states[len(sm.states)-1].onResume(ecs)
	}
}

// Pause the active state and add new states to the stack
func (sm *StateMachine) _Push(ecs e.Ecs, newStates []state) {
	if len(newStates) > 0 {
		sm.states[len(sm.states)-1].onPause(ecs)

		for _, state := range newStates[:len(newStates)-1] {
			state.onStart(ecs)
			state.onPause(ecs)
		}
		newStates[len(newStates)-1].onStart(ecs)

		sm.states = append(sm.states, newStates...)
	}
}

// Remove the active state and replace it by a new one
func (sm *StateMachine) _Switch(ecs e.Ecs, newStates []state) {
	if len(newStates) != 1 {
		utils.LogError(fmt.Errorf("switch transition accept only one new state"))
	}

	sm.states[len(sm.states)-1].onStop(ecs)
	newStates[0].onStart(ecs)
	sm.states[len(sm.states)-1] = newStates[0]
}

// Remove all states and insert a new stack
func (sm *StateMachine) _Replace(ecs e.Ecs, newStates []state) {
	for len(sm.states) > 0 {
		sm.states[len(sm.states)-1].onStop(ecs)
		sm.states = sm.states[:len(sm.states)-1]
	}

	if len(newStates) > 0 {
		for _, state := range newStates[:len(newStates)-1] {
			state.onStart(ecs)
			state.onPause(ecs)
		}
		newStates[len(newStates)-1].onStart(ecs)
	}
	sm.states = newStates
}

// Remove all states and quit
func (sm *StateMachine) _Quit(ecs e.Ecs) {
	for len(sm.states) > 0 {
		sm.states[len(sm.states)-1].onStop(ecs)
		sm.states = sm.states[:len(sm.states)-1]
	}
	os.Exit(0)
}
