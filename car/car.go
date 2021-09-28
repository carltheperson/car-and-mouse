package car

import (
	"fmt"
	"syscall/js"
)

type Car struct {
}

func (c *Car) Draw(ctx js.Value) {

}

func (c *Car) Update() {
	fmt.Println("Updating car")
}
