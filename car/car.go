package car

import (
	"syscall/js"

	"github.com/carltheperson/car-and-mouse/math"
)

const pxsPerMpf = 100

type Car struct {
	x         int
	y         int
	lastX     int
	lastY     int
	width     int
	height    int
	direction math.Vector2D
}

func NewCar(x int, y int) *Car {

	return &Car{
		x:         x,
		y:         y,
		width:     50,
		height:    100,
		direction: math.Vector2D{},
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
		transformedPoint := math.RotatePoint(point, center, math.ConvertDirectionVectorToRadians(c.direction))
		ctx.Call("lineTo", int(transformedPoint.A), int(transformedPoint.B))
	}

	ctx.Call("closePath")
	ctx.Call("fill")
}

func (c *Car) Update(mouseX int, mouseY int, mpf float64) {
	c.lastX = c.x
	c.lastY = c.y

	mouseV := math.Vector2D{A: float64(mouseX - (c.x + c.width/2)), B: float64(mouseY - (c.y + c.height/2))}
	c.direction = mouseV.GetUnitVector()

	c.x += int(c.direction.A * mpf * pxsPerMpf)
	c.y += int(c.direction.B * mpf * pxsPerMpf)
}

func (c *Car) ShouldDraw() bool {
	return c.lastX != c.x || c.lastY != c.y
}
