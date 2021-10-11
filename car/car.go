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

	spritePath = "/assets/car_sprite2.png"
)

type Car struct {
	x         int
	y         int
	lastX     int
	lastY     int
	width     int
	height    int
	sprite    js.Value
	direction float64
	speed     float64
	Obstacles *[]*obstacle.Obstacle
	game      *game.Game
}

func NewCar(x int, y int, game *game.Game, obstacles *[]*obstacle.Obstacle) *Car {
	sprite := js.Global().Get("Image").New()
	sprite.Set("src", spritePath)
	return &Car{
		x:         x,
		y:         y,
		width:     sprite.Get("width").Int(),
		height:    sprite.Get("height").Int(),
		direction: 0.0,
		speed:     minSpeed / 2,
		Obstacles: obstacles,
		game:      game,
		sprite:    sprite,
	}
}

func (c *Car) Draw(ctx js.Value) {
	ctx.Call("beginPath")
	ctx.Set("fillStyle", "#2e2e2e")
	center := c.getCenter()

	ctx.Call("setTransform", 1, 0, 0, 1, center.A, center.B)
	ctx.Call("rotate", stdMath.Mod(c.direction, 2*stdMath.Pi)+stdMath.Pi)
	ctx.Call("drawImage", c.sprite, -c.width/2, -c.height/2)
	ctx.Call("setTransform", 1, 0, 0, 1, 0, 0)
	ctx.Call("closePath")
	ctx.Call("fill")
}

func (c *Car) getCenter() math.Vector2D {
	return math.Vector2D{A: float64(c.x + c.width/2), B: float64(c.y + c.height/2)}
}

func (c *Car) IsOutsideCanvas() bool {
	center := c.getCenter()
	x := center.A
	y := center.B
	if x > float64(game.CanvasWidth) || x < 0 || y > float64(game.CanvasHeight) || y < 0 {
		return true
	}
	return false
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
	for _, o := range *c.Obstacles {
		if int(math.GetDistanceBetweenTwoPoints(c.getCenter(), math.Vector2D{A: float64(o.X), B: float64(o.Y)})) < o.Diameter/2+c.width/2 {
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

	if c.IsTouchingMouse(mouseX, mouseY) || c.IsTouchingObstacle() || c.IsOutsideCanvas() {
		c.game.State = game.StateGameOver
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
