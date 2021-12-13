package main

import (
	"runtime"

	"./game"
	"./game/ui2d"
)


func main() {
	game := game.NewGame(3, "game/maps/level1.map")

	for i := 0; i < 3; i++ {
		go func(i int) {
			runtime.LockOSThread()
			ui := ui2d.NewUI(game.InputChan, game.LevelChans[i])
			ui.Run()
		}(i)
	}
	game.Run()
}