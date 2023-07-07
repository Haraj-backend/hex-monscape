package shared

import (
	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type MonsterRow struct {
	ID            string `db:"id"`
	Name          string `db:"name"`
	Health        int    `db:"health"`
	MaxHealth     int    `db:"max_health"`
	Attack        int    `db:"attack"`
	Defense       int    `db:"defense"`
	Speed         int    `db:"speed"`
	AvatarURL     string `db:"avatar_url"`
	IsPartnerable int    `db:"is_partnerable"`
}

func (r *MonsterRow) ToMonster() *entity.Monster {
	return &entity.Monster{
		ID:   r.ID,
		Name: r.Name,
		BattleStats: entity.BattleStats{
			Health:    r.Health,
			MaxHealth: r.MaxHealth,
			Attack:    r.Attack,
			Defense:   r.Defense,
			Speed:     r.Speed,
		},
		AvatarURL: r.AvatarURL,
	}
}

type MonsterRows []MonsterRow

func (r MonsterRows) ToMonsters() []entity.Monster {
	var monsters []entity.Monster
	for _, row := range r {
		monsters = append(monsters, *row.ToMonster())
	}
	return monsters
}

func ToMonsterRow(monster *entity.Monster) *MonsterRow {
	return &MonsterRow{
		ID:        monster.ID,
		Name:      monster.Name,
		Health:    monster.BattleStats.Health,
		MaxHealth: monster.BattleStats.MaxHealth,
		Attack:    monster.BattleStats.Attack,
		Defense:   monster.BattleStats.Defense,
		Speed:     monster.BattleStats.Speed,
		AvatarURL: monster.AvatarURL,
	}
}
