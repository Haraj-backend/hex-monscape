package mysql

import (
	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type monsterRow struct {
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

func newMonsterRow(monster *entity.Monster) *monsterRow {
	return &monsterRow{
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

func (r *monsterRow) toMonster() *entity.Monster {
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

type monsterRows []monsterRow

func (r monsterRows) toMonsters() []entity.Monster {
	var monsters []entity.Monster
	for _, row := range r {
		monsters = append(monsters, *row.toMonster())
	}

	return monsters
}

type gameRow struct {
	ID         string          `db:"id"`
	PlayerName string          `db:"player_name"`
	CreatedAt  int64           `db:"created_at"`
	BattleWon  int             `db:"battle_won"`
	Scenario   entity.Scenario `db:"scenario"`
	Partner    *monsterRow     `db:"partner"`
}

func newGameRow(game *entity.Game) *gameRow {
	return &gameRow{
		ID:         game.ID,
		PlayerName: game.PlayerName,
		CreatedAt:  game.CreatedAt,
		BattleWon:  game.BattleWon,
		Scenario:   game.Scenario,
		Partner:    newMonsterRow(game.Partner),
	}
}

func (r *gameRow) toGame() *entity.Game {
	return &entity.Game{
		ID:         r.ID,
		PlayerName: r.PlayerName,
		CreatedAt:  r.CreatedAt,
		BattleWon:  r.BattleWon,
		Scenario:   r.Scenario,
		Partner:    r.Partner.toMonster(),
	}
}

type battleRow struct {
	GameID            string `db:"game_id"`
	State             string `db:"state"`
	PartnerMonsterID  string `db:"partner_monster_id"`
	PartnerName       string `db:"partner_name"`
	PartnerMaxHealth  int    `db:"partner_max_health"`
	PartnerHealth     int    `db:"partner_health"`
	PartnerAttack     int    `db:"partner_attack"`
	PartnerDefense    int    `db:"partner_defense"`
	PartnerSpeed      int    `db:"partner_speed"`
	PartnerAvatarURL  string `db:"partner_avatar_url"`
	PartnerLastDamage int    `db:"partner_last_damage"`
	EnemyMonsterID    string `db:"enemy_monster_id"`
	EnemyName         string `db:"enemy_name"`
	EnemyMaxHealth    int    `db:"enemy_max_health"`
	EnemyHealth       int    `db:"enemy_health"`
	EnemyAttack       int    `db:"enemy_attack"`
	EnemyDefense      int    `db:"enemy_defense"`
	EnemySpeed        int    `db:"enemy_speed"`
	EnemyAvatarURL    string `db:"enemy_avatar_url"`
	EnemyLastDamage   int    `db:"enemy_last_damage"`
}

func newBattleRow(b entity.Battle) battleRow {
	return battleRow{
		GameID:            b.GameID,
		State:             string(b.State),
		PartnerMonsterID:  b.Partner.ID,
		PartnerName:       b.Partner.Name,
		PartnerMaxHealth:  b.Partner.BattleStats.MaxHealth,
		PartnerHealth:     b.Partner.BattleStats.Health,
		PartnerAttack:     b.Partner.BattleStats.Attack,
		PartnerDefense:    b.Partner.BattleStats.Defense,
		PartnerSpeed:      b.Partner.BattleStats.Speed,
		PartnerAvatarURL:  b.Partner.AvatarURL,
		PartnerLastDamage: b.LastDamage.Partner,
		EnemyMonsterID:    b.Enemy.ID,
		EnemyName:         b.Enemy.Name,
		EnemyMaxHealth:    b.Enemy.BattleStats.MaxHealth,
		EnemyHealth:       b.Enemy.BattleStats.Health,
		EnemyAttack:       b.Enemy.BattleStats.Attack,
		EnemyDefense:      b.Enemy.BattleStats.Defense,
		EnemySpeed:        b.Enemy.BattleStats.Speed,
		EnemyAvatarURL:    b.Enemy.AvatarURL,
		EnemyLastDamage:   b.LastDamage.Enemy,
	}
}

func (r battleRow) toBattle() *entity.Battle {
	return &entity.Battle{
		GameID: r.GameID,
		State:  entity.State(r.State),
		Partner: &entity.Monster{
			ID:   r.PartnerMonsterID,
			Name: r.PartnerName,
			BattleStats: entity.BattleStats{
				Health:    r.PartnerHealth,
				MaxHealth: r.PartnerMaxHealth,
				Attack:    r.PartnerAttack,
				Defense:   r.PartnerDefense,
				Speed:     r.PartnerSpeed,
			},
			AvatarURL: r.PartnerAvatarURL,
		},
		Enemy: &entity.Monster{
			ID:   r.EnemyMonsterID,
			Name: r.EnemyName,
			BattleStats: entity.BattleStats{
				Health:    r.EnemyHealth,
				MaxHealth: r.EnemyMaxHealth,
				Attack:    r.EnemyAttack,
				Defense:   r.EnemyDefense,
				Speed:     r.EnemySpeed,
			},
			AvatarURL: r.EnemyAvatarURL,
		},
		LastDamage: entity.LastDamage{
			Partner: r.PartnerLastDamage,
			Enemy:   r.EnemyLastDamage,
		},
	}
}
