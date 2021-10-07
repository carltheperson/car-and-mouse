package game

import (
	"fmt"
	"syscall/js"
	"time"
)

const (
	CanvasWidth  = 800
	CanvasHeight = 800

	skipFrequency = 1
)

const (
	StateNormal = iota
	StateGameOver
)

type Entity interface {
	Draw(ctx js.Value)
	Update(mouseX int, mouseY int, mpf float64)
	ShouldDraw() bool
}

type Game struct {
	body                   js.Value
	document               js.Value
	canvas                 js.Value
	ctx                    js.Value
	prompt                 js.Value
	shouldReRenderPrompt   bool
	WindowWidth            int
	WindowHeight           int
	mouseX                 int
	mouseY                 int
	Entities               *[]Entity
	State                  int
	lastState              int
	Points                 int
	lastPoints             int
	addInitialEntitiesFunc func()
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
	game.WindowWidth = int(windowScreen.Get("width").Float())
	game.WindowHeight = int(windowScreen.Get("height").Float())
	game.canvas.Set("width", CanvasWidth)
	game.canvas.Set("height", CanvasHeight)
	game.ctx.Set("fillStyle", "white")
	game.ctx.Call("fillRect", 0, 0, game.WindowWidth, game.WindowHeight)
	game.Points = 0
	game.lastPoints = -1
	game.State = StateNormal
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

func (g *Game) SetAddInitialEntitiesFunc(function func()) {
	g.addInitialEntitiesFunc = function
}

func (g *Game) restartGame() {
	*g.Entities = []Entity{}
	if g.addInitialEntitiesFunc == nil {
		panic("No function was set to initialize entities. Call setAddInitialEntitiesFunc to correct this")
	}
	g.addInitialEntitiesFunc()
	g.Points = 0
	g.State = StateNormal
}

func (g *Game) setTryAgainButtonEventListener() {
	eventListener := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.restartGame()
		return nil
	})
	g.document.Call("getElementById", "try-again-button").Call("addEventListener", "click", eventListener)
}

func (g *Game) showPromptForState() {
	if g.State != g.lastState {
		g.shouldReRenderPrompt = true
	}
	if !g.shouldReRenderPrompt {
		g.lastState = g.State
		return
	}

	switch g.State {
	case StateNormal:
		if g.lastPoints != g.Points {
			g.prompt.Set("innerHTML", "Points "+fmt.Sprint(g.Points))
		}
		g.lastPoints = g.Points
		g.shouldReRenderPrompt = true
	case StateGameOver:
		html := fmt.Sprintf(`GAME OVER! Points %d
		<br><br>
		<button id="try-again-button">Try again</button>`, g.Points)
		g.prompt.Set("innerHTML", html)
		g.setTryAgainButtonEventListener()
		g.shouldReRenderPrompt = false
	}
	g.lastState = g.State

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
		g.showPromptForState()

		if g.State == StateGameOver {
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
			g.ctx.Call("clearRect", 0, 0, g.WindowWidth, g.WindowHeight)
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
