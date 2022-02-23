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
	partner := newSamplePokemon()
	enemy := newSamplePokemon()
	game, _ := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	})

	testCases := []struct {
		Name    string
		Enemy   *entity.Pokemon
		Game    *entity.Game
		IsError bool
	}{
		{
			Name:    "Test Get Battle Valid",
			Enemy:   enemy,
			Game:    game,
			IsError: false,
		},
		{
			Name:    "Test Game Not Found",
			Enemy:   enemy,
			Game:    nil,
			IsError: true,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battleStorage := newMockBattleStorage()
			gameStorage := newMockGameStorage()
			pokemonStorage := newMockPokemonStorage()
			if testCase.Enemy != nil {
				pokemonStorage.enemyMap[testCase.Enemy.ID] = *testCase.Enemy
			}
			gameID := ""

			if testCase.Game != nil {
				_ = gameStorage.SaveGame(context.Background(), *testCase.Game)
				gameID = testCase.Game.ID
			}

			svc, _ := NewService(ServiceConfig{
				GameStorage:    gameStorage,
				BattleStorage:  battleStorage,
				PokemonStorage: pokemonStorage,
			})
			battle, err := svc.StartBattle(context.Background(), gameID)
			assert.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if err == nil {
				assert.Equal(t, game.ID, battle.GameID, "gameID is not valid")
				assert.Equal(t, enemy.ID, battle.Enemy.ID, "enemyID is not valid")

				storedBattle, err := battleStorage.GetBattle(context.Background(), gameID)
				assert.NoError(t, err, "unable to get stored battle")
				assert.Equal(t, battle, storedBattle, "battle stored is not valid")
			}
		})
	}
}

func TestServiceGetBattle(t *testing.T) {
	partner := newSamplePokemon()
	game, _ := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	})

	battle, _ := NewBattle(BattleConfig{
		GameID:  game.ID,
		Partner: partner,
		Enemy:   newSamplePokemon(),
	})

	testCases := []struct {
		Name    string
		Battle  *Battle
		Game    *entity.Game
		IsError bool
	}{
		{
			Name:    "Test Get Battle Valid",
			Battle:  battle,
			Game:    game,
			IsError: false,
		},
		{
			Name:    "Test Game Not Found",
			Battle:  battle,
			Game:    nil,
			IsError: true,
		},
		{
			Name:    "Test Battle Not Found",
			Battle:  nil,
			Game:    game,
			IsError: true,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battleStorage := newMockBattleStorage()
			gameStorage := newMockGameStorage()
			pokemonStorage := newMockPokemonStorage()
			gameID := ""

			if testCase.Game != nil {
				_ = gameStorage.SaveGame(context.Background(), *testCase.Game)
				gameID = testCase.Game.ID
			}
			if testCase.Battle != nil {
				_ = battleStorage.SaveBattle(context.Background(), *testCase.Battle)
				gameID = testCase.Battle.GameID
			}

			svc, _ := NewService(ServiceConfig{
				GameStorage:    gameStorage,
				BattleStorage:  battleStorage,
				PokemonStorage: pokemonStorage,
			})
			newBattle, err := svc.GetBattle(context.Background(), gameID)
			assert.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if err == nil {
				assert.Equal(t, battle, newBattle, "battle is not valid")
			}
		})
	}
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
	battle, err = svc.DecideTurn(context.Background(), battle.GameID)
	if err != nil {
		t.Fatalf("unable to decide turn, due: %v", err)
	}

	storedBattle, err := battleStorage.GetBattle(context.Background(), battle.GameID)
	assert.NoError(t, err, "unable to get stored battle")
	assert.Equal(t, battle, storedBattle, "invalid battle stored")
	assert.Equal(t, DECIDE_TURN, storedBattle.State, "invalid battle state")
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
	surrenderBattle, err := svc.Surrender(context.Background(), battle.GameID)
	if err != nil {
		t.Fatalf("unable to surrender, due: %v", err)
	}
	expectedBattle := battle
	expectedBattle.State = LOSE
	assert.Equal(t, expectedBattle, surrenderBattle, "invalid battle stored")
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
