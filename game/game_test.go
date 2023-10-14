package game

import (
	"testing"
)


func TestLoadLevels(t *testing.T) {
	levels := loadLevels()
	if levels["level1"].Player.Character.Name != "GoMan" {
		t.Errorf("The name %s doesnt match %s", levels["level1"].Player.Character.Name, "GoMan")
	}
}

func TestInRange(t *testing.T) {
	level := &Level{}
	level.Map = make([][]Tile, 6)
	for i := range level.Map {
		level.Map[i] = make([]Tile, 6)
	}
	x := 0
	y := -1
	res := inRange(level, Pos{x, y})
	if (res) {
		t.Errorf("Current position {%d, %d} outside the field", x, y)
	}
	x = 6
	y = 6
	res = inRange(level, Pos{x, y})
	if (res) {
		t.Errorf("Current position {%d, %d} outside the field", x, y)
	}
	x = 0
	y = 0
	res = inRange(level, Pos{x, y})
	if (!res) {
		t.Errorf("Current position {%d, %d} side the field", x, y)
	}
	x = 5
	y = 5
	res = inRange(level, Pos{x, y})
	if (!res) {
		t.Errorf("Current position {%d, %d} side the field", x, y)
	}
}
