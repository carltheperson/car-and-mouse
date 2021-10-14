package car

import (
	stdMath "math"

	"github.com/carltheperson/car-and-mouse/game"
	"github.com/carltheperson/car-and-mouse/math"
)

func (c *Car) getCenter() math.Vector {
	return math.Vector{A: float64(c.x + c.width/2), B: float64(c.y + c.height/2)}
}

func (c *Car) IsOutsideCanvas() bool {
	for _, corner := range c.getTransformedCorners() {
		if corner.A > float64(game.CanvasWidth)+allowedCollisionOverlap || corner.A < 0-allowedCollisionOverlap || corner.B > game.CanvasHeight+allowedCollisionOverlap || corner.B < 0-allowedCollisionOverlap {
			return true
		}
	}
	return false
}

func (c *Car) IsTouchingMouse(mouseX, mouseY int) bool {
	center := c.getCenter()
	mouse := math.Vector{A: float64(mouseX), B: float64(mouseY)}
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
		obs := math.Vector{A: float64(o.X), B: float64(o.Y)}
		if int(math.GetDistanceBetweenTwoPoints(c.getCenter(), obs))+allowedCollisionOverlap/2 < o.Diameter/2+c.width/2 {
			return true
		}
		for _, corner := range c.getTransformedCorners() {
			if int(math.GetDistanceBetweenTwoPoints(corner, obs))+allowedCollisionOverlap/2 < o.Diameter/2 {
				return true
			}
		}
	}
	return false
}

func (c *Car) getTransformedCorners() []math.Vector {
	center := c.getCenter()
	points := []math.Vector{
		{A: float64(c.x), B: float64(c.y)},
		{A: float64(c.x + c.width), B: float64(c.y)},
		{A: float64(c.x + c.width), B: float64(c.y + c.height)},
		{A: float64(c.x), B: float64(c.y + c.height)},
	}

	transformedPoints := []math.Vector{}
	for _, point := range points {
		transformedPoints = append(transformedPoints, math.RotatePoint(point, center, stdMath.Mod(c.direction, 2*stdMath.Pi)))
	}

	return transformedPoints
}
