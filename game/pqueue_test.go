package game

import (
	"testing"
)


func TestParent(t *testing.T) {
	pq := pqueue{
		{Pos{0, 0}, 1},
		{Pos{0, 1}, 2},
		{Pos{0, 2}, 3},
		{Pos{0, 3}, 3},
		{Pos{0, 4}, 3},
		{Pos{0, 5}, 2},
		{Pos{0, 6}, 1},
		{Pos{0, 7}, 1},
		{Pos{0, 8}, 4},
		{Pos{0, 9}, 3},
	}
	index, pqParent := pq.parent(7)
	if index != 3 {
		t.Errorf("Incorrect parent node index")
	}
	pos := Pos{0, 3}
	if pqParent.Pos != pos {
		t.Errorf("Incorrect parent node")
	}
}
