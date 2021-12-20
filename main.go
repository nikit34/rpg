package main

import (
	"./game"
	"./game/ui2d"
)


func main() {
	game := game.NewGame(1, "game/maps/level1.map")

	go func() {
		game.Run()
	}()
	}
	ui := ui2d.NewUI(game.InputChan, game.LevelChans[0])
	ui.Run()
}