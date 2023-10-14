package game

import "testing"


func TestLoadLevels(t *testing.T) {
	levels := loadLevels()
	if levels["level1"].Player.Character.Name != "GoMan" {
		t.Errorf("The name %s doesnt match %s", levels["level1"].Player.Character.Name, "GoMan")
	}
}
