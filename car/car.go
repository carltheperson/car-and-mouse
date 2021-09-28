package car

import (
	"fmt"
	"syscall/js"
)

type Car struct {
}

func (c *Car) Draw(ctx js.Value) {

}

func (c *Car) Update(mouseX int, mouseY int) {
	fmt.Println("Updating car")
}

func (c *Car) GetShouldDraw() bool {
	return false
}
