package entity_test

import (
	"testing"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateGameConfig(t *testing.T) {
	sampleMonster := testutil.NewTestMonster()
	currentTs := time.Now().Unix()
	testCases := []struct {
		Name    string
		Config  entity.GameConfig
		IsError bool
	}{
		{
			Name: "Empty Player Name",
			Config: entity.GameConfig{
				PlayerName: "",
				Partner:    sampleMonster,
				CreatedAt:  currentTs,
			},
			IsError: true,
		},
		{
			Name: "Empty Partner",
			Config: entity.GameConfig{
				PlayerName: "Riandy R.N",
				Partner:    nil,
				CreatedAt:  currentTs,
			},
			IsError: true,
		},
		{
			Name: "Empty Created At",
			Config: entity.GameConfig{
				PlayerName: "Riandy R.N",
				Partner:    sampleMonster,
				CreatedAt:  0,
			},
			IsError: true,
		},
		{
			Name: "All Filled",
			Config: entity.GameConfig{
				PlayerName: "Riandy R.N",
				Partner:    sampleMonster,
				CreatedAt:  currentTs,
			},
			IsError: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := testCase.Config.Validate()
			assert.Equal(t, testCase.IsError, (err != nil), "unexpected error")
		})
	}
}

func TestNewGame(t *testing.T) {
	// define function for validating game
	validateGame := func(t *testing.T, game entity.Game, cfg entity.GameConfig) {
		assert.NotEmpty(t, game.ID, "game id is empty")
		assert.Equal(t, cfg.PlayerName, game.PlayerName, "mismatch player name")
		assert.Equal(t, cfg.Partner, game.Partner, "mismatch partner")
		assert.Equal(t, cfg.CreatedAt, game.CreatedAt, "mismatch created_at")
		assert.Equal(t, 0, game.BattleWon, "mismatch battle_won")
		assert.Equal(t, entity.ScenarioBattle1, game.Scenario, "mismatch scenario")
	}
	// define test cases
	testCases := []struct {
		Name    string
		Config  entity.GameConfig
		IsError bool
	}{
		{
			Name:    "Invalid Config",
			Config:  entity.GameConfig{},
			IsError: true,
		},
		{
			Name: "Valid Config",
			Config: entity.GameConfig{
				PlayerName: "Riandy R.N",
				Partner:    testutil.NewTestMonster(),
				CreatedAt:  time.Now().Unix(),
			},
			IsError: false,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			game, err := entity.NewGame(testCase.Config)
			assert.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if game == nil {
				return
			}
			validateGame(t, *game, testCase.Config)
		})
	}
}

func TestAdvanceScenario(t *testing.T) {
	// initialize new game
	game := initNewGame()
	// test scenario
	scenarios := []entity.Scenario{entity.ScenarioBattle1, entity.ScenarioBattle2, entity.ScenarioBattle3, entity.ScenarioEndGame}
	for i := 0; i < len(scenarios); i++ {
		// not won any battle, scenario should be same as previous
		game.AdvanceScenario()
		require.Equal(t, scenarios[i], game.Scenario, "scenario should not advancing")
		if i == len(scenarios)-1 {
			continue
		}
		// won battle, scenario should be advancing
		game.BattleWon++
		game.AdvanceScenario()
		require.Equal(t, scenarios[i+1], game.Scenario, "scenario is not advancing")
	}
}

func TestIncBattleWon(t *testing.T) {
	game := initNewGame()
	initBattleWon := game.BattleWon
	initGameScenario := game.Scenario

	game.IncBattleWon()
	require.Equal(t, initBattleWon+1, game.BattleWon, "mismatch number of battle won")
	require.NotEqual(t, initGameScenario, game.Scenario, "scenario is not advancing")
}

func initNewGame() *entity.Game {
	game, _ := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    testutil.NewTestMonster(),
		CreatedAt:  time.Now().Unix(),
	})
	return game
}
