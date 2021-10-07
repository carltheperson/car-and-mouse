package obstacle

import (
	"fmt"
	stdMath "math"
	"math/rand"
	"syscall/js"
	"time"

	"github.com/carltheperson/car-and-mouse/math"
)

const (
	maxDiameter = 120
	minDiameter = 40
	minSpeed    = 30.0
	maxSpeed    = 50.0
)

type Obstacle struct {
	x            int
	y            int
	diameter     int
	direction    math.Vector2D
	speed        float64
	canvasWidth  int
	canvasHeight int
}

func NewObstacle(canvasWidth, canvasHeight int) *Obstacle {
	obs := Obstacle{}
	obs.setRandomValues(canvasWidth, canvasHeight)
	return &obs
}

func (o *Obstacle) Reset() {
	o.setRandomValues(o.canvasWidth, o.canvasHeight)
}

func (o *Obstacle) setRandomValues(canvasWidth, canvasHeight int) {
	rand.Seed(time.Now().UnixNano())
	side := rand.Intn(5) + 1
	switch side {
	case 1:
		o.x = rand.Intn(canvasWidth)
		o.y = 0
	case 2:
		o.x = rand.Intn(canvasWidth)
		o.y = canvasHeight
	case 3:
		o.x = 0
		o.y = rand.Intn(canvasHeight)
	case 4:
		o.x = canvasHeight
		o.y = rand.Intn(canvasHeight)
	}

	// Creating direction by pointing obstacle to random randomPoint inside canvas
	randomPoint := math.Vector2D{A: float64(rand.Intn(canvasWidth)), B: float64(rand.Intn(canvasHeight))}
	o.direction = math.Vector2D{A: randomPoint.A - float64(o.x), B: randomPoint.B - float64(o.y)}

	o.diameter = minDiameter + rand.Intn(maxDiameter-minDiameter)
	o.speed = minSpeed + rand.Float64()*(maxSpeed-minSpeed)

	fmt.Println(randomPoint)
	fmt.Println(o.direction)

}

func (o *Obstacle) Draw(ctx js.Value) {
	ctx.Call("beginPath")
	ctx.Call("arc", o.x, o.y, o.diameter/2, 0, 2*stdMath.Pi, false)
	ctx.Set("fillStyle", "red")
	ctx.Call("fill")
}

func (o *Obstacle) Update(mouseX int, mouseY int, mpf float64) {
	directionUnitVector := o.direction.GetUnitVector()

	o.x += int(directionUnitVector.A * mpf * o.speed)
	o.y += int(directionUnitVector.B * mpf * o.speed)
}

func (o *Obstacle) ShouldDraw() bool {
	return true
}
