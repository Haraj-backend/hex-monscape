package battle

import (
	"context"
	"testing"
	"time"

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
	gameStorage := newMockGameStorage()
	battleStorage := newMockBattleStorage()
	pokemonStorage := newMockPokemonStorage()

	enemy := newSamplePokemon()
	pokemonStorage.enemyMap[enemy.ID] = *enemy

	partner := newSamplePokemon()
	game, err := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	})
	if err != nil {
		t.Fatalf("unable to init new game, due: %v", err)
	}

	err = gameStorage.SaveGame(context.Background(), *game)
	if err != nil {
		t.Fatalf("unable to save game, due: %v", err)
	}

	svc, err := NewService(ServiceConfig{
		GameStorage:    gameStorage,
		BattleStorage:  battleStorage,
		PokemonStorage: pokemonStorage,
	})
	if err != nil {
		t.Fatalf("unable to init new service, due: %v", err)
	}
	battle, err := svc.StartBattle(context.Background(), game.ID)
	if err != nil {
		t.Fatalf("unable to start battle, due: %v", err)
	}
	assert.Equal(t, game.ID, battle.GameID, "gameID is not valid")
	assert.Equal(t, enemy.ID, battle.Enemy.ID, "enemyID is not valid")
}

func TestServiceGetBattle(t *testing.T) {
	battleStorage := newMockBattleStorage()
	gameStorage := newMockGameStorage()
	pokemonStorage := newMockPokemonStorage()

	partner := newSamplePokemon()
	game, err := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	})
	if err != nil {
		t.Fatalf("unable to init new game, due: %v", err)
	}

	battle, _ := NewBattle(BattleConfig{
		GameID:  game.ID,
		Partner: partner,
		Enemy:   newSamplePokemon(),
	})

	err = gameStorage.SaveGame(context.Background(), *game)
	if err != nil {
		t.Fatalf("unable to save game, due: %v", err)
	}
	err = battleStorage.SaveBattle(context.Background(), *battle)
	if err != nil {
		t.Fatalf("unable to save battle, due: %v", err)
	}

	svc, err := NewService(ServiceConfig{
		GameStorage:    gameStorage,
		BattleStorage:  battleStorage,
		PokemonStorage: pokemonStorage,
	})
	if err != nil {
		t.Fatalf("unable to init new service, due: %v", err)
	}
	newBattle, err := svc.GetBattle(context.Background(), battle.GameID)
	if err != nil {
		t.Fatalf("unable to get battle, due: %v", err)
	}
	assert.Equal(t, battle, newBattle, "battle is not valid")
}

func TestServiceDecideTurn(t *testing.T) {
	battleStorage := newMockBattleStorage()
	gameStorage := newMockGameStorage()
	pokemonStorage := newMockPokemonStorage()

	partner := newSamplePokemon()
	game, err := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	})
	if err != nil {
		t.Fatalf("unable to init new game, due: %v", err)
	}

	battle, _ := NewBattle(BattleConfig{
		GameID:  game.ID,
		Partner: partner,
		Enemy:   newSamplePokemon(),
	})

	err = gameStorage.SaveGame(context.Background(), *game)
	if err != nil {
		t.Fatalf("unable to save game, due: %v", err)
	}
	err = battleStorage.SaveBattle(context.Background(), *battle)
	if err != nil {
		t.Fatalf("unable to save battle, due: %v", err)
	}

	svc, err := NewService(ServiceConfig{
		GameStorage:    gameStorage,
		BattleStorage:  battleStorage,
		PokemonStorage: pokemonStorage,
	})
	if err != nil {
		t.Fatalf("unable to init new service, due: %v", err)
	}
	_, err = svc.DecideTurn(context.Background(), battle.GameID)
	if err != nil {
		t.Fatalf("unable to decide turn, due: %v", err)
	}
}

func TestServiceAttack(t *testing.T) {
	battleStorage := newMockBattleStorage()
	gameStorage := newMockGameStorage()
	pokemonStorage := newMockPokemonStorage()

	partner := newSamplePokemon()
	game, err := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	})
	if err != nil {
		t.Fatalf("unable to init new game, due: %v", err)
	}

	battle, _ := NewBattle(BattleConfig{
		GameID:  game.ID,
		Partner: partner,
		Enemy:   newSamplePokemon(),
	})
	battle.State = PARTNER_TURN

	err = gameStorage.SaveGame(context.Background(), *game)
	if err != nil {
		t.Fatalf("unable to save game, due: %v", err)
	}
	err = battleStorage.SaveBattle(context.Background(), *battle)
	if err != nil {
		t.Fatalf("unable to save battle, due: %v", err)
	}

	svc, err := NewService(ServiceConfig{
		GameStorage:    gameStorage,
		BattleStorage:  battleStorage,
		PokemonStorage: pokemonStorage,
	})
	if err != nil {
		t.Fatalf("unable to init new service, due: %v", err)
	}
	_, err = svc.Attack(context.Background(), battle.GameID)
	if err != nil {
		t.Fatalf("unable to attack, due: %v", err)
	}
}

func TestServiceSurrender(t *testing.T) {
	battleStorage := newMockBattleStorage()
	gameStorage := newMockGameStorage()
	pokemonStorage := newMockPokemonStorage()

	partner := newSamplePokemon()
	game, err := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	})
	if err != nil {
		t.Fatalf("unable to init new game, due: %v", err)
	}

	battle, _ := NewBattle(BattleConfig{
		GameID:  game.ID,
		Partner: partner,
		Enemy:   newSamplePokemon(),
	})
	battle.State = PARTNER_TURN

	err = gameStorage.SaveGame(context.Background(), *game)
	if err != nil {
		t.Fatalf("unable to save game, due: %v", err)
	}
	err = battleStorage.SaveBattle(context.Background(), *battle)
	if err != nil {
		t.Fatalf("unable to save battle, due: %v", err)
	}

	svc, err := NewService(ServiceConfig{
		GameStorage:    gameStorage,
		BattleStorage:  battleStorage,
		PokemonStorage: pokemonStorage,
	})
	if err != nil {
		t.Fatalf("unable to init new service, due: %v", err)
	}
	_, err = svc.Surrender(context.Background(), battle.GameID)
	if err != nil {
		t.Fatalf("unable to surrender, due: %v", err)
	}
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
