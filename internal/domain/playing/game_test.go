package playing

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riandyrn/pokebattle/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestValidateGameConfig(t *testing.T) {
	samplePokemon := newSamplePokemon()
	currentTs := time.Now().Unix()
	testCases := []struct {
		Name    string
		Config  GameConfig
		IsError bool
	}{
		{
			Name: "Empty Player Name",
			Config: GameConfig{
				PlayerName: "",
				Partner:    samplePokemon,
				CreatedAt:  currentTs,
			},
			IsError: true,
		},
		{
			Name: "Empty Partner",
			Config: GameConfig{
				PlayerName: "Riandy R.N",
				Partner:    nil,
				CreatedAt:  currentTs,
			},
			IsError: true,
		},
		{
			Name: "Empty Created At",
			Config: GameConfig{
				PlayerName: "Riandy R.N",
				Partner:    samplePokemon,
				CreatedAt:  0,
			},
			IsError: true,
		},
		{
			Name: "All Filled",
			Config: GameConfig{
				PlayerName: "Riandy R.N",
				Partner:    samplePokemon,
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
	validateGame := func(t *testing.T, game Game, cfg GameConfig) {
		assert.NotEmpty(t, game.ID, "game id is empty")
		assert.Equal(t, cfg.PlayerName, game.PlayerName, "mismatch player name")
		assert.Equal(t, cfg.Partner, game.Partner, "mismatch partner")
		assert.Equal(t, cfg.CreatedAt, game.CreatedAt, "mismatch created_at")
		assert.Equal(t, 0, game.BattleWon, "mismatch battle_won")
		assert.Equal(t, BATTLE_1, game.Scenario, "mismatch scenario")
	}
	// define test cases
	testCases := []struct {
		Name    string
		Config  GameConfig
		IsError bool
	}{
		{
			Name:    "Invalid Config",
			Config:  GameConfig{},
			IsError: true,
		},
		{
			Name: "Valid Config",
			Config: GameConfig{
				PlayerName: "Riandy R.N",
				Partner:    newSamplePokemon(),
				CreatedAt:  time.Now().Unix(),
			},
			IsError: false,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			game, err := NewGame(testCase.Config)
			assert.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if game == nil {
				return
			}
			validateGame(t, *game, testCase.Config)
		})
	}
}

func TestAdvanceScenario(t *testing.T) {
	// TODO
}

func TestIncBattleWon(t *testing.T) {
	// TODO
}

func newSamplePokemon() *entity.Pokemon {
	currentTs := time.Now().Unix()
	return &entity.Pokemon{
		ID:   uuid.NewString(),
		Name: fmt.Sprintf("pokemon_%v", currentTs),
		BattleStats: entity.BattleStats{
			Health:    100,
			MaxHealth: 100,
			Attack:    100,
			Defense:   100,
			Speed:     100,
		},
		AvatarURL: fmt.Sprintf("https://example.com/%v", currentTs),
	}
}
