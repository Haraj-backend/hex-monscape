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
	// define test cases
	testCases := []struct {
		Name                string
		State               State
		Partner             entity.Monster
		Enemy               entity.Monster
		IsError             bool
		ExpectedEnemyHealth int
	}{
		{
			Name:  "Validate State PARTNER_TURN",
			State: PARTNER_TURN,
			Partner: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("pokemon_%v", time.Now().Unix()),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			Enemy: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("pokemon_%v", time.Now().Unix()),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   50,
					Speed:     100,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			IsError:             false,
			ExpectedEnemyHealth: 50,
		},
		{
			Name:    "Validate State DECIDE_TURN",
			State:   DECIDE_TURN,
			Partner: *newSamplePokemon(),
			Enemy:   *newSamplePokemon(),
			IsError: true,
		},
		{
			Name:    "Validate State WIN",
			State:   WIN,
			Partner: *newSamplePokemon(),
			Enemy:   *newSamplePokemon(),
			IsError: true,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battle, _ := NewBattle(BattleConfig{
				GameID:  "b1c87c5c-2ac3-471d-9880-4812552ee15d",
				Partner: &testCase.Partner,
				Enemy:   &testCase.Enemy,
			})

			battle.State = testCase.State
			err := battle.PartnerAttack()
			assert.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if !testCase.IsError {
				assert.Equal(t, battle.Enemy.BattleStats.Health, testCase.ExpectedEnemyHealth, "enemy health is not valid")
			}
		})
	}
}

func TestPartnerSurrender(t *testing.T) {
	battle := initNewBattle()
	// define test cases
	testCases := []struct {
		Name    string
		State   State
		IsError bool
	}{
		{
			Name:    "Validate State PARTNER_TURN",
			State:   PARTNER_TURN,
			IsError: false,
		},
		{
			Name:    "Validate State DECIDE_TURN",
			State:   DECIDE_TURN,
			IsError: true,
		},
		{
			Name:    "Validate State WIN",
			State:   WIN,
			IsError: true,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battle.State = testCase.State
			err := battle.PartnerSurrender()
			assert.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if !testCase.IsError {
				assert.Equal(t, LOSE, battle.State)
			}
		})
	}
}

func TestEnemyAttack(t *testing.T) {
	// define test cases
	testCases := []struct {
		Name                  string
		State                 State
		Partner               entity.Monster
		Enemy                 entity.Monster
		IsError               bool
		ExpectedPartnerHealth int
	}{
		{
			Name:  "Validate State ENEMY_TURN",
			State: ENEMY_TURN,
			Partner: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("pokemon_%v", time.Now().Unix()),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   50,
					Speed:     100,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			Enemy: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("pokemon_%v", time.Now().Unix()),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			IsError:               false,
			ExpectedPartnerHealth: 50,
		},
		{
			Name:    "Validate State PARTNER_TURN",
			State:   PARTNER_TURN,
			Partner: *newSamplePokemon(),
			Enemy:   *newSamplePokemon(),
			IsError: true,
		},
		{
			Name:    "Validate State DECIDE_TURN",
			State:   DECIDE_TURN,
			Partner: *newSamplePokemon(),
			Enemy:   *newSamplePokemon(),
			IsError: true,
		},
		{
			Name:    "Validate State WIN",
			State:   WIN,
			Partner: *newSamplePokemon(),
			Enemy:   *newSamplePokemon(),
			IsError: true,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battle, _ := NewBattle(BattleConfig{
				GameID:  "b1c87c5c-2ac3-471d-9880-4812552ee15d",
				Partner: &testCase.Partner,
				Enemy:   &testCase.Enemy,
			})
			battle.State = testCase.State
			err := battle.EnemyAttack()
			assert.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if !testCase.IsError {
				assert.Equal(t, battle.Partner.BattleStats.Health, testCase.ExpectedPartnerHealth, "partner health is not valid")
			}
		})
	}
}

func TestIsEnded(t *testing.T) {
	battle := initNewBattle()
	// define test cases
	testCases := []struct {
		Name     string
		State    State
		Expected bool
	}{
		{
			Name:     "Validate State PARTNER_TURN",
			State:    PARTNER_TURN,
			Expected: false,
		},
		{
			Name:     "Validate State DECIDE_TURN",
			State:    DECIDE_TURN,
			Expected: false,
		},
		{
			Name:     "Validate State ENEMY_TURN",
			State:    ENEMY_TURN,
			Expected: false,
		},
		{
			Name:     "Validate State WIN",
			State:    WIN,
			Expected: true,
		},
		{
			Name:     "Validate State LOSE",
			State:    LOSE,
			Expected: true,
		},
	}

	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battle.State = testCase.State
			actual := battle.IsEnded()
			assert.Equal(t, testCase.Expected, actual, "unexpected dead")
		})
	}
}

func TestDecideTurn(t *testing.T) {
	// define test cases
	testCases := []struct {
		Name          string
		State         State
		Partner       entity.Monster
		Enemy         entity.Monster
		IsError       bool
		ExpectedState State
	}{
		{
			Name:    "Validate State PARTNER_TURN",
			State:   PARTNER_TURN,
			Partner: *newSamplePokemon(),
			Enemy:   *newSamplePokemon(),
			IsError: true,
		},
		{
			Name:    "Validate State WIN",
			State:   WIN,
			Partner: *newSamplePokemon(),
			Enemy:   *newSamplePokemon(),
			IsError: true,
		},
		{
			Name:  "Validate PARTNER_TURN",
			State: DECIDE_TURN,
			Partner: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("pokemon_%v", time.Now().Unix()),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			Enemy: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("pokemon_%v", time.Now().Unix()),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     0,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			IsError:       false,
			ExpectedState: PARTNER_TURN,
		},
		{
			Name:  "Validate ENEMY_TURN",
			State: DECIDE_TURN,
			Partner: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("pokemon_%v", time.Now().Unix()),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     0,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			Enemy: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("pokemon_%v", time.Now().Unix()),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			IsError:       false,
			ExpectedState: ENEMY_TURN,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battle, _ := NewBattle(BattleConfig{
				GameID:  "b1c87c5c-2ac3-471d-9880-4812552ee15d",
				Partner: &testCase.Partner,
				Enemy:   &testCase.Enemy,
			})
			battle.State = testCase.State
			state, err := battle.DecideTurn()
			assert.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if !testCase.IsError {
				assert.Equal(t, testCase.ExpectedState, state, "expected state is not valid")
			}
		})
	}
}

func initNewBattle() *Battle {
	game, _ := NewBattle(BattleConfig{
		GameID:  "b1c87c5c-2ac3-471d-9880-4812552ee15d",
		Partner: newSamplePokemon(),
		Enemy:   newSamplePokemon(),
	})
	return game
}

func newSamplePokemon() *entity.Monster {
	currentTs := time.Now().Unix()
	return &entity.Monster{
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
