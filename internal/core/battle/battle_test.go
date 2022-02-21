package battle

import (
	"fmt"
	"testing"
	"time"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBattle(t *testing.T) {
	// define function for validating battle
	validateBattle := func(t *testing.T, battle Battle, cfg BattleConfig) {
		assert.NotEmpty(t, battle.GameID, "GameID is empty")
		assert.NotEmpty(t, battle.Partner, "Partner is empty")
		assert.NotEmpty(t, battle.Enemy, "Enemy is empty")
	}
	// define test cases
	testCases := []struct {
		Name    string
		Config  BattleConfig
		IsError bool
	}{
		{
			Name:    "Invalid Config",
			Config:  BattleConfig{},
			IsError: true,
		},
		{
			Name: "Valid Config",
			Config: BattleConfig{
				GameID:  "b1c87c5c-2ac3-471d-9880-4812552ee15d",
				Partner: newSamplePokemon(),
				Enemy:   newSamplePokemon(),
			},
			IsError: false,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battle, err := NewBattle(testCase.Config)
			assert.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if battle == nil {
				return
			}
			validateBattle(t, *battle, testCase.Config)
		})
	}
}

func TestPartnerAttack(t *testing.T) {
	// TODO
	battle := initNewBattle()
	battle.PartnerAttack()
}

func TestPartnerSurrender(t *testing.T) {
	// TODO
}

func TestEnemyAttack(t *testing.T) {
	// TODO
}

func TestIsEnded(t *testing.T) {
	// TODO
}

func TestDecideTurn(t *testing.T) {
	// TODO
}

func initNewBattle() *Battle {
	game, _ := NewBattle(BattleConfig{
		GameID:  "b1c87c5c-2ac3-471d-9880-4812552ee15d",
		Partner: newSamplePokemon(),
		Enemy:   newSamplePokemon(),
	})
	return game
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
