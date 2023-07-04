package shared

import "github.com/Haraj-backend/hex-monscape/internal/core/entity"

type MonsterRow struct {
	ID            string         `dynamodbav:"id"`
	Name          string         `dynamodbav:"name"`
	BattleStats   BattleStatsRow `dynamodbav:"battle_stats"`
	AvatarURL     string         `dynamodbav:"avatar_url"`
	IsPartnerable int            `dynamodbav:"is_partnerable"`
}

func ToMonsterRow(m entity.Monster) MonsterRow {
	return MonsterRow{
		ID:   m.ID,
		Name: m.Name,
		BattleStats: BattleStatsRow{
			Health:    m.BattleStats.Health,
			MaxHealth: m.BattleStats.MaxHealth,
			Attack:    m.BattleStats.Attack,
			Defense:   m.BattleStats.Defense,
			Speed:     m.BattleStats.Speed,
		},
		AvatarURL: m.AvatarURL,
	}
}

func (r MonsterRow) ToMonster() *entity.Monster {
	return &entity.Monster{
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

type BattleStatsRow struct {
	Health    int `dynamodbav:"health"`
	MaxHealth int `dynamodbav:"max_health"`
	Attack    int `dynamodbav:"attack"`
	Defense   int `dynamodbav:"defense"`
	Speed     int `dynamodbav:"speed"`
}
