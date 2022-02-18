package battling

import (
	"context"
	"errors"

	"gopkg.in/validator.v2"
)

var (
	ErrGameNotFound       = errors.New("game is not found")
	ErrBattleNotFound     = errors.New("battle is not found")
	ErrInvalidBattleState = errors.New("invalid battle state")
)

type Service struct {
	gameStorage    GameStorage
	battleStorage  BattleStorage
	pokemonStorage PokemonStorage
}

// StartBattle is used for starting new battle for given game id. Battle could only
// be started when previous battle is not found or already ended.
func (s *Service) StartBattle(ctx context.Context, gameID string) (*Battle, error) {
	// TODO
	return nil, nil
}

// GetBattle returns battle for given game id.
func (s *Service) GetBattle(ctx context.Context, gameID string) (*Battle, error) {
	// TODO
	return nil, nil
}

// DecideTurn is used for deciding turn for the battle. There are 3 possible outcome
// battle states from this action:
//
// - DECIDE_TURN => partner has been attacked by enemy but still not lose, so the state
// 					returned back to DECIDE_TURN
// - LOSE => partner has been attacked by enemy and already lose
// - PARTNER_TURN => it is partner turn, client may commence attack or surrender
func (s *Service) DecideTurn(ctx context.Context, gameID string) (*Battle, error) {
	// TODO
	return nil, nil
}

// Attack is used for executing attack to enemy. The possible battle state outcome
// is DECIDE_TURN or WIN.
func (s *Service) Attack(ctx context.Context, gameID string) (*Battle, error) {
	// TODO
	return nil, nil
}

// Surrender is used for executing surrender action. Battle will immediately ended
// with player losing the battle.
func (s *Service) Surrender(ctx context.Context, gameID string) (*Battle, error) {
	// TODO
	return nil, nil
}

type ServiceConfig struct {
	GameStorage    GameStorage    `validator:"nonnil"`
	BattleStorage  BattleStorage  `validator:"nonnil"`
	PokemonStorage PokemonStorage `validator:"nonnil"`
}

func (c ServiceConfig) Validate() error {
	return validator.Validate(c)
}

// NewService returns new instance of Service.
func NewService(cfg ServiceConfig) (*Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	svc := &Service{
		gameStorage:    cfg.GameStorage,
		battleStorage:  cfg.BattleStorage,
		pokemonStorage: cfg.PokemonStorage,
	}
	return svc, nil
}
