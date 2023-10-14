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
	testCases := []struct {
		x int
		y int
		condition bool
		msg string
	} {
		{0, -1, true, "Current position {%d, %d} outside the field"},
		{6, 6, true, "Current position {%d, %d} outside the field"},
		{0, 0, false, "Current position {%d, %d} side the field"},
		{5, 5, false, "Current position {%d, %d} side the field"},
	}
	for _, tc := range testCases {
		res := inRange(level, Pos{tc.x, tc.y})
		if (res == tc.condition) {
			t.Errorf(tc.msg, tc.x, tc.y)
		}
	}
}
