package game


type Monster struct {
	Pos
	Rune rune
	Name string
	Hitpoints int
	Strength int
	Speed float64
}

func NewRat(p Pos) *Monster {
	return &Monster{p, 'R', "Rat", 5, 5, 2.0}
}

func NewSpider(p Pos) *Monster {
	return &Monster{p, 'S', "Spider", 10, 10, 1.0}
}

func (m *Monster) Update(level *Level){
	playerPos := level.Player.Pos
	positions := level.astar(m.Pos, playerPos)

	if len(positions) > 1 {
		m.Move(positions[1], level)
	}
}

func (m *Monster) Move(to Pos, level *Level) {
	_, exists := level.Monsters[to]
	if !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
	}
}