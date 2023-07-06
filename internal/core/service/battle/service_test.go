package battle_test

/*
	The purpose of testing the Service component is to ensure it has correct
	implementation of business logic.

	The common pitfall when creating test for Service component is we tend to use
	concrete implementation for the dependency components (e.g actual GameStorage
	for MySQL). Not only this will increase the test complexity but also it will
	increase the possibility of getting false test result. The reason is simply
	because service such as MySQL has its own constraints & has much higher chance
	of failing rather than its mock counterpart (e.g network failure).

	So to avoid this pitfall, our first go to choice is to use mock implementation
	for the dependency when testing the Service component. This way we can control
	more the behavior of the dependency components to fit our test scenarios.
*/

import (
	"context"
	"testing"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/battle"
	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	// define mock dependencies
	gameStorage := newMockGameStorage()
	battleStorage := newMockBattleStorage()
	monsterStorage := newMockMonsterStorage()

	// define test cases
	testCases := []struct {
		Name    string
		Config  battle.ServiceConfig
		IsError bool
	}{
		{
			Name: "Test Missing Game Storage",
			Config: battle.ServiceConfig{
				GameStorage:    nil,
				BattleStorage:  battleStorage,
				MonsterStorage: monsterStorage,
			},
			IsError: true,
		},
		{
			Name: "Test Missing Monster Storage",
			Config: battle.ServiceConfig{
				GameStorage:    gameStorage,
				BattleStorage:  battleStorage,
				MonsterStorage: nil,
			},
			IsError: true,
		},
		{
			Name: "Test Valid Config",
			Config: battle.ServiceConfig{
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
			_, err := battle.NewService(testcase.Config)
			require.Equal(t, testcase.IsError, (err != nil), "unexpected error")
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

			svc, _ := battle.NewService(battle.ServiceConfig{
				GameStorage:    gameStorage,
				BattleStorage:  battleStorage,
				MonsterStorage: monsterStorage,
			})
			battle, err := svc.StartBattle(context.Background(), gameID)
			require.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if err == nil {
				require.Equal(t, game.ID, battle.GameID, "gameID is not valid")
				require.Equal(t, enemy.ID, battle.Enemy.ID, "enemyID is not valid")

				storedBattle, err := battleStorage.GetBattle(context.Background(), gameID)
				require.NoError(t, err, "unable to get stored battle")
				require.Equal(t, battle, storedBattle, "battle stored is not valid")
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

	bt, _ := entity.NewBattle(entity.BattleConfig{
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
			Battle:  bt,
			Game:    game,
			IsError: false,
		},
		{
			Name:    "Test Game Not Found",
			Battle:  bt,
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

			svc, _ := battle.NewService(battle.ServiceConfig{
				GameStorage:    gameStorage,
				BattleStorage:  battleStorage,
				MonsterStorage: monsterStorage,
			})
			newBattle, err := svc.GetBattle(context.Background(), gameID)
			require.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if err == nil {
				require.Equal(t, bt, newBattle, "battle is not valid")
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

	bt, _ := entity.NewBattle(entity.BattleConfig{
		GameID:  game.ID,
		Partner: partner,
		Enemy:   testutil.NewTestMonster(),
	})

	err = gameStorage.SaveGame(context.Background(), *game)
	if err != nil {
		t.Fatalf("unable to save game, due: %v", err)
	}
	err = battleStorage.SaveBattle(context.Background(), *bt)
	if err != nil {
		t.Fatalf("unable to save battle, due: %v", err)
	}

	svc, err := battle.NewService(battle.ServiceConfig{
		GameStorage:    gameStorage,
		BattleStorage:  battleStorage,
		MonsterStorage: monsterStorage,
	})
	if err != nil {
		t.Fatalf("unable to init new service, due: %v", err)
	}
	bt, err = svc.DecideTurn(context.Background(), bt.GameID)
	if err != nil {
		t.Fatalf("unable to decide turn, due: %v", err)
	}

	storedBattle, err := battleStorage.GetBattle(context.Background(), bt.GameID)
	require.NoError(t, err, "unable to get stored battle")
	require.Equal(t, bt, storedBattle, "invalid battle stored")
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

	bt, _ := entity.NewBattle(entity.BattleConfig{
		GameID:  game.ID,
		Partner: partner,
		Enemy:   testutil.NewTestMonster(),
	})
	bt.State = entity.StatePartnerTurn

	err = gameStorage.SaveGame(context.Background(), *game)
	if err != nil {
		t.Fatalf("unable to save game, due: %v", err)
	}
	err = battleStorage.SaveBattle(context.Background(), *bt)
	if err != nil {
		t.Fatalf("unable to save battle, due: %v", err)
	}

	svc, err := battle.NewService(battle.ServiceConfig{
		GameStorage:    gameStorage,
		BattleStorage:  battleStorage,
		MonsterStorage: monsterStorage,
	})
	if err != nil {
		t.Fatalf("unable to init new service, due: %v", err)
	}
	_, err = svc.Attack(context.Background(), bt.GameID)
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

	bt, _ := entity.NewBattle(entity.BattleConfig{
		GameID:  game.ID,
		Partner: partner,
		Enemy:   testutil.NewTestMonster(),
	})
	bt.State = entity.StatePartnerTurn

	err = gameStorage.SaveGame(context.Background(), *game)
	if err != nil {
		t.Fatalf("unable to save game, due: %v", err)
	}
	err = battleStorage.SaveBattle(context.Background(), *bt)
	if err != nil {
		t.Fatalf("unable to save battle, due: %v", err)
	}

	svc, err := battle.NewService(battle.ServiceConfig{
		GameStorage:    gameStorage,
		BattleStorage:  battleStorage,
		MonsterStorage: monsterStorage,
	})
	if err != nil {
		t.Fatalf("unable to init new service, due: %v", err)
	}
	surrenderBattle, err := svc.Surrender(context.Background(), bt.GameID)
	if err != nil {
		t.Fatalf("unable to surrender, due: %v", err)
	}
	expectedBattle := bt
	expectedBattle.State = entity.StateLose
	require.Equal(t, expectedBattle, surrenderBattle, "invalid battle stored")
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
