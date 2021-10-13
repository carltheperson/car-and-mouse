package game

import (
	"fmt"
	"syscall/js"
	"time"
)

const (
	CanvasWidth  = 800
	CanvasHeight = 800

	skipFrequency            = 1
	highscoreLocalStorageKey = "highstore"
)

const (
	StateNormal = iota
	StateGameOver
)

type Entity interface {
	Draw(ctx js.Value)
	Update(mouseX int, mouseY int, mpf float64)
}

type Game struct {
	body                   js.Value
	document               js.Value
	window                 js.Value
	canvas                 js.Value
	ctx                    js.Value
	promptElement          js.Value
	scoreElement           js.Value
	highscoreElement       js.Value
	highscore              int
	shouldReRenderPrompt   bool
	WindowWidth            int
	WindowHeight           int
	mouseX                 int
	mouseY                 int
	Entities               *[]Entity
	State                  int
	lastState              int
	Score                  int
	lastScore              int
	addInitialEntitiesFunc func()
}

func NewGame(canvasId string) Game {
	game := Game{}
	game.window = js.Global().Get("window")
	game.document = js.Global().Get("document")
	game.body = game.document.Get("body")
	game.canvas = game.document.Call("getElementById", canvasId)
	game.promptElement = game.document.Call("getElementById", "prompt")
	game.scoreElement = game.document.Call("getElementById", "score")
	game.highscoreElement = game.document.Call("getElementById", "highscore")
	game.highscore = game.getHighscore()
	game.shouldReRenderPrompt = true
	game.ctx = game.canvas.Call("getContext", "2d")
	windowScreen := js.Global().Get("window").Get("screen")
	game.WindowWidth = int(windowScreen.Get("width").Float())
	game.WindowHeight = int(windowScreen.Get("height").Float())
	game.canvas.Set("width", CanvasWidth)
	game.canvas.Set("height", CanvasHeight)
	game.ctx.Set("fillStyle", "white")
	game.ctx.Call("fillRect", 0, 0, game.WindowWidth, game.WindowHeight)
	game.Score = 0
	game.lastScore = -1
	game.State = StateNormal
	game.Entities = &[]Entity{}
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

func (g *Game) renderGameState() {
	if g.State != g.lastState {
		g.shouldReRenderPrompt = true
	}
	if !g.shouldReRenderPrompt {
		g.lastState = g.State
		return
	}

	switch g.State {
	case StateNormal:
		if g.lastScore != g.Score {
			if g.Score > g.highscore {
				g.highscore = g.Score
				g.setHighscoreInLocalStorage(g.highscore)
			}

			g.scoreElement.Set("innerText", "Score "+fmt.Sprint(g.Score))
			g.highscoreElement.Set("innerText", "Highscore "+fmt.Sprint(g.highscore))
			g.promptElement.Set("innerHTML", "")
		}
		g.lastScore = g.Score
		g.shouldReRenderPrompt = true

	case StateGameOver:
		html := "<h1 style=\"margin: 3px\">GAME OVER!</h1><br><button id=\"try-again-button\">Try again</button>"
		g.promptElement.Set("innerHTML", html)
		g.setTryAgainButtonEventListener()
		g.shouldReRenderPrompt = false
	}
	g.lastState = g.State
}

func (g *Game) RunMainLoop() {
	c := make(chan bool)

	mouseMoveEventListener := g.getMouseMoveEventListener()
	g.document.Call("addEventListener", "mousemove", mouseMoveEventListener)

	var renderFrame js.Func
	var mpf float64
	frameTime := time.Now()
	frameCount := 1

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.renderGameState()

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

		g.ctx.Call("clearRect", 0, 0, g.WindowWidth, g.WindowHeight)

		for _, entity := range *g.Entities {
			entity.Update(g.mouseX, g.mouseY, mpf)
			entity.Draw(g.ctx)

		}

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})

	defer renderFrame.Release()
	js.Global().Call("requestAnimationFrame", renderFrame)

	<-c
}
