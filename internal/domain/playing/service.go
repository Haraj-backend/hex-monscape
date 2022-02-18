package playing

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/riandyrn/pokebattle/internal/domain/entity"
	"gopkg.in/validator.v2"
)

var (
	ErrGameNotFound    = errors.New("game is not found")
	ErrPartnerNotFound = errors.New("partner is not found")
)

type ServiceConfig struct {
	GameStorage    GameStorage    `validate:"nonnil"`
	PartnerStorage PartnerStorage `validate:"nonnil"`
}

func (c ServiceConfig) Validate() error {
	return validator.Validate(c)
}

type Service struct {
	gameStorage    GameStorage
	partnerStorage PartnerStorage
}

// GetAvailablePartners returns pokemons that available to be selected as player partner.
func (s *Service) GetAvailablePartners(ctx context.Context) ([]entity.Pokemon, error) {
	partners, err := s.partnerStorage.GetAvailablePartners(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get available partners due: %w", err)
	}
	return partners, nil
}

// NewGame is used for initiating new game in storage. If the given `partnerID` not found in storage,
// it returns `ErrPartnerNotFound`. Upon success it returns game instance that being saved on storage.
func (s *Service) NewGame(ctx context.Context, playerName string, partnerID string) (*Game, error) {
	// get partner instance
	partner, err := s.partnerStorage.GetPartner(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch partner instance due: %w", err)
	}
	if partner == nil {
		return nil, ErrPartnerNotFound
	}
	// initiate new game instance
	cfg := GameConfig{
		PlayerName: playerName,
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	}
	game, err := NewGame(cfg)
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

// GetGame returns game instance from storage from given game id. Upon game is not found, it returns
// `ErrGameNotFound`.
func (s *Service) GetGame(ctx context.Context, gameID string) (*Game, error) {
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

// AdvanceScenario is used for advancing scenario for given game id. Returns the advanced
// game instance upon success. When given game id is not found returns `ErrGameNotFound`.
func (s *Service) AdvanceScenario(ctx context.Context, gameID string) (*Game, error) {
	// get game instance
	game, err := s.GetGame(ctx, gameID)
	if err != nil {
		return nil, err
	}
	// advance game scenario
	game.AdvanceScenario()
	// save game instance
	err = s.gameStorage.SaveGame(ctx, *game)
	if err != nil {
		return nil, fmt.Errorf("unable to save game to storage due: %w", err)
	}
	return game, nil
}

func NewService(cfg ServiceConfig) (*Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &Service{
		gameStorage:    cfg.GameStorage,
		partnerStorage: cfg.PartnerStorage,
	}
	return s, nil
}
