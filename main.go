package main

import (
	"github.com/nikit34/rpg/game"
	"github.com/nikit34/rpg/ui2d"
)


func main() {
	game := game.NewGame()

	go func() {
		game.Run()
	}()

	ui := ui2d.NewUI(game.InputChan, game.LevelChans)
	ui.Run()
}