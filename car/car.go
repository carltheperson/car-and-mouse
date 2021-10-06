package car

import (
	stdMath "math"
	"syscall/js"

	"github.com/carltheperson/car-and-mouse/math"
)

const (
	maxSpeed = 60
	minSpeed = 40

	maxTurningDif = 0.09
)

type Car struct {
	x         int
	y         int
	lastX     int
	lastY     int
	width     int
	height    int
	direction float64
	speed     float64
}

func NewCar(x int, y int) *Car {

	return &Car{
		x:         x,
		y:         y,
		width:     25,
		height:    50,
		direction: 0.0,
		speed:     maxSpeed / 2,
	}
}

func (c *Car) Draw(ctx js.Value) {
	ctx.Call("beginPath")
	ctx.Set("fillStyle", "#2e2e2e")
	center := math.Vector2D{A: float64(c.x + c.width/2), B: float64(c.y + c.height/2)}

	points := []math.Vector2D{
		{A: float64(c.x), B: float64(c.y)},
		{A: float64(c.x + c.width), B: float64(c.y)},
		{A: float64(c.x + c.width), B: float64(c.y + c.height)},
		{A: float64(c.x + c.width/2), B: float64(c.y + c.height + 10)},
		{A: float64(c.x), B: float64(c.y + c.height)},
	}
	for _, point := range points {
		transformedPoint := math.RotatePoint(point, center, stdMath.Mod(c.direction, 2*stdMath.Pi))
		ctx.Call("lineTo", int(transformedPoint.A), int(transformedPoint.B))
	}

	ctx.Call("closePath")
	ctx.Call("fill")
}

func (c *Car) ShouldDraw() bool {
	return c.lastX != c.x || c.lastY != c.y
}

func (c *Car) Update(mouseX int, mouseY int, mpf float64) {
	c.lastX = c.x
	c.lastY = c.y

	mouseVector := math.Vector2D{A: float64(mouseX - (c.x + c.width/2)), B: float64(mouseY - (c.y + c.height/2))}
	mouseRadians := math.ConvertDirectionVectorToRadians(mouseVector.GetUnitVector())

	directionDifference := math.GetDirectionDifference(c.direction, mouseRadians)
	regulatedDirectedDifference := getRegulatedDirectionDifference(directionDifference, maxTurningDif)

	turningFraction := stdMath.Abs(regulatedDirectedDifference) / maxTurningDif
	if turningFraction > 0.5 {
		c.speed -= (c.speed - minSpeed) * turningFraction * mpf * 0.4
	} else {
		c.speed += (maxSpeed - c.speed) * (1 - turningFraction) * mpf * 0.4
	}

	c.direction = c.direction - regulatedDirectedDifference*(mpf*10)
	directionVector := math.ConvertRadiansToDirectionVector(c.direction)
	directionUnitVector := directionVector.GetUnitVector()

	c.x -= int(directionUnitVector.B * mpf * c.speed)
	c.y += int(directionUnitVector.A * mpf * c.speed)
}

// getRegulatedDirectionDifference will keep the difference to a max size
func getRegulatedDirectionDifference(currentDifference float64, maxDifference float64) float64 {
	if currentDifference < 0 {
		currentDifference = stdMath.Max(currentDifference, -maxDifference)
	} else if currentDifference > 0 {
		currentDifference = stdMath.Min(currentDifference, maxDifference)
	}
	return currentDifference
}
