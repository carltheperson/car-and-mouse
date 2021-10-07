package obstacle

import (
	stdMath "math"
	"math/rand"
	"syscall/js"
	"time"

	"github.com/carltheperson/car-and-mouse/game"
	"github.com/carltheperson/car-and-mouse/math"
)

var randomSeed = time.Now().UnixNano()

const (
	maxDiameter         = 120
	minDiameter         = 40
	minSpeed            = 30.0
	maxSpeed            = 75.0
	innerSpawningOffset = 50
	maxSpawnDelay       = 25
)

type Obstacle struct {
	X               int
	Y               int
	Diameter        int
	game            *game.Game
	direction       math.Vector2D
	spawningDelay   float64
	OnObstacleReset func(obs Obstacle)
	speed           float64
	canvasWidth     int
	canvasHeight    int
}

func NewObstacle(canvasWidth, canvasHeight int, game *game.Game, onObstacleReset func(obs Obstacle)) *Obstacle {
	obs := Obstacle{game: game, OnObstacleReset: onObstacleReset}
	obs.setRandomValues(canvasWidth, canvasHeight)
	return &obs
}

func (o *Obstacle) Reset() {
	o.setRandomValues(o.canvasWidth, o.canvasHeight)
}

func (o *Obstacle) setRandomValues(canvasWidth, canvasHeight int) {
	rand.Seed(randomSeed)
	randomSeed += time.Now().UnixNano()
	side := rand.Intn(5) + 1
	switch side {
	case 1:
		o.X = rand.Intn(canvasWidth)
		o.Y = 0
	case 2:
		o.X = rand.Intn(canvasWidth)
		o.Y = canvasHeight
	case 3:
		o.X = 0
		o.Y = rand.Intn(canvasHeight)
	case 4:
		o.X = canvasHeight
		o.Y = rand.Intn(canvasHeight)
	}

	// Creating direction by pointing obstacle to random randomPoint inside canvas
	randomPoint := math.Vector2D{A: float64(innerSpawningOffset + rand.Intn(canvasWidth-innerSpawningOffset*2)), B: float64(innerSpawningOffset + rand.Intn(canvasHeight-innerSpawningOffset*2))}
	o.direction = math.Vector2D{A: randomPoint.A - float64(o.X), B: randomPoint.B - float64(o.Y)}

	o.Diameter = minDiameter + rand.Intn(maxDiameter-minDiameter)
	o.speed = minSpeed + rand.Float64()*(maxSpeed-minSpeed)
	o.canvasWidth = canvasWidth
	o.canvasHeight = canvasHeight

	o.spawningDelay = float64(rand.Intn(maxSpawnDelay))
}

func (o *Obstacle) Draw(ctx js.Value) {
	if o.spawningDelay > 0 {
		return
	}
	ctx.Call("beginPath")
	ctx.Call("arc", o.X, o.Y, o.Diameter/2, 0, 2*stdMath.Pi, false)
	ctx.Set("fillStyle", "salmon")
	ctx.Call("fill")
}

func (o *Obstacle) Update(mouseX int, mouseY int, mpf float64) {
	if o.spawningDelay > 0 {
		o.spawningDelay -= mpf
		return
	}

	directionUnitVector := o.direction.GetUnitVector()

	o.X += int(directionUnitVector.A * mpf * o.speed)
	o.Y += int(directionUnitVector.B * mpf * o.speed)

	if o.X-o.Diameter/2 > o.canvasWidth || o.Y-o.Diameter/2 > o.canvasHeight || o.X+o.Diameter/2 < 0 || o.Y+o.Diameter/2 < 0 {
		o.OnObstacleReset(*o)
		o.Reset()
	}
}

func (o *Obstacle) ShouldDraw() bool {
	return true
}
