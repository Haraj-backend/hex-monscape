package gamestrg

import (
	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/shared"
)

type GameRow struct {
	ID         string             `db:"id"`
	PlayerName string             `db:"player_name"`
	CreatedAt  int64              `db:"created_at"`
	BattleWon  int                `db:"battle_won"`
	Scenario   entity.Scenario    `db:"scenario"`
	Partner    *shared.MonsterRow `db:"partner"`
}

func (r *GameRow) ToGame() *entity.Game {
	return &entity.Game{
		ID:         r.ID,
		PlayerName: r.PlayerName,
		CreatedAt:  r.CreatedAt,
		BattleWon:  r.BattleWon,
		Scenario:   r.Scenario,
		Partner:    r.Partner.ToMonster(),
	}
}

func NewGameRow(game *entity.Game) *GameRow {
	return &GameRow{
		ID:         game.ID,
		PlayerName: game.PlayerName,
		CreatedAt:  game.CreatedAt,
		BattleWon:  game.BattleWon,
		Scenario:   game.Scenario,
		Partner:    shared.ToMonsterRow(game.Partner),
	}
}
