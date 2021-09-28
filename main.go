package main

import (
	"github.com/carltheperson/car-and-mouse/car"
	"github.com/carltheperson/car-and-mouse/game"
)

func main() {
	g := game.NewGame("canvas")
	g.Entities = append(g.Entities, car.NewCar(100, 100))
	g.RunMainLoop()
}
