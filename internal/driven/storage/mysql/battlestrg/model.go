package battlestrg

import (
	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
)

type battleRow struct {
	GameID         string          `db:"game_id"`
	PartnerLastDmg int             `db:"partner_last_damage"`
	EnemyLastDmg   int             `db:"enemy_last_damage"`
	State          battle.State    `db:"state"`
	Partner        *shared.PokeRow `db:"partner"`
	Enemy          *shared.PokeRow `db:"enemy"`
}

func (r *battleRow) ToBattle() *battle.Battle {
	return &battle.Battle{
		GameID: r.GameID,
		LastDamage: battle.LastDamage{
			Partner: r.PartnerLastDmg,
			Enemy:   r.EnemyLastDmg,
		},
		State:   r.State,
		Partner: r.Partner.ToPokemon(),
		Enemy:   r.Enemy.ToPokemon(),
	}
}
