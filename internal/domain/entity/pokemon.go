package entity

type Pokemon struct {
	ID          string
	Name        string
	BattleStats BattleStats
	AvatarURL   string
}

func (p Pokemon) IsDead() bool {
	// TODO
	return false
}

func (p *Pokemon) InflictDamage(enemy Pokemon) (int, error) {
	// TODO
	return 0, nil
}

type BattleStats struct {
	Health    int
	MaxHealth int
	Attack    int
	Defense   int
	Speed     int
}
