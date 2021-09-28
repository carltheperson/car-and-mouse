package game

import (
	"fmt"
	"syscall/js"
)

type Game struct {
	body     js.Value
	document js.Value
	canvas   js.Value
}

func NewGame(canvasId string) Game {
	game := Game{}
	game.document = js.Global().Get("document")
	game.body = game.document.Get("body")
	game.canvas = game.document.Call("getElementById", canvasId)
	return game
}

func (G *Game) RunMainLoop() {
	c := make(chan bool)

	var renderFrame js.Func

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Hello")
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})

	defer renderFrame.Release()
	js.Global().Call("requestAnimationFrame", renderFrame)

	<-c
}
