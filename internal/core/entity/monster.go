package entity

const MinDamage = 5

type Monster struct {
	ID          string
	Name        string
	BattleStats BattleStats
	AvatarURL   string
}

// ResetBattleStats is used to reset partner health
// after lost or surrender from the battle
func (p *Monster) ResetBattleStats() {
	p.BattleStats.Health = p.BattleStats.MaxHealth
}

func (p Monster) IsDead() bool {
	return p.BattleStats.Health <= 0
}

// InflictDamage is used for inflicting damage to self based
// on given enemy. Returned the damage amount.
func (p *Monster) InflictDamage(enemy Monster) (int, error) {
	dmg := max(enemy.BattleStats.Attack-p.BattleStats.Defense, MinDamage)
	p.BattleStats.Health -= dmg
	return dmg, nil
}

type BattleStats struct {
	Health    int
	MaxHealth int
	Attack    int
	Defense   int
	Speed     int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
