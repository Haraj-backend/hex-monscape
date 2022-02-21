package battling

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

type Service struct {
	gameStorage    GameStorage
	battleStorage  BattleStorage
	pokemonStorage PokemonStorage
}

// StartBattle is used for starting new battle for given game id. Battle could only
// be started when there is no previous battle or it is already ended.
func (s *Service) StartBattle(ctx context.Context, gameID string) (*Battle, error) {
	// get existing battle, make sure there is no ongoing battle
	game, battle, err := s.getBattle(ctx, gameID)
	if err != nil && !errors.Is(err, ErrBattleNotFound) {
		return nil, err
	}
	if err != nil && game == nil {
		return nil, fmt.Errorf("unable to start battle due: %w", err)
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

// GetBattle returns ongoing battle for given game id.
func (s *Service) GetBattle(ctx context.Context, gameID string) (*Battle, error) {
	_, battle, err := s.getBattle(ctx, gameID)
	return battle, err
}

func (s *Service) getBattle(ctx context.Context, gameID string) (*entity.Game, *Battle, error) {
	game, err := s.gameStorage.GetGame(ctx, gameID)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get game due: %w", err)
	}
	if game == nil {
		return nil, nil, ErrGameNotFound
	}
	battle, err := s.battleStorage.GetBattle(ctx, gameID)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get battle due: %w", err)
	}
	if battle != nil && battle.IsEnded() {
		return nil, nil, ErrBattleNotFound
	}
	return game, battle, nil
}

// DecideTurn is used for deciding turn for the battle. There are 3 possible outcome
// battle states from this action:
//
// - DECIDE_TURN => partner has been attacked by enemy but still not lose, so the state
// 					returned back to DECIDE_TURN
// - LOSE => partner has been attacked by enemy and already lose
// - PARTNER_TURN => it is partner turn, client may commence attack or surrender
func (s *Service) DecideTurn(ctx context.Context, gameID string) (*Battle, error) {
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

// Attack is used for executing attack to enemy. The possible battle state outcome
// is DECIDE_TURN or WIN.
func (s *Service) Attack(ctx context.Context, gameID string) (*Battle, error) {
	game, battle, err := s.getBattle(ctx, gameID)
	if err != nil {
		return nil, err
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

// Surrender is used for executing surrender action. Battle will immediately ended
// with player losing the battle.
func (s *Service) Surrender(ctx context.Context, gameID string) (*Battle, error) {
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
