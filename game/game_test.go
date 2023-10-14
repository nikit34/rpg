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

func TestCanWalk(t *testing.T) {
	level := &Level{}
	level.Map = make([][]Tile, 2)
	for i := range level.Map {
		level.Map[i] = make([]Tile, 2)
	}
	level.Map[0][0].Rune = DirtFloor
	level.Map[0][1].OverlayRune = ClosedDoor
	level.Map[1][0].Rune = DirtFloor
	level.Map[1][1].Rune = StoneWall
	testCases := []struct {
		x int
		y int
		condition bool
		msg string
	} {
		{0, 0, false, "You can go this way {%d, %d}"},
		{0, 1, false, "You can go this way {%d, %d}"},
		{1, 0, true, "You can't go this way {%d, %d}"},
		{1, 1, true, "You can't go this way {%d, %d}"},
	}
	for _, tc := range testCases {
		res := canWalk(level, Pos{tc.x, tc.y})
		if (res == tc.condition) {
			t.Errorf(tc.msg, tc.x, tc.y)
		}
	}
}

func TestCanSeeThrough(t *testing.T) {
	level := &Level{}
	level.Map = make([][]Tile, 2)
	for i := range level.Map {
		level.Map[i] = make([]Tile, 2)
	}
	level.Map[0][0].Rune = DirtFloor
	level.Map[0][1].OverlayRune = ClosedDoor
	level.Map[1][0].Rune = DirtFloor
	level.Map[1][1].Rune = StoneWall
	testCases := []struct {
		x int
		y int
		condition bool
		msg string
	} {
		{0, 0, false, "You can see {%d, %d}"},
		{0, 1, false, "You can see {%d, %d}"},
		{1, 0, true, "You can't see {%d, %d}"},
		{1, 1, true, "You can't see {%d, %d}"},
	}
	for _, tc := range testCases {
		res := canSeeThrough(level, Pos{tc.x, tc.y})
		if (res == tc.condition) {
			t.Errorf(tc.msg, tc.x, tc.y)
		}
	}
}
