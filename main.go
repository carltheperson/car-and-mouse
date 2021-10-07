package main

import (
	"github.com/carltheperson/car-and-mouse/car"
	"github.com/carltheperson/car-and-mouse/game"
	"github.com/carltheperson/car-and-mouse/obstacle"
)

func main() {
	g := game.NewGame("canvas")

	obstacles := []*obstacle.Obstacle{}

	car := car.NewCar(300, 300, &g, &obstacles)
	*g.Entities = append(*g.Entities, car)

	onObstacleReset := func(obs obstacle.Obstacle) {
		g.Points += 1
		if g.Points < 5 || g.Points%2 == 0 {
			newObstacle := obstacle.NewObstacle(800, 800, &g, obs.OnObstacleReset)
			*g.Entities = append(*g.Entities, newObstacle)
			*car.Obstacles = append(*car.Obstacles, newObstacle)
		}
	}

	firstObstacle := obstacle.NewObstacle(800, 800, &g, onObstacleReset)
	*g.Entities = append(*g.Entities, firstObstacle)
	*car.Obstacles = append(*car.Obstacles, firstObstacle)

	g.RunMainLoop()
}
