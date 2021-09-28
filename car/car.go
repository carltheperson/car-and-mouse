package car

import (
	"fmt"
	"syscall/js"
)

type Car struct {
	x      int
	y      int
	lastX  int
	lastY  int
	width  int
	height int
}

func NewCar(x int, y int) *Car {
	return &Car{
		x:      x,
		y:      y,
		width:  50,
		height: 100,
	}
}

func (c *Car) Draw(ctx js.Value) {
	fmt.Println(c.x, c.y)
	ctx.Call("beginPath")
	ctx.Call("rect", c.x, c.y, c.width, c.height)
	ctx.Call("stroke")
}

func (c *Car) Update(mouseX int, mouseY int) {
	c.lastX = c.x
	c.lastY = c.y
	if mouseX > c.x {
		c.x += 2
	}
	if mouseX < c.x {
		c.x -= 2
	}
	if mouseY > c.y {
		c.y += 2
	}
	if mouseY < c.y {
		c.y -= 2
	}
}

func (c *Car) GetShouldDraw() bool {
	return c.lastX != c.x || c.lastY != c.y
}
