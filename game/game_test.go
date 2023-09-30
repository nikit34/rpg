package game_test

import (
	"testing"

	"github.com/nikit34/rpg/game"
)


func TestLoadLevels(t *testing.T) {
	levels := game.ExportLoadLevels()
	if levels["level1"].Player.Character.Name != "GoMan" {
		t.Errorf("The name %s doesnt match %s", levels["level1"].Player.Character.Name, "GoMan")
	}
}
