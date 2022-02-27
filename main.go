package main

import (
	"github.com/rpg_go/game"
	"github.com/rpg_go/ui2d"
)


func main() {
	game := game.NewGame()

	go func() {
		game.Run()
	}()

	ui := ui2d.NewUI(game.InputChan, game.LevelChans)
	ui.Run()
}