package playing

import (
	"context"

	"github.com/riandyrn/pokebattle/internal/domain/entity"
)

type Storage interface {
	GetGame(ctx context.Context, gameID string) (*Game, error)
	SaveGame(ctx context.Context, game Game) error
	GetPartners(ctx context.Context) ([]entity.Pokemon, error)
}
