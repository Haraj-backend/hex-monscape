package pokemon

import (
	"context"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
)

type PokemonStorage struct{}

func (ps *PokemonStorage) GetAvailablePartners(ctx context.Context) ([]entity.Pokemon, error) {
	// TODO
	return nil, nil
}

func (ps *PokemonStorage) GetPartner(ctx context.Context, partnerID string) (*entity.Pokemon, error) {
	// TODO
	return nil, nil
}

func (ps *PokemonStorage) GetPossibleEnemies(ctx context.Context) ([]entity.Pokemon, error) {
	// TODO
	return nil, nil
}
