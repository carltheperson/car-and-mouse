package game

import "syscall/js"

func (g *Game) SetRestartFunc(function func()) {
	g.addInitialEntitiesFunc = function
}

func (g *Game) restartGame() {
	*g.Entities = []Entity{}
	if g.addInitialEntitiesFunc == nil {
		panic("No function was set to initialize entities. Call SetRestartFunc to correct this")
	}
	g.addInitialEntitiesFunc()
	g.Score = 0
	g.lastScore = -1
	g.State = StateNormal
}

func (g *Game) setTryAgainButtonEventListener() {
	eventListener := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.restartGame()
		return nil
	})
	g.document.Call("getElementById", "try-again-button").Call("addEventListener", "click", eventListener)
}
