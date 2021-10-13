package game

import (
	"fmt"
	"strconv"
)

func (g *Game) getHighscore() int {
	highscoreStringJS := g.window.Get("localStorage").Call("getItem", highscoreLocalStorageKey)
	if highscoreStringJS.IsNull() {
		g.setHighscoreInLocalStorage(0)
		return 0
	}
	highscoreString := highscoreStringJS.String()
	highscore, err := strconv.Atoi(highscoreString)
	if err != nil {
		panic(err)
	}
	return highscore
}

func (g *Game) setHighscoreInLocalStorage(newHighscore int) {
	g.window.Get("localStorage").Call("setItem", highscoreLocalStorageKey, fmt.Sprint(newHighscore))
}
