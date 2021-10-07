package car

import (
	stdMath "math"
	"syscall/js"

	"github.com/carltheperson/car-and-mouse/game"
	"github.com/carltheperson/car-and-mouse/math"
	"github.com/carltheperson/car-and-mouse/obstacle"
)

const (
	maxSpeed = 60
	minSpeed = 40

	maxTurningDif = 0.075
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
	obstacles []*obstacle.Obstacle
	game      *game.Game
}

func NewCar(x int, y int, game *game.Game, obstacles []*obstacle.Obstacle) *Car {

	return &Car{
		x:         x,
		y:         y,
		width:     25,
		height:    50,
		direction: 0.0,
		speed:     maxSpeed / 2,
		obstacles: obstacles,
		game:      game,
	}
}

func (c *Car) Draw(ctx js.Value) {
	ctx.Call("beginPath")
	ctx.Set("fillStyle", "#2e2e2e")
	center := c.getCenter()

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

func (c *Car) getCenter() math.Vector2D {
	return math.Vector2D{A: float64(c.x + c.width/2), B: float64(c.y + c.height/2)}
}

func (c *Car) IsTouchingMouse(mouseX, mouseY int) bool {
	center := c.getCenter()
	mouse := math.Vector2D{A: float64(mouseX), B: float64(mouseY)}
	transformedMousePoint := math.RotatePoint(mouse, center, stdMath.Mod(c.direction, 2*stdMath.Pi))
	mX := int(transformedMousePoint.A)
	mY := int(transformedMousePoint.B)
	if mX > c.x && mX < c.x+c.width && mY > c.y && mY < c.y+c.height {
		return true
	}
	return false
}

func (c *Car) IsTouchingObstacle() bool {
	for _, o := range c.obstacles {
		if int(math.GetDistanceBetweenTwoPoints(c.getCenter(), math.Vector2D{A: float64(o.X), B: float64(o.Y)})) < o.Diameter/2+c.height/3 {
			return true
		}
	}
	return false
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
		c.speed -= (c.speed - minSpeed) * turningFraction * mpf * 0.25
	} else {
		c.speed += (maxSpeed - c.speed) * (1 - turningFraction) * mpf * 0.25
	}

	if c.IsTouchingMouse(mouseX, mouseY) || c.IsTouchingObstacle() {
		*c.game.State = game.StateGameOver
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
