package game

import (
	"fmt"
	"syscall/js"
	"time"
)

const (
	skipFrequency = 1
)

const (
	StateNormal = iota
	StateRestarting
	StateGameOver
	StateWon
)

type Entity interface {
	Draw(ctx js.Value)
	Update(mouseX int, mouseY int, mpf float64)
	ShouldDraw() bool
}

type Game struct {
	body                 js.Value
	document             js.Value
	canvas               js.Value
	ctx                  js.Value
	prompt               js.Value
	shouldReRenderPrompt bool
	windowWidth          int
	windowHeight         int
	mouseX               int
	mouseY               int
	Entities             *[]Entity
	State                *int
	lastState            int
	Points               int
	lastPoints           int
}

func NewGame(canvasId string) Game {
	game := Game{}
	game.document = js.Global().Get("document")
	game.body = game.document.Get("body")
	game.canvas = game.document.Call("getElementById", canvasId)
	game.prompt = game.document.Call("getElementById", "prompt")
	game.shouldReRenderPrompt = true
	game.ctx = game.canvas.Call("getContext", "2d")
	windowScreen := js.Global().Get("window").Get("screen")
	game.windowWidth = int(windowScreen.Get("width").Float())
	game.windowHeight = int(windowScreen.Get("height").Float())
	game.canvas.Set("width", 800)
	game.canvas.Set("height", 800)
	game.ctx.Set("fillStyle", "white")
	game.ctx.Call("fillRect", 0, 0, game.windowWidth, game.windowHeight)
	game.Points = 0
	game.lastPoints = -1
	game.State = new(int)
	*game.State = StateNormal
	game.Entities = &[]Entity{}
	time.Sleep(500 * time.Millisecond)
	return game
}

func (g *Game) getCanvasXAndY() (float64, float64) {
	x := g.canvas.Call("getBoundingClientRect").Get("left").Float()
	y := g.canvas.Call("getBoundingClientRect").Get("top").Float()
	return x, y
}

func (g *Game) getMouseMoveEventListener() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		canvasX, canvasY := g.getCanvasXAndY()
		g.mouseX = int(event.Get("clientX").Float()) - int(canvasX)
		g.mouseY = int(event.Get("clientY").Float()) - int(canvasY)
		return nil
	})
}

func (g *Game) setTextInPrompt(text string) {
	g.prompt.Set("innerHTML", text)
}

func (g *Game) showPromptForState(state int) {
	if state != g.lastState {
		g.shouldReRenderPrompt = true
	}
	if !g.shouldReRenderPrompt {
		g.lastState = state
		return
	}

	g.lastState = state

	switch state {
	case StateNormal:
		if g.lastPoints != g.Points {
			g.setTextInPrompt("Points " + fmt.Sprint(g.Points))
		}
		g.lastPoints = g.Points
		g.shouldReRenderPrompt = true
	case StateRestarting:
		g.setTextInPrompt("We are restarting")
		g.shouldReRenderPrompt = false
	case StateGameOver:
		g.setTextInPrompt("GAME OVER! Points " + fmt.Sprint(g.Points))
		g.shouldReRenderPrompt = false

	}

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
		g.showPromptForState(*g.State)

		if *g.State == StateGameOver {
			js.Global().Call("requestAnimationFrame", renderFrame)
			return nil
		}

		if frameCount%skipFrequency == 0 && skipFrequency != 1 && skipFrequency != 0 {
			js.Global().Call("requestAnimationFrame", renderFrame)
			return nil
		}

		mpf = 1 / float64(time.Since(frameTime).Milliseconds())
		frameTime = time.Now()

		shouldDraw := false
		for _, entity := range *g.Entities {
			if entity.ShouldDraw() {
				shouldDraw = true
			}
		}

		if shouldDraw {
			g.ctx.Call("clearRect", 0, 0, g.windowWidth, g.windowHeight)
		}

		for _, entity := range *g.Entities {
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
