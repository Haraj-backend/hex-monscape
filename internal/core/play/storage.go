package play

import (
	"context"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
)

type GameStorage interface {
	// GetGame returns game instance for given gameID from storage. Returns nil
	// when given gameID is not found in database.
	GetGame(ctx context.Context, gameID string) (*entity.Game, error)

	// Save is used for saving game instance in storage.
	SaveGame(ctx context.Context, game entity.Game) error
}

type PartnerStorage interface {
	// GetAvailablePartners returns list of pokemon that selectable as partner.
	// Returns nil when there is no partners available.
	GetAvailablePartners(ctx context.Context) ([]entity.Pokemon, error)

	// GetPartner returns partner instance from given partner id. Returns nil
	// when given partnerID is not found.
	GetPartner(ctx context.Context, partnerID string) (*entity.Pokemon, error)
}
