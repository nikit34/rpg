package game

import (
	"bufio"
	"os"
)


type Input struct {
	up bool
	down bool
	left bool
	right bool
}

type GameUI interface {
	DrawThenGetInput(*Level) Input
}

type Tile rune

const (
	StoneWall Tile = '#'
	DirtFloor Tile = '.'
	Door Tile = '|'
	Blank Tile = 0
	Pending Tile = -1
)

type Entity struct {
	X, Y int
}

type Player struct {
	Entity
}

type Level struct {
	Map [][]Tile
	Player Player
}


func loadLevelFromFile(filename string) *Level {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	levelLines := make([]string, 0)
	longestRow := 0
	index := 0
	for scanner.Scan() {
		levelLines = append(levelLines, scanner.Text())
		if len(levelLines[index]) > longestRow {
			longestRow = len(levelLines[index])
		}
		index++
	}
	level := &Level{}
	level.Map := make([][]Tile, len(levelLines))
	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}
	for y := 0; y < len(level.Map); y++ {
		line := levelLines[y]
		for x, c := range line {
			var t Tile
			switch c {
			case ' ', '\t', '\n', '\r':
				t = Blank
			case '#':
				t = StoneWall
			case '|':
				t = Door
			case '.':
				t = DirtFloor
			case 'P':
				level.Player.X = x * 32
				level.Player.Y = y * 32
				t = Pending
			default:
				panic("Invalid character in map")
			}
			level.Map[y][x] = t
		}
	}
	for y, row := range level.Map {
		for x, tile := range row {
			if tile == Pending {
			SearchLoop:
				for searchX := x - 1; searchX <= x+1; searchX++ {
					for searchY := y - 1; searchY <= y+1; searchY++ {
						searchTile := level.Map[searchY][searchX]
						switch searchTile {
						case DirtFloor:
							level.Map[y][x] = DirtFloor
							break SearchLoop
						default:
							panic("Error in searchTile")
						}
					}
				}
			}
		}
	}
	return level
}


func Run(ui GameUI) {
	level := loadLevelFromFile("games/maps/level1.map")
	for {
		_ = ui.DrawThenGetInput(level)
	}
}