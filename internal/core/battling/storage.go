package battling

import (
	"context"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/playing"
)

type BattleStorage interface {
	GetBattle(ctx context.Context, gameID string) (*Battle, error)
	SaveBattle(ctx context.Context, b Battle) error
}

type GameStorage interface {
	GetGame(ctx context.Context, gameID string) (*playing.Game, error)
	SaveGame(ctx context.Context, game playing.Game) error
}

type PokemonStorage interface {
	GetPossibleEnemies(ctx context.Context) ([]entity.Pokemon, error)
}
