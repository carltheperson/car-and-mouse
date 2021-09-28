package game

import (
	"fmt"
	"syscall/js"
)

type Game struct {
	body         js.Value
	document     js.Value
	canvas       js.Value
	ctx          js.Value
	windowWidth  float64
	windowHeight float64
	mouseX       float64
	mouseY       float64
}

func NewGame(canvasId string) Game {
	game := Game{}
	game.document = js.Global().Get("document")
	game.body = game.document.Get("body")
	game.canvas = game.document.Call("getElementById", canvasId)
	game.ctx = game.canvas.Call("getContext", "2d")
	game.windowWidth = game.body.Get("clientWidth").Float()
	game.windowHeight = game.body.Get("clientHeight").Float()
	return game
}

func (g *Game) getMouseMoveEventListener() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		g.mouseX = event.Get("clientX").Float()
		g.mouseY = event.Get("clientY").Float()
		fmt.Println(g.mouseX)
		return nil
	})
}

func (g *Game) RunMainLoop() {
	c := make(chan bool)

	mouseMoveEventListener := g.getMouseMoveEventListener()
	defer mouseMoveEventListener.Release()
	g.document.Call("addEventListener", "mousemove", mouseMoveEventListener)

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
