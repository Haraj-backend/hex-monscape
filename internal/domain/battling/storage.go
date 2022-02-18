package battling

import (
	"context"

	"github.com/Haraj-backend/hex-pokebattle/internal/domain/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/domain/playing"
)

type BattleStorage interface {
	GetBattle(ctx context.Context, gameID string) (*Battle, error)
	SaveBattle(ctx context.Context, b Battle) error
}

type GameStorage interface {
	GetGame(ctx context.Context, gameID string) (*playing.Game, error)
}

type PokemonStorage interface {
	GetPossibleEnemies(ctx context.Context) ([]entity.Pokemon, error)
}
