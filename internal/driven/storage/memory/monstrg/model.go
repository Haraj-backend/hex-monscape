package monstrg

import "github.com/Haraj-backend/hex-monscape/internal/core/entity"

type monsterRow struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BattleStats struct {
		Health    int `json:"health"`
		MaxHealth int `json:"max_health"`
		Attack    int `json:"attack"`
		Defense   int `json:"defense"`
		Speed     int `json:"speed"`
	} `json:"battle_stats"`
	AvatarURL     string `json:"avatar_url"`
	IsPartnerable bool   `json:"is_partnerable"`
}

func (r monsterRow) toMonster() entity.Monster {
	return entity.Monster{
		ID:   r.ID,
		Name: r.Name,
		BattleStats: entity.BattleStats{
			Health:    r.BattleStats.Health,
			MaxHealth: r.BattleStats.MaxHealth,
			Attack:    r.BattleStats.Attack,
			Defense:   r.BattleStats.Defense,
			Speed:     r.BattleStats.Speed,
		},
		AvatarURL: r.AvatarURL,
	}
}
