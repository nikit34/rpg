package game

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var (
    _, b, _, _ = runtime.Caller(0)
    basepath   = filepath.Dir(b)
)

type Game struct {
	LevelChans chan *Level
	InputChan chan *Input
	Levels map[string]*Level
	CurrentLevel *Level
}

func NewGame() *Game {
	levelChans := make(chan *Level)
	inputChan := make(chan *Input)
	levels := loadLevels()

	game := &Game{levelChans, inputChan, levels, nil}
	game.loadWorldFile()
	game.CurrentLevel.lineOfSight()

	return game
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	TakeAll
	QuitGame
	CloseWindow
	TakeItem
	DropItem
	EquipItem
	Search
)

type Input struct {
	Typ InputType
	Item *Item
	LevelChannel chan *Level
}

type Tile struct {
	Rune rune
	OverlayRune rune
	Visible bool
	Seen bool
}

const (
	StoneWall rune = '#'
	DirtFloor = '.'
	ClosedDoor = '|'
	OpenDoor = '/'
	UpStair = 'u'
	DownStair = 'd'
	Blank = 0
	Pending = -1
)

type Pos struct {
	X, Y int
}

type LevelPos struct {
	*Level
	Pos
}

type Entity struct {
	Pos
	Name string
	Rune rune
}

type Player struct {
	Character
}

type Character struct {
	Entity
	Hitpoints int
	Strength int
	Speed float64
	ActionPoints float64
	SightRange int
	Items []*Item
	Helmet *Item
	Weapon *Item
}

type GameEvent int

const (
	Move GameEvent = iota
	DoorOpen
	Attack
	Hit
	Portal
	PickUp
	Drop
)

type Level struct {
	Map [][]Tile
	Player *Player
	Monsters map[Pos]*Monster
	Items map[Pos][]*Item
	Portals map[Pos]*LevelPos
	Events []string
	EventPos int
	Debug map[Pos]bool
	LastEvent GameEvent
}

func (level *Level) DropItem(itemDrop *Item, character *Character) {
	pos := character.Pos
	items := character.Items
	for i, item := range items {
		if item == itemDrop {
			character.Items = append(character.Items[:i], character.Items[i+1:]...)
			level.Items[pos] = append(level.Items[pos], item)
			level.AddEvent(character.Name + " dropped 1x " + item.Name)
			return
		}
	}
	panic("Tried to drop an item we don't have.")
}

func (level *Level) MoveItem(itemToMove *Item, character *Character) {
	pos := character.Pos
	items := level.Items[pos]
	for i, item := range items {
		if item == itemToMove {
			items = append(items[:i], items[i+1:]...)
			level.Items[pos] = items
			character.Items = append(character.Items, item)
			level.AddEvent(character.Name + " picked up 1x " + item.Name)
			return
		}
	}
	panic("Tried to move an item we're not on top of")
}

func (level *Level) Attack(c1, c2 *Character) {
	c1.ActionPoints--
	c1AttackPower := c1.Strength

	if c1.Weapon != nil {
		c1AttackPower = int(float64(c1AttackPower) * c1.Weapon.power)
	}
	damage := c1AttackPower

	if c2.Helmet != nil {
		damage = int(float64(damage) * (1.0 - c2.Helmet.power))
	}

	c2.Hitpoints -= damage

	if c2.Hitpoints > 0 {
		level.AddEvent(c1.Name + " Attacked " + c2.Name + " for " + strconv.Itoa(damage))
	} else {
		level.AddEvent(c1.Name + " Killed " + c2.Name)
	}
}

func (level *Level) AddEvent(event string) {
	level.Events[level.EventPos] = event
	level.EventPos++
	if level.EventPos == len(level.Events) {
		level.EventPos = 0
	}
}

func (level *Level) lineOfSight() {
	pos := level.Player.Pos
	dist := level.Player.SightRange
	for y := pos.Y - dist; y <= pos.Y + dist; y++ {
		for x := pos.X - dist; x <= pos.X + dist; x++ {
			xDelta := pos.X - x
			yDelta := pos.Y - y
			d := math.Sqrt(float64(xDelta * xDelta + yDelta * yDelta))
			if d <= float64(dist){
				level.bresenham(pos, Pos{x, y})
			}
		}
	}
}

func (level *Level) setVisibleSeen(steep bool, x, y int) Pos {
	var pos Pos
	if steep {
		pos = Pos{y, x}
	} else {
		pos = Pos{x, y}
	}
	level.Map[pos.Y][pos.X].Visible = true
	level.Map[pos.Y][pos.X].Seen = true
	return pos
}

func (level *Level) bresenham(start Pos, end Pos) {
	steep := math.Abs(float64(end.Y - start.Y)) > math.Abs(float64(end.X - start.X))
	if steep {
		start.X, start.Y = start.Y, start.X
		end.X, end.Y = end.Y, end.X
	}
	deltaY := int(math.Abs(float64(end.Y - start.Y)))

	err := 0
	y := start.Y
	ystep := 1
	if start.Y >= end.Y {
		ystep = -1
	}
	if start.X > end.X {
		deltaX := start.X - end.X
		for x := start.X; x > end.X; x-- {
			pos := level.setVisibleSeen(steep, x, y)
			if !canSeeThrough(level, pos) {
				return
			}
			err += deltaY
			if 2 * err >= deltaX {
				y += ystep
				err -= deltaX
			}
		}
	} else {
		deltaX := end.X - start.X
		for x := start.X; x < end.X; x++ {
			pos := level.setVisibleSeen(steep, x, y)
			if !canSeeThrough(level, pos) {
				return
			}
			err += deltaY
			if 2 * err >= deltaX {
				y += ystep
				err -= deltaX
			}
		}
	}
}

func (game *Game) loadWorldFile() {
	file, err := os.Open("game/maps/world.txt")
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	for rowIndex, row := range rows {
		if rowIndex == 0 {
			game.CurrentLevel = game.Levels[row[0]]
			if game.CurrentLevel == nil {
				fmt.Println("Couldn't find current level name in world file.")
				panic(nil)
			}
			continue
		}
		levelWithPortal := game.Levels[row[0]]
		if levelWithPortal == nil {
			fmt.Println("Couldn't find current level name in world file.")
			panic(nil)
		}
		x, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(row[2], 10, 64)
		if err != nil {
			panic(err)
		}
		pos := Pos{int(x), int(y)}

		levelToTeleportTo := game.Levels[row[3]]
		if levelToTeleportTo == nil {
			fmt.Println("Couldn't find current level name in world file.")
			panic(nil)
		}
		x, err = strconv.ParseInt(row[4], 10, 64)
		if err != nil {
			panic(err)
		}
		y, err = strconv.ParseInt(row[5], 10, 64)
		if err != nil {
			panic(err)
		}
		posToTeleportTo := Pos{int(x), int(y)}
		levelWithPortal.Portals[pos] = &LevelPos{levelToTeleportTo, posToTeleportTo}
	}
}

func getSeparaterPath() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}
	return "/"
}

func loadLevels() map[string]*Level {
	player := &Player{
		Character: Character {
			Strength: 5,
			Hitpoints: 100,
			Speed: 1.0,
			ActionPoints: 0,
			SightRange: 7,
			Entity: Entity {
				Name: "GoMan",
				Rune: '@',
			},
		},
	}

	levels := make(map[string]*Level)

	filenames, err := filepath.Glob(basepath + "/maps/*.map")
	if err != nil {
		panic(err)
	}
	for _, filename := range filenames {
		extIndex := strings.LastIndex(filename, ".map")
		lastSlashIndex := strings.LastIndex(filename, getSeparaterPath())
		levelName := filename[lastSlashIndex + 1 : extIndex]
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
		level.Debug = make(map[Pos]bool)
		level.Events = make([]string, 10)
		level.Player = player

		level.Map = make([][]Tile, len(levelLines))
		level.Monsters = make(map[Pos] *Monster)
		level.Items = make(map[Pos][] *Item)
		level.Portals = make(map[Pos] *LevelPos)

		for i := range level.Map {
			level.Map[i] = make([]Tile, longestRow)
		}

		for y := range level.Map {
			line := levelLines[y]
			for x, c := range line {
				pos := Pos{x, y}
				var t Tile
				switch c {
				case ' ', '\t', '\n', '\r':
					t.OverlayRune = Blank
					t.Rune = Blank
				case '#':
					t.OverlayRune = Blank
					t.Rune = StoneWall
				case '|':
					t.OverlayRune = ClosedDoor
					t.Rune = Pending
				case '/':
					t.OverlayRune = Blank
					t.Rune = OpenDoor
				case 'u':
					t.OverlayRune = UpStair
					t.Rune = Pending
				case 'd':
					t.OverlayRune = DownStair
					t.Rune = Pending
				case 's':
					t.OverlayRune = Blank
					level.Items[pos] = append(level.Items[pos], NewSword(pos))
					level.Items[pos] = append(level.Items[pos], NewHelmet(pos))
					t.Rune = Pending
				case 'h':
					t.OverlayRune = Blank
					level.Items[pos] = append(level.Items[pos], NewHelmet(pos))
					t.Rune = Pending
				case '.':
					t.OverlayRune = Blank
					t.Rune = DirtFloor
				case '@':
					t.OverlayRune = Blank
					level.Player.X = x
					level.Player.Y = y
					t.Rune = Pending
				case 'R':
					t.OverlayRune = Blank
					level.Monsters[pos] = NewRat(pos)
					t.Rune = Pending
				case 'S':
					t.OverlayRune = Blank
					level.Monsters[pos] = NewSpider(pos)
					t.Rune = Pending
				default:
					panic("Invalid character in map!")
				}
				level.Map[y][x] = t
			}
		}
		for y, row := range level.Map {
			for x, tile := range row {
				if tile.Rune == Pending {
					level.Map[y][x].Rune = level.bfsFloor(Pos{x, y})
				}
			}
		}
		levels[levelName] = level
	}
	return levels
}

func inRange(level *Level, pos Pos) bool {
	return pos.X < len(level.Map[0]) && pos.Y < len(level.Map) && pos.X >= 0 && pos.Y >= 0
}

func canWalk(level *Level, pos Pos) bool {
	if inRange(level, pos) {
		t := level.Map[pos.Y][pos.X]
		switch t.Rune {
		case StoneWall, Blank:
			return false
		}
		switch t.OverlayRune {
		case ClosedDoor:
			return false
		}
		_, exists := level.Monsters[pos]
		return !exists
	}
	return false
}

func canSeeThrough(level *Level, pos Pos) bool {
	if inRange(level, pos) {
		t := level.Map[pos.Y][pos.X]
		switch t.Rune {
		case StoneWall, Blank:
			return false
		}
		switch t.OverlayRune {
		case ClosedDoor:
			return false
		default:
			return true
		}
	}
	return false
}

func checkDoor(level *Level, pos Pos) {
	t := level.Map[pos.Y][pos.X]
	if t.OverlayRune == ClosedDoor {
		level.Map[pos.Y][pos.X].OverlayRune = OpenDoor
		level.LastEvent = OpenDoor
		level.lineOfSight()
	}
}

func (game *Game) Move(to Pos){
	level := game.CurrentLevel
	player := level.Player
	levelAndPos := level.Portals[to]
	if levelAndPos != nil {
		game.CurrentLevel = levelAndPos.Level
		game.CurrentLevel.Player.Pos = levelAndPos.Pos
		game.CurrentLevel.lineOfSight()
	} else {
		player.Pos = to
		level.LastEvent = Move
		for y, row := range level.Map {
			for x := range row {
				level.Map[y][x].Visible = false
			}
		}
		level.lineOfSight()
	}
}

func (game *Game) resolveMovement(pos Pos) {
	level := game.CurrentLevel
	monster, exists := level.Monsters[pos]
	if exists {
		level.Attack(&level.Player.Character, &monster.Character)
		level.LastEvent = Attack
		if monster.Hitpoints <= 0 {
			monster.Kill(level)
		}
		if level.Player.Hitpoints <= 0 {
			panic("ded")
		}
	} else if canWalk(level, pos) {
		game.Move(pos)
	} else {
		checkDoor(level, pos)
	}
}

func equip(c *Character, itemToEquip *Item) {
	for i, item := range c.Items {
		if item == itemToEquip {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			if itemToEquip.Typ == Helmet {
				c.Helmet = itemToEquip
			} else if itemToEquip.Typ == Weapon {
				c.Weapon = itemToEquip
			}
			return
		}
	}
	panic("Tried to equip something you don't have.")
}

func (game *Game) handleInput(input *Input) {
	level := game.CurrentLevel
	p := level.Player
	switch input.Typ {
	case Up:
		newPos := Pos{p.X, p.Y - 1}
		game.resolveMovement(newPos)
	case Down:
		newPos := Pos{p.X, p.Y + 1}
		game.resolveMovement(newPos)
	case Left:
		newPos := Pos{p.X - 1, p.Y}
		game.resolveMovement(newPos)
	case Right:
		newPos := Pos{p.X + 1, p.Y}
		game.resolveMovement(newPos)
	case TakeItem:
		level.MoveItem(input.Item, &p.Character)
		level.LastEvent = PickUp
	case DropItem:
		level.DropItem(input.Item, &p.Character)
		level.LastEvent = Drop
	case TakeAll:
		for _, item := range level.Items[p.Pos] {
			level.MoveItem(item, &p.Character)
		}
		level.LastEvent = PickUp
	case EquipItem:
		equip(&level.Player.Character, input.Item)
	case CloseWindow:
		close(input.LevelChannel)
	}
}

func (level *Level) bfsFloor(start Pos) rune {
	frontier := make([]Pos, 0, 8)
	frontier = append(frontier, start)
	visited := make(map[Pos]bool)
	visited[start] = true

	for len(frontier) > 0 {
		current := frontier[0]
		currentTile := level.Map[current.Y][current.X]
		switch currentTile.Rune {
		case DirtFloor:
			return DirtFloor
		default:
		}

		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
			}
		}
	}
	return DirtFloor
}

func getNeighbors(level *Level, pos Pos) []Pos {
	neighbors := make([]Pos, 0, 4)
	dirs := make([]Pos, 0, 4)
	dirs = append(dirs, Pos{pos.X - 1, pos.Y})
	dirs = append(dirs, Pos{pos.X + 1, pos.Y})
	dirs = append(dirs, Pos{pos.X, pos.Y - 1})
	dirs = append(dirs, Pos{pos.X, pos.Y + 1})

	for _, dir := range dirs {
		if canWalk(level, dir) {
			neighbors = append(neighbors, dir)
		}
	}
	return neighbors
}

func (level *Level) astar(start, goal Pos) []Pos {
	frontier := make(pqueue, 0, 8)
	frontier = frontier.push(start, 1)

	cameFrom := make(map[Pos]Pos)
	cameFrom[start] = start

	costSoFar := make(map[Pos]int)
	costSoFar[start] = 0

	var current Pos
	for len(frontier) > 0 {
		frontier, current = frontier.pop()

		if current == goal {
			path := make([]Pos, 0)

			p := current
			for p != start {
				path = append(path, p)
				p = cameFrom[p]
			}
			path = append(path, p)

			for i, j := 0, len(path) - 1; i < j; i, j = i + 1, j - 1 {
				path[i], path[j] = path[j], path[i]
			}
			return path
		}

		for _, next := range getNeighbors(level, current) {
			newCost := costSoFar[current] + 1
			_, exists := costSoFar[next]
			if !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				xDist := int(math.Abs(float64(goal.X - next.X)))
				yDist := int(math.Abs(float64(goal.Y - next.Y)))
				priority := newCost + xDist + yDist
				frontier = frontier.push(next, priority)
				cameFrom[next] = current
			}
		}
	}
	return nil
}

func (game *Game) Run() {

	game.LevelChans <- game.CurrentLevel

	for input := range game.InputChan {
		if input.Typ == QuitGame {
			return
		}

		game.handleInput(input)

		for _, monster := range game.CurrentLevel.Monsters {
			monster.Update(game.CurrentLevel)
		}

		game.LevelChans <- game.CurrentLevel
	}
}
