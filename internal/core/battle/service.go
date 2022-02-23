package battle

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"gopkg.in/validator.v2"
)

var (
	ErrGameNotFound       = errors.New("game is not found")
	ErrBattleNotFound     = errors.New("battle is not found")
	ErrInvalidBattleState = errors.New("invalid battle state")
)

type Service interface {
	// StartBattle is used for starting new battle for given game id. Battle could only
	// be started when there is no previous battle or it is already ended.
	StartBattle(ctx context.Context, gameID string) (*Battle, error)

	// GetBattle returns ongoing battle for given game id.
	GetBattle(ctx context.Context, gameID string) (*Battle, error)

	// DecideTurn is used for deciding turn for the battle. There are 3 possible outcome
	// battle states from this action:
	//
	// - DECIDE_TURN => partner has been attacked by enemy but still not lose, so the state
	// 					returned back to DECIDE_TURN
	// - LOSE => partner has been attacked by enemy and already lose
	// - PARTNER_TURN => it is partner turn, client may commence attack or surrender
	DecideTurn(ctx context.Context, gameID string) (*Battle, error)

	// Attack is used for executing attack to enemy. The possible battle state outcome
	// is DECIDE_TURN or WIN.
	Attack(ctx context.Context, gameID string) (*Battle, error)

	// Surrender is used for executing surrender action. Battle will immediately ended
	// with player losing the battle.
	Surrender(ctx context.Context, gameID string) (*Battle, error)
}

type service struct {
	gameStorage    GameStorage
	battleStorage  BattleStorage
	pokemonStorage PokemonStorage
}

func (s *service) StartBattle(ctx context.Context, gameID string) (*Battle, error) {
	// get existing battle, make sure there is no ongoing battle
	game, battle, err := s.getGameAndBattleInstance(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("unable to get game and battle instance due: %w", err)
	}
	if game == nil {
		return nil, ErrGameNotFound
	}
	if battle != nil {
		return nil, ErrInvalidBattleState
	}
	// get possible enemies, choose it randomly
	enemies, err := s.pokemonStorage.GetPossibleEnemies(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get possible enemies due: %w", err)
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	enemy := enemies[rnd.Intn(len(enemies))]
	// create new battle instance
	battle, err = NewBattle(BattleConfig{
		GameID:  gameID,
		Partner: game.Partner,
		Enemy:   &enemy,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create new battle instance due: %w", err)
	}
	// save newly created battle instance to storage
	err = s.battleStorage.SaveBattle(ctx, *battle)
	if err != nil {
		return nil, fmt.Errorf("unable to save battle into storage due: %w", err)
	}
	return battle, nil
}

func (s *service) GetBattle(ctx context.Context, gameID string) (*Battle, error) {
	game, battle, err := s.getGameAndBattleInstance(ctx, gameID)
	if game == nil {
		return nil, ErrGameNotFound
	}
	if battle == nil {
		return nil, ErrBattleNotFound
	}
	return battle, err
}

// getGameAndBattleInstance returns game & battle for given game id, if game is not found both game
// and battle will be returned nil. If battle is not found, the battle will be returned nil.
func (s *service) getGameAndBattleInstance(ctx context.Context, gameID string) (*entity.Game, *Battle, error) {
	game, err := s.gameStorage.GetGame(ctx, gameID)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get game due: %w", err)
	}
	if game == nil {
		return nil, nil, nil
	}
	battle, err := s.battleStorage.GetBattle(ctx, gameID)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get battle due: %w", err)
	}
	if battle != nil && battle.IsEnded() {
		battle = nil
	}
	return game, battle, nil
}

func (s *service) DecideTurn(ctx context.Context, gameID string) (*Battle, error) {
	battle, err := s.GetBattle(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if battle.State != DECIDE_TURN {
		return nil, ErrInvalidBattleState
	}
	_, err = battle.DecideTurn()
	if err != nil {
		return nil, fmt.Errorf("unable to decide turn due: %w", err)
	}
	if battle.State == ENEMY_TURN {
		err = battle.EnemyAttack()
		if err != nil {
			return nil, fmt.Errorf("unable to make enemy attack due: %w", err)
		}
	}
	err = s.battleStorage.SaveBattle(ctx, *battle)
	if err != nil {
		return nil, fmt.Errorf("unable to save battle due: %w", err)
	}
	return battle, nil
}

func (s *service) Attack(ctx context.Context, gameID string) (*Battle, error) {
	game, battle, err := s.getGameAndBattleInstance(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("unable to get game and battle instance due: %w", err)
	}
	if game == nil {
		return nil, ErrGameNotFound
	}
	if battle == nil {
		return nil, ErrBattleNotFound
	}
	if battle.State != PARTNER_TURN {
		return nil, ErrInvalidBattleState
	}
	err = battle.PartnerAttack()
	if err != nil {
		return nil, fmt.Errorf("unable to decide turn due: %w", err)
	}
	err = s.battleStorage.SaveBattle(ctx, *battle)
	if err != nil {
		return nil, fmt.Errorf("unable to save battle due: %w", err)
	}
	if battle.State == WIN {
		game.IncBattleWon()
		err = s.gameStorage.SaveGame(ctx, *game)
		if err != nil {
			return nil, fmt.Errorf("unable to save game due: %w", err)
		}
	}
	return battle, nil
}

func (s *service) Surrender(ctx context.Context, gameID string) (*Battle, error) {
	battle, err := s.GetBattle(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if battle.State != PARTNER_TURN {
		return nil, ErrInvalidBattleState
	}
	err = battle.PartnerSurrender()
	if err != nil {
		return nil, fmt.Errorf("unable to decide turn due: %w", err)
	}
	err = s.battleStorage.SaveBattle(ctx, *battle)
	if err != nil {
		return nil, fmt.Errorf("unable to save battle due: %w", err)
	}
	return battle, nil
}

type ServiceConfig struct {
	GameStorage    GameStorage    `validate:"nonnil"`
	BattleStorage  BattleStorage  `validate:"nonnil"`
	PokemonStorage PokemonStorage `validate:"nonnil"`
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
	svc := &service{
		gameStorage:    cfg.GameStorage,
		battleStorage:  cfg.BattleStorage,
		pokemonStorage: cfg.PokemonStorage,
	}
	return svc, nil
}
