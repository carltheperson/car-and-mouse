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
			g.Score += 1
			if (g.Score < 5 || g.Score%2 == 0) && len(*car.Obstacles) < 15 {
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
