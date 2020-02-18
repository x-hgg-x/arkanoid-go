package main

import (
	"github.com/ByteArena/ecs"
	"github.com/hajimehoshi/ebiten"
)

const (
	windowWidth  = 720
	windowHeight = 600
	scaleFactor  = 1
)

func initEcs() *ecs.Manager {
	manager := ecs.NewManager()
	return manager
}

func update(screen *ebiten.Image) error {
	return nil
}

func main() {
	if err := ebiten.Run(update, windowWidth, windowHeight, scaleFactor, "Arkanoid"); err != nil {
		panic(err)
	}
}
