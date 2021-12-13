package main

import (
	"runtime"

	"./game"
	"./game/ui2d"
)


func main() {
	numWindows := 1
	game := game.NewGame(numWindows, "game/maps/level1.map")

	for i := 0; i < numWindows; i++ {
		go func(i int) {
			runtime.LockOSThread()
			ui := ui2d.NewUI(game.InputChan, game.LevelChans[i])
			ui.Run()
		}(i)
	}
	game.Run()
}