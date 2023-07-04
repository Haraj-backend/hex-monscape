package battle

import (
	"context"
	"testing"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	// define mock dependencies
	gameStorage := newMockGameStorage()
	battleStorage := newMockBattleStorage()
	monsterStorage := newMockMonsterStorage()

	// define function for validating new game instance
	validateFunc := func(t *testing.T, svc Service, cfg ServiceConfig) {
		assert.Equal(t, cfg.GameStorage, svc.(*service).gameStorage)
		assert.Equal(t, cfg.BattleStorage, svc.(*service).battleStorage)
		assert.Equal(t, cfg.MonsterStorage, svc.(*service).monsterStorage)
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
				MonsterStorage: monsterStorage,
			},
			IsError: true,
		},
		{
			Name: "Test Missing Monster Storage",
			Config: ServiceConfig{
				GameStorage:    gameStorage,
				BattleStorage:  battleStorage,
				MonsterStorage: nil,
			},
			IsError: true,
		},
		{
			Name: "Test Valid Config",
			Config: ServiceConfig{
				GameStorage:    gameStorage,
				BattleStorage:  battleStorage,
				MonsterStorage: monsterStorage,
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
	partner := testutil.NewTestMonster()
	enemy := testutil.NewTestMonster()
	game, _ := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	})

	testCases := []struct {
		Name    string
		Enemy   *entity.Monster
		Game    *entity.Game
		IsError bool
	}{
		{
			Name:    "Test Get entity.Battle Valid",
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
			monsterStorage := newMockMonsterStorage()
			if testCase.Enemy != nil {
				monsterStorage.enemyMap[testCase.Enemy.ID] = *testCase.Enemy
			}
			gameID := ""

			if testCase.Game != nil {
				_ = gameStorage.SaveGame(context.Background(), *testCase.Game)
				gameID = testCase.Game.ID
			}

			svc, _ := NewService(ServiceConfig{
				GameStorage:    gameStorage,
				BattleStorage:  battleStorage,
				MonsterStorage: monsterStorage,
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
	partner := testutil.NewTestMonster()
	game, _ := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	})

	battle, _ := entity.NewBattle(entity.BattleConfig{
		GameID:  game.ID,
		Partner: partner,
		Enemy:   testutil.NewTestMonster(),
	})

	testCases := []struct {
		Name    string
		Battle  *entity.Battle
		Game    *entity.Game
		IsError bool
	}{
		{
			Name:    "Test Get entity.Battle Valid",
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
			Name:    "Test entity.Battle Not Found",
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
			monsterStorage := newMockMonsterStorage()
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
				MonsterStorage: monsterStorage,
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
	monsterStorage := newMockMonsterStorage()

	partner := testutil.NewTestMonster()
	game, err := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	})
	if err != nil {
		t.Fatalf("unable to init new game, due: %v", err)
	}

	battle, _ := entity.NewBattle(entity.BattleConfig{
		GameID:  game.ID,
		Partner: partner,
		Enemy:   testutil.NewTestMonster(),
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
		MonsterStorage: monsterStorage,
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
}

func TestServiceAttack(t *testing.T) {
	battleStorage := newMockBattleStorage()
	gameStorage := newMockGameStorage()
	monsterStorage := newMockMonsterStorage()

	partner := testutil.NewTestMonster()
	game, err := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	})
	if err != nil {
		t.Fatalf("unable to init new game, due: %v", err)
	}

	battle, _ := entity.NewBattle(entity.BattleConfig{
		GameID:  game.ID,
		Partner: partner,
		Enemy:   testutil.NewTestMonster(),
	})
	battle.State = entity.StatePartnerTurn

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
		MonsterStorage: monsterStorage,
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
	monsterStorage := newMockMonsterStorage()

	partner := testutil.NewTestMonster()
	game, err := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    partner,
		CreatedAt:  time.Now().Unix(),
	})
	if err != nil {
		t.Fatalf("unable to init new game, due: %v", err)
	}

	battle, _ := entity.NewBattle(entity.BattleConfig{
		GameID:  game.ID,
		Partner: partner,
		Enemy:   testutil.NewTestMonster(),
	})
	battle.State = entity.StatePartnerTurn

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
		MonsterStorage: monsterStorage,
	})
	if err != nil {
		t.Fatalf("unable to init new service, due: %v", err)
	}
	surrenderBattle, err := svc.Surrender(context.Background(), battle.GameID)
	if err != nil {
		t.Fatalf("unable to surrender, due: %v", err)
	}
	expectedBattle := battle
	expectedBattle.State = entity.StateLose
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
	data map[string]entity.Battle
}

func (gs *mockBattleStorage) GetBattle(ctx context.Context, gameID string) (*entity.Battle, error) {
	battle, ok := gs.data[gameID]
	if !ok {
		return nil, nil
	}
	return &battle, nil
}

func (gs *mockBattleStorage) SaveBattle(ctx context.Context, b entity.Battle) error {
	gs.data[b.GameID] = b
	return nil
}

func newMockBattleStorage() *mockBattleStorage {
	return &mockBattleStorage{
		data: map[string]entity.Battle{},
	}
}

type mockMockStorage struct {
	enemyMap map[string]entity.Monster
}

func (gs *mockMockStorage) GetPossibleEnemies(ctx context.Context) ([]entity.Monster, error) {
	var enemies []entity.Monster
	for _, enemy := range gs.enemyMap {
		enemies = append(enemies, enemy)
	}
	return enemies, nil
}

func newMockMonsterStorage() *mockMockStorage {
	return &mockMockStorage{
		enemyMap: map[string]entity.Monster{},
	}
}
