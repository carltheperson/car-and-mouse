package main

import (
	"github.com/carltheperson/car-and-mouse/car"
	"github.com/carltheperson/car-and-mouse/game"
	"github.com/carltheperson/car-and-mouse/obstacle"
)

func main() {
	g := game.NewGame("canvas")
	obstacles := []*obstacle.Obstacle{obstacle.NewObstacle(800, 800), obstacle.NewObstacle(800, 800)}
	car := car.NewCar(300, 300, &g, obstacles)

	g.Entities = append(g.Entities, car)
	g.Entities = append(g.Entities, obstacles[0])
	g.Entities = append(g.Entities, obstacles[1])

	g.RunMainLoop()
}
