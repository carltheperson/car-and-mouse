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
	innerSpawningOffset = 0
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
}

func NewObstacle(game *game.Game, onObstacleReset func(obs Obstacle)) *Obstacle {
	obs := Obstacle{game: game, OnObstacleReset: onObstacleReset}
	obs.setRandomValues()
	return &obs
}

func (o *Obstacle) Reset() {
	o.setRandomValues()
}

func (o *Obstacle) setRandomValues() {
	rand.Seed(randomSeed)
	randomSeed += time.Now().UnixNano()
	side := rand.Intn(5) + 1
	switch side {
	case 1:
		o.X = rand.Intn(game.CanvasWidth)
		o.Y = 0
	case 2:
		o.X = rand.Intn(game.CanvasWidth)
		o.Y = game.CanvasHeight
	case 3:
		o.X = 0
		o.Y = rand.Intn(game.CanvasHeight)
	case 4:
		o.X = game.CanvasHeight
		o.Y = rand.Intn(game.CanvasHeight)
	}

	// Creating direction by pointing obstacle to random randomPoint inside canvas
	randomPoint := math.Vector2D{A: float64(innerSpawningOffset + rand.Intn(game.CanvasWidth-innerSpawningOffset*2)), B: float64(innerSpawningOffset + rand.Intn(game.CanvasHeight-innerSpawningOffset*2))}
	o.direction = math.Vector2D{A: randomPoint.A - float64(o.X), B: randomPoint.B - float64(o.Y)}

	o.Diameter = minDiameter + rand.Intn(maxDiameter-minDiameter)
	o.speed = minSpeed + rand.Float64()*(maxSpeed-minSpeed)

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
	ctx.Set("lineWidth", 2)
	ctx.Set("strokeStyle", "red")
	ctx.Call("stroke")
}

func (o *Obstacle) Update(mouseX int, mouseY int, mpf float64) {
	if o.spawningDelay > 0 {
		o.spawningDelay -= mpf
		return
	}

	directionUnitVector := o.direction.GetUnitVector()

	o.X += int(directionUnitVector.A * mpf * o.speed)
	o.Y += int(directionUnitVector.B * mpf * o.speed)

	if o.X-(o.Diameter/2) > game.CanvasWidth || o.Y-(o.Diameter/2) > game.CanvasHeight || o.X+(o.Diameter/2) < 0 || o.Y+(o.Diameter/2) < 0 {
		o.OnObstacleReset(*o)
		o.Reset()
	}
}
