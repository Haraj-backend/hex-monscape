package battle

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	// define mock dependencies
	gameStorage := newMockGameStorage()
	battleStorage := newMockBattleStorage()
	pokemonStorage := newMockPokemonStorage()

	// define function for validating new game instance
	validateFunc := func(t *testing.T, svc Service, cfg ServiceConfig) {
		assert.Equal(t, cfg.GameStorage, svc.(*service).gameStorage)
		assert.Equal(t, cfg.BattleStorage, svc.(*service).battleStorage)
		assert.Equal(t, cfg.PokemonStorage, svc.(*service).pokemonStorage)
	}
	// define test cases
	testCases := []struct {
		Name    string
		Config  ServiceConfig
		IsError bool
	}{
		{
			Name: "Test Missing Game Storage",
			Config: ServiceConfig{
				GameStorage:    nil,
				BattleStorage:  battleStorage,
				PokemonStorage: pokemonStorage,
			},
			IsError: true,
		},
		{
			Name: "Test Missing Pokemon Storage",
			Config: ServiceConfig{
				GameStorage:    gameStorage,
				BattleStorage:  battleStorage,
				PokemonStorage: nil,
			},
			IsError: true,
		},
		{
			Name: "Test Valid Config",
			Config: ServiceConfig{
				GameStorage:    gameStorage,
				BattleStorage:  battleStorage,
				PokemonStorage: pokemonStorage,
			},
			IsError: false,
		},
	}
	// execute test cases
	for _, testcase := range testCases {
		t.Run(testcase.Name, func(t *testing.T) {
			svc, err := NewService(testcase.Config)
			assert.Equal(t, testcase.IsError, (err != nil), "unexpected error")
			if svc == nil {
				return
			}
			validateFunc(t, svc, testcase.Config)
		})
	}
}

func TestServiceStartBattle(t *testing.T) {
	// TODO
}

func TestServiceGetBattle(t *testing.T) {
	// TODO
}

func TestServiceDecideTurn(t *testing.T) {
	// TODO
}

func TestServiceAttack(t *testing.T) {
	// TODO
}

func TestServiceSurrender(t *testing.T) {
	// TODO
}

type mockGameStorage struct {
	data map[string]entity.Game
}

func (gs *mockGameStorage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	game, ok := gs.data[gameID]
	if !ok {
		return nil, nil
	}
	return &game, nil
}

func (gs *mockGameStorage) SaveGame(ctx context.Context, game entity.Game) error {
	gs.data[game.ID] = game
	return nil
}

func newMockGameStorage() *mockGameStorage {
	return &mockGameStorage{
		data: map[string]entity.Game{},
	}
}

type mockBattleStorage struct {
	data map[string]Battle
}

func (gs *mockBattleStorage) GetBattle(ctx context.Context, gameID string) (*Battle, error) {
	battle, ok := gs.data[gameID]
	if !ok {
		return nil, nil
	}
	return &battle, nil
}

func (gs *mockBattleStorage) SaveBattle(ctx context.Context, b Battle) error {
	gs.data[b.GameID] = b
	return nil
}

func newMockBattleStorage() *mockBattleStorage {
	return &mockBattleStorage{
		data: map[string]Battle{},
	}
}

type mockPokemonStorage struct {
	enemyMap map[string]entity.Pokemon
}

func (gs *mockPokemonStorage) GetPossibleEnemies(ctx context.Context) ([]entity.Pokemon, error) {
	var enemies []entity.Pokemon
	for _, enemy := range gs.enemyMap {
		enemies = append(enemies, enemy)
	}
	return enemies, nil
}

func newMockPokemonStorage() *mockPokemonStorage {
	return &mockPokemonStorage{
		enemyMap: map[string]entity.Pokemon{},
	}
}

func newNewService() Service {
	// initialize service
	cfg := ServiceConfig{
		GameStorage:    newMockGameStorage(),
		BattleStorage:  newMockBattleStorage(),
		PokemonStorage: newMockPokemonStorage(),
	}
	svc, _ := NewService(cfg)
	return svc
}
