package battling

import (
	"context"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
)

type BattleStorage interface {
	GetBattle(ctx context.Context, gameID string) (*Battle, error)
	SaveBattle(ctx context.Context, b Battle) error
}

type GameStorage interface {
	GetGame(ctx context.Context, gameID string) (*entity.Game, error)
	SaveGame(ctx context.Context, game entity.Game) error
}

type PokemonStorage interface {
	GetPossibleEnemies(ctx context.Context) ([]entity.Pokemon, error)
}
