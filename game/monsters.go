package game


type Monster struct {
	Character
}

func NewRat(p Pos) *Monster {
	return &Monster {
		Character: Character{
			Entity: Entity{
				Pos: p,
				Name: "Rat",
				Rune: 'R',
			},
			Hitpoints: 20,
			Strength: 0,
			Speed: 1.5,
			ActionPoints: 0.0,
			SightRange: 10.0,
			Items: []*Item{NewSword(Pos{})},
		},
	}
}

func NewSpider(p Pos) *Monster {
	return &Monster{
		Character: Character{
			Entity: Entity{
				Pos: p,
				Name: "Spider",
				Rune: 'S',
			},
			Hitpoints: 200,
			Strength: 0,
			Speed: 1.0,
			ActionPoints: 0.0,
			SightRange: 10.0,
		},
	}
}

func (m *Monster) Update(level *Level){
	m.ActionPoints += m.Speed
	playerPos := level.Player.Pos
	apInt := int(m.ActionPoints)
	positions := level.astar(m.Pos, playerPos)
	if len(positions) == 0 {
		m.Pass()
		return
	}
	moveIndex := 1
	for i := 0; i < apInt; i++ {
		if moveIndex < len(positions) {
			m.Move(positions[moveIndex], level)
			moveIndex++
			m.ActionPoints--
		}
	}
}

func (m *Monster) Pass() {
	m.ActionPoints -= m.Speed
}

func (m *Monster) Kill(level *Level) {
	delete(level.Monsters, m.Pos)

	groundItems := level.Items[m.Pos]
	for _, item := range m.Items {
		item.Pos = m.Pos
		groundItems = append(groundItems, item)
	}
	level.Items[m.Pos] = groundItems
}


func (m *Monster) Move(to Pos, level *Level) {
	_, exists := level.Monsters[to]
	if !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
		return
	}
	if to == level.Player.Pos {
		level.Attack(&m.Character, &level.Player.Character)
		if m.Hitpoints <= 0 {
			m.Kill(level)
		}
		if level.Player.Hitpoints <= 0 {
			panic("DIED")
		}
	}
}
