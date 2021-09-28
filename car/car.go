package car

import (
	"syscall/js"
)

const pxsPerMpf = 100

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
	ctx.Call("beginPath")
	ctx.Call("rect", c.x, c.y, c.width, c.height)
	ctx.Call("stroke")
}

func (c *Car) Update(mouseX int, mouseY int, mpf float64) {
	c.lastX = c.x
	c.lastY = c.y
	if mouseX > c.x+c.width/2 {
		c.x += int(pxsPerMpf * mpf)
	}
	if mouseX < c.x+c.width/2 {
		c.x -= int(pxsPerMpf * mpf)
	}
	if mouseY > c.y+c.height/2 {
		c.y += int(pxsPerMpf * mpf)
	}
	if mouseY < c.y+c.height/2 {
		c.y -= int(pxsPerMpf * mpf)
	}
}

func (c *Car) GetShouldDraw() bool {
	return c.lastX != c.x || c.lastY != c.y
}
