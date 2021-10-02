package game

import (
	"syscall/js"
	"time"
)

const (
	skipFrequency = 1
)

type Entity interface {
	Draw(ctx js.Value)
	Update(mouseX int, mouseY int, mpf float64)
	ShouldDraw() bool
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
	game.ctx = game.canvas.Call("getContext", "2d")
	windowScreen := js.Global().Get("window").Get("screen")
	game.windowWidth = int(windowScreen.Get("width").Float())
	game.windowHeight = int(windowScreen.Get("height").Float())
	game.canvas.Set("width", game.windowWidth)
	game.canvas.Set("height", game.windowHeight)
	game.ctx.Set("fillStyle", "white")
	game.ctx.Call("fillRect", 0, 0, game.windowWidth, game.windowHeight)
	time.Sleep(500 * time.Millisecond)
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
	var mpf float64
	frameTime := time.Now()
	frameCount := 1

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if frameCount%skipFrequency == 0 && skipFrequency != 1 && skipFrequency != 0 {
			js.Global().Call("requestAnimationFrame", renderFrame)
			return nil
		}

		mpf = 1 / float64(time.Since(frameTime).Milliseconds())
		frameTime = time.Now()

		shouldDraw := false
		for _, entity := range g.Entities {
			if entity.ShouldDraw() {
				shouldDraw = true
			}
		}

		if shouldDraw {
			g.ctx.Call("clearRect", 0, 0, g.windowWidth, g.windowHeight)
		}

		for _, entity := range g.Entities {
			entity.Update(g.mouseX, g.mouseY, mpf)
			if shouldDraw {
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
