package main

import (
	"github.com/carltheperson/car-and-mouse/car"
	"github.com/carltheperson/car-and-mouse/game"
	"github.com/carltheperson/car-and-mouse/obstacle"
)

func main() {
	g := game.NewGame("canvas")
	g.Entities = append(g.Entities, car.NewCar(300, 300, g.State))
	g.Entities = append(g.Entities, obstacle.NewObstacle(800, 800))
	g.RunMainLoop()
}
