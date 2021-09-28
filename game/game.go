package game

import (
	"syscall/js"
)

type Entity interface {
	Draw(ctx js.Value)
	Update(mouseX int, mouseY int)
	GetShouldDraw() bool
}

type Game struct {
	body         js.Value
	document     js.Value
	canvas       js.Value
	ctx          js.Value
	windowWidth  int
	windowHeight int
	mouseX       int
	mouseY       int
	Entities     []Entity
}

func NewGame(canvasId string) Game {
	game := Game{}
	game.document = js.Global().Get("document")
	game.body = game.document.Get("body")
	game.canvas = game.document.Call("getElementById", canvasId)
	game.ctx = game.canvas.Call("getContext", "2d", map[string]interface{}{"alpha": false})
	windowScreen := js.Global().Get("window").Get("screen")
	game.windowWidth = int(windowScreen.Get("width").Float())
	game.windowHeight = int(windowScreen.Get("height").Float())
	game.canvas.Set("width", game.windowWidth)
	game.canvas.Set("height", game.windowHeight)
	game.ctx.Set("fillStyle", "white")
	game.ctx.Call("fillRect", 0, 0, game.windowWidth, game.windowHeight)
	return game
}

func (g *Game) getMouseMoveEventListener() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		g.mouseX = int(event.Get("clientX").Float())
		g.mouseY = int(event.Get("clientY").Float())
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
		// g.ctx.Call("beginPath")
		// g.ctx.Call("rect", g.mouseX-20, g.mouseY-20, 20, 20)
		// g.ctx.Call("stroke")

		for _, entity := range g.Entities {
			entity.Update(g.mouseX, g.mouseY)
			if entity.GetShouldDraw() {
				entity.Draw(g.ctx)
			}
		}

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})

	defer renderFrame.Release()
	js.Global().Call("requestAnimationFrame", renderFrame)

	<-c
}
