package game

import (
	"reflect"
	"testing"
)

var pq = pqueue{
	{Pos{0, 0}, 1},
	{Pos{0, 1}, 2},
	{Pos{0, 2}, 2},
	{Pos{0, 3}, 3},
	{Pos{0, 4}, 3},
	{Pos{0, 5}, 3},
	{Pos{0, 6}, 3},
	{Pos{0, 7}, 4},
	{Pos{0, 8}, 5},
	{Pos{0, 9}, 5},
}

func TestParent(t *testing.T) {
	index, pqParent := pq.parent(7)
	if index != 3 {
		t.Errorf("Incorrect parent node index")
	}
	pos := Pos{0, 3}
	if pqParent.Pos != pos {
		t.Errorf("Incorrect parent node")
	}
}

func TestLeft(t *testing.T) {
	res, index, pqLeft := pq.left(4)
	if !res {
		t.Errorf("Incorrect result")
	}
	if index != 9 {
		t.Errorf("Incorrect left node index")
	}
	pos := Pos{0, 9}
	if pqLeft.Pos != pos {
		t.Errorf("Incorrect left node")
	}
}

func TestRight(t *testing.T) {
	res, index, pqRight := pq.right(4)
	if res {
		t.Errorf("Incorrect result")
	}
	if index != 0 {
		t.Errorf("Incorrect right node index")
	}
	pos := Pos{0, 0}
	if pqRight.Pos != pos {
		t.Errorf("Incorrect right node")
	}
}

func TestPush(t *testing.T) {
	pqCopy := make(pqueue, 10)
	copy(pqCopy, pq)
	pqItem := priorityPos{
		Pos{1, 10},
		2,
	}
	pqCopy = pqCopy.push(pqItem.Pos, pqItem.priority)
	if (!reflect.DeepEqual(pqCopy[:4], pq[:4])) {
		t.Error("Adding an element changed the original queue in an unexpected way")
	}
	if (!reflect.DeepEqual(pqCopy[5:10], pq[5:])) {
		t.Error("Adding an element changed the original queue in an unexpected way")
	}
	if (pqCopy[4] != pqItem) {
		t.Error("The added element is not equal to the last one")
	}
}

