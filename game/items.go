package game


type Item struct {
	Entity
}

func NewSword(p Pos) *Item {
	return &Item{
		Entity: Entity{
			Pos:  p,
			Name: "Sword",
			Rune: 's',
		},
	}
}

func NewHelmet(p Pos) *Item {
	return &Item{
		Entity: Entity{
			Pos:  p,
			Name: "Helmet",
			Rune: 'h',
		},
	}
}
