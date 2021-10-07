package main

import (
	"github.com/carltheperson/car-and-mouse/car"
	"github.com/carltheperson/car-and-mouse/game"
	"github.com/carltheperson/car-and-mouse/obstacle"
)

func main() {
	g := game.NewGame("canvas")

	addInitialEntities := func() {
		obstacles := []*obstacle.Obstacle{}

		car := car.NewCar(400, 400, &g, &obstacles)

		onObstacleReset := func(obs obstacle.Obstacle) {
			g.Points += 1
			if (g.Points < 5 || g.Points%3 == 0 || g.Points%3 == 1) && len(*car.Obstacles) < 15 {
				newObstacle := obstacle.NewObstacle(&g, obs.OnObstacleReset)
				*g.Entities = append(*g.Entities, newObstacle)
				*car.Obstacles = append(*car.Obstacles, newObstacle)
			}
		}
		firstObstacle := obstacle.NewObstacle(&g, onObstacleReset)
		*g.Entities = append(*g.Entities, firstObstacle)
		*g.Entities = append(*g.Entities, car)
		*car.Obstacles = append(*car.Obstacles, firstObstacle)
	}
	addInitialEntities()

	g.SetAddInitialEntitiesFunc(addInitialEntities)

	g.RunMainLoop()
}
