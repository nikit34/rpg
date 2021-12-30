package game


type itemType int

const (
	Weapon itemType = iota
	Helmet
	Other
)

type Item struct {
	Typ itemType
	Entity
	power float64
}

func NewSword(p Pos) *Item {
	return &Item{
		Typ: Weapon,
		Entity: Entity{
			Pos:  p,
			Name: "Sword",
			Rune: 's',
		},
		power: 2.0,
	}
}

func NewHelmet(p Pos) *Item {
	return &Item{
		Typ: Helmet,
		Entity: Entity{
			Pos:  p,
			Name: "Helmet",
			Rune: 'h',
		},
		power: 0.50,
	}
}
