package car

import (
	stdMath "math"
	"sort"
	"syscall/js"

	"github.com/carltheperson/car-and-mouse/math"
)

const pxsPerMpf = 60

const maxDif = 0.09

type Car struct {
	x          int
	y          int
	lastX      int
	lastY      int
	width      int
	height     int
	direction  float64
	correcting bool
}

func NewCar(x int, y int) *Car {

	return &Car{
		x:          x,
		y:          y,
		width:      25,
		height:     50,
		direction:  0.0,
		correcting: true,
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

func (c *Car) Update(mouseX int, mouseY int, mpf float64) {
	c.lastX = c.x
	c.lastY = c.y

	mouseVector := math.Vector2D{A: float64(mouseX - (c.x + c.width/2)), B: float64(mouseY - (c.y + c.height/2))}
	mouseRadians := math.ConvertDirectionVectorToRadians(mouseVector.GetUnitVector())

	directionDifference := getDirectionDifference(c.direction, mouseRadians)

	if directionDifference < 0 {
		directionDifference = stdMath.Max(directionDifference, -maxDif)
	} else if directionDifference > 0 {
		directionDifference = stdMath.Min(directionDifference, maxDif)
	}

	c.direction = c.direction - directionDifference*(mpf*10)
	directionVector := math.ConvertRadiansToDirectionVector(c.direction)
	directionUnitVector := directionVector.GetUnitVector()

	c.x -= int(directionUnitVector.B * mpf * pxsPerMpf)
	c.y += int(directionUnitVector.A * mpf * pxsPerMpf)
}

func getDirectionDifference(d1 float64, d2 float64) float64 {
	frac1 := stdMath.Mod(d1, stdMath.Pi*2)
	frac2 := stdMath.Mod(d2, stdMath.Pi*2)

	option1 := frac1 - frac2
	option2 := ((stdMath.Pi*2 - frac1) + frac2) * -1
	option3 := (stdMath.Pi*2 - frac2) + frac1
	options := []float64{stdMath.Abs(option1), stdMath.Abs(option2), stdMath.Abs(option3)}
	sort.Float64s(options)
	smallestOption := options[0]

	switch smallestOption {
	case stdMath.Abs(option1):
		return option1
	case stdMath.Abs(option2):
		return option2
	}
	return option3

}

func (c *Car) ShouldDraw() bool {
	return c.lastX != c.x || c.lastY != c.y
}
