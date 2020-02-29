package states

import (
	"fmt"
	"os"

	"arkanoid/lib/ecs"
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
	onStart(world ecs.World)
	// Executed when the state exits
	onStop(world ecs.World)
	// Executed when a new state is pushed over this one
	onPause(world ecs.World)
	// Executed when the state become active again (states pushed over this one have been popped)
	onResume(world ecs.World)
	// Executed on every frame when the state is active
	update(world ecs.World, screen *ebiten.Image) transition
}

// StateMachine contains a stack of states.
// Only the top state is active.
type StateMachine struct {
	states []state
}

// Init creates a new state machine with an initial state
func Init(s state, world ecs.World) StateMachine {
	s.onStart(world)
	return StateMachine{[]state{s}}
}

// Update updates the state machine
func (sm *StateMachine) Update(world ecs.World, screen *ebiten.Image) {
	if len(sm.states) < 1 {
		os.Exit(0)
	}

	switch t := sm.states[len(sm.states)-1].update(world, screen); t.transType {
	case transPop:
		sm._Pop(world)
	case transPush:
		sm._Push(world, t.newStates)
	case transSwitch:
		sm._Switch(world, t.newStates)
	case transReplace:
		sm._Replace(world, t.newStates)
	case transQuit:
		sm._Quit(world)
	}
}

// Remove the active state and resume the next state
func (sm *StateMachine) _Pop(world ecs.World) {
	sm.states[len(sm.states)-1].onStop(world)
	sm.states = sm.states[:len(sm.states)-1]

	if len(sm.states) > 0 {
		sm.states[len(sm.states)-1].onResume(world)
	}
}

// Pause the active state and add new states to the stack
func (sm *StateMachine) _Push(world ecs.World, newStates []state) {
	if len(newStates) > 0 {
		sm.states[len(sm.states)-1].onPause(world)

		for _, state := range newStates[:len(newStates)-1] {
			state.onStart(world)
			state.onPause(world)
		}
		newStates[len(newStates)-1].onStart(world)

		sm.states = append(sm.states, newStates...)
	}
}

// Remove the active state and replace it by a new one
func (sm *StateMachine) _Switch(world ecs.World, newStates []state) {
	if len(newStates) != 1 {
		utils.LogError(fmt.Errorf("switch transition accept only one new state"))
	}

	sm.states[len(sm.states)-1].onStop(world)
	newStates[0].onStart(world)
	sm.states[len(sm.states)-1] = newStates[0]
}

// Remove all states and insert a new stack
func (sm *StateMachine) _Replace(world ecs.World, newStates []state) {
	for len(sm.states) > 0 {
		sm.states[len(sm.states)-1].onStop(world)
		sm.states = sm.states[:len(sm.states)-1]
	}

	if len(newStates) > 0 {
		for _, state := range newStates[:len(newStates)-1] {
			state.onStart(world)
			state.onPause(world)
		}
		newStates[len(newStates)-1].onStart(world)
	}
	sm.states = newStates
}

// Remove all states and quit
func (sm *StateMachine) _Quit(world ecs.World) {
	for len(sm.states) > 0 {
		sm.states[len(sm.states)-1].onStop(world)
		sm.states = sm.states[:len(sm.states)-1]
	}
	os.Exit(0)
}
