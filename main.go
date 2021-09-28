package main

import "github.com/carltheperson/car-and-mouse/game"

func main() {
	g := game.NewGame("canvas")
	g.RunMainLoop()
}
