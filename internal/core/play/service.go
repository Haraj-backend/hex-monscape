package play

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/shared/telemetry"
	"gopkg.in/validator.v2"
)

var (
	ErrGameNotFound    = errors.New("game is not found")
	ErrPartnerNotFound = errors.New("partner is not found")
)

type Service interface {
	// GetAvailablePartners returns pokemons that available to be selected as player partner.
	GetAvailablePartners(ctx context.Context) ([]entity.Pokemon, error)

	// NewGame is used for initiating new game in storage. If the given `partnerID` not found in storage,
	// it returns `ErrPartnerNotFound`. Upon success it returns game instance that being saved on storage.
	NewGame(ctx context.Context, playerName string, partnerID string) (*entity.Game, error)

	// GetGame returns game instance from storage from given game id. Upon game is not found, it returns
	// `ErrGameNotFound`.
	GetGame(ctx context.Context, gameID string) (*entity.Game, error)
}

type service struct {
	gameStorage    GameStorage
	partnerStorage PartnerStorage
}

func (s *service) GetAvailablePartners(ctx context.Context) ([]entity.Pokemon, error) {
	tr := telemetry.GetTracer()
	span := tr.Trace(ctx, "GetAvailablePartners")
	defer span.End()

	partners, err := s.partnerStorage.GetAvailablePartners(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get available partners due: %w", err)
	}
	return partners, nil
}

func (s *service) NewGame(ctx context.Context, playerName string, partnerID string) (*entity.Game, error) {
	tr := telemetry.GetTracer()
	span := tr.Trace(ctx, "NewGame")
	defer span.End()

	// get partner instance
	partner, err := s.partnerStorage.GetPartner(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch partner instance due: %w", err)
	}
	if partner == nil {
		return nil, ErrPartnerNotFound
	}
	// initiate new game instance
	cfg := entity.GameConfig{
		PlayerName: playerName,
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	}
	game, err := entity.NewGame(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize game instance due: %w", err)
	}
	// store the game instance on storage
	err = s.gameStorage.SaveGame(ctx, *game)
	if err != nil {
		return nil, fmt.Errorf("unable to save game instance due: %w", err)
	}
	return game, nil
}

func (s *service) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	tr := telemetry.GetTracer()
	span := tr.Trace(ctx, "GetGame")
	defer span.End()

	// get game instance from storage
	game, err := s.gameStorage.GetGame(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("unable to get game instance due: %w", err)
	}
	if game == nil {
		return nil, ErrGameNotFound
	}
	return game, nil
}

type ServiceConfig struct {
	GameStorage    GameStorage    `validate:"nonnil"`
	PartnerStorage PartnerStorage `validate:"nonnil"`
}

func (c ServiceConfig) Validate() error {
	return validator.Validate(c)
}

// NewService returns new instance of service.
func NewService(cfg ServiceConfig) (Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &service{
		gameStorage:    cfg.GameStorage,
		partnerStorage: cfg.PartnerStorage,
	}
	return s, nil
}
