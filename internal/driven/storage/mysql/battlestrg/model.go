package battlestrg

import (
	"github.com/Haraj-backend/hex-monscape/internal/core/battle"
	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type battleRow struct {
	GameID            string `db:"game_id"`
	State             string `db:"state"`
	PartnerPokemonID  string `db:"partner_monster_id"`
	PartnerName       string `db:"partner_name"`
	PartnerMaxHealth  int    `db:"partner_max_health"`
	PartnerHealth     int    `db:"partner_health"`
	PartnerAttack     int    `db:"partner_attack"`
	PartnerDefense    int    `db:"partner_defense"`
	PartnerSpeed      int    `db:"partner_speed"`
	PartnerAvatarURL  string `db:"partner_avatar_url"`
	PartnerLastDamage int    `db:"partner_last_damage"`
	EnemyPokemonID    string `db:"enemy_monster_id"`
	EnemyName         string `db:"enemy_name"`
	EnemyMaxHealth    int    `db:"enemy_max_health"`
	EnemyHealth       int    `db:"enemy_health"`
	EnemyAttack       int    `db:"enemy_attack"`
	EnemyDefense      int    `db:"enemy_defense"`
	EnemySpeed        int    `db:"enemy_speed"`
	EnemyAvatarURL    string `db:"enemy_avatar_url"`
	EnemyLastDamage   int    `db:"enemy_last_damage"`
}

func (r battleRow) ToBattle() *battle.Battle {
	return &battle.Battle{
		GameID: r.GameID,
		State:  battle.State(r.State),
		Partner: &entity.Monster{
			ID:   r.PartnerPokemonID,
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
			ID:   r.EnemyPokemonID,
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
		LastDamage: battle.LastDamage{
			Partner: r.PartnerLastDamage,
			Enemy:   r.EnemyLastDamage,
		},
	}
}

func newBattleRow(b battle.Battle) battleRow {
	return battleRow{
		GameID:            b.GameID,
		State:             string(b.State),
		PartnerPokemonID:  b.Partner.ID,
		PartnerName:       b.Partner.Name,
		PartnerMaxHealth:  b.Partner.BattleStats.MaxHealth,
		PartnerHealth:     b.Partner.BattleStats.Health,
		PartnerAttack:     b.Partner.BattleStats.Attack,
		PartnerDefense:    b.Partner.BattleStats.Defense,
		PartnerSpeed:      b.Partner.BattleStats.Speed,
		PartnerAvatarURL:  b.Partner.AvatarURL,
		PartnerLastDamage: b.LastDamage.Partner,
		EnemyPokemonID:    b.Enemy.ID,
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
