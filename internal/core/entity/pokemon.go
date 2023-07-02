package entity

const minDamage = 5

type Monster struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	BattleStats BattleStats `json:"battle_stats"`
	AvatarURL   string      `json:"avatar_url"`
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
	dmg := max(enemy.BattleStats.Attack-p.BattleStats.Defense, minDamage)
	p.BattleStats.Health -= dmg
	return dmg, nil
}

type BattleStats struct {
	Health    int `json:"health"`
	MaxHealth int `json:"max_health"`
	Attack    int `json:"attack"`
	Defense   int `json:"defense"`
	Speed     int `json:"speed"`
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
