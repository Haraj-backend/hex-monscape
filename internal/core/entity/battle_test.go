package entity_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewBattle(t *testing.T) {
	// define function for validating battle
	validateBattle := func(t *testing.T, battle entity.Battle, cfg entity.BattleConfig) {
		require.NotEmpty(t, battle.GameID, "GameID is empty")
		require.NotEmpty(t, battle.Partner, "Partner is empty")
		require.NotEmpty(t, battle.Enemy, "Enemy is empty")
	}
	// define test cases
	testCases := []struct {
		Name    string
		Config  entity.BattleConfig
		IsError bool
	}{
		{
			Name:    "Invalid Config",
			Config:  entity.BattleConfig{},
			IsError: true,
		},
		{
			Name: "Valid Config",
			Config: entity.BattleConfig{
				GameID:  "b1c87c5c-2ac3-471d-9880-4812552ee15d",
				Partner: testutil.NewTestMonster(),
				Enemy:   testutil.NewTestMonster(),
			},
			IsError: false,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battle, err := entity.NewBattle(testCase.Config)
			require.Equal(t, testCase.IsError, (err != nil), "unexpected error")
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
		State               entity.State
		Partner             entity.Monster
		Enemy               entity.Monster
		IsError             bool
		ExpectedEnemyHealth int
	}{
		{
			Name:  "Validate State PARTNER_TURN",
			State: entity.StatePartnerTurn,
			Partner: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("monster_%v", time.Now().Unix()),
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
				Name: fmt.Sprintf("monster_%v", time.Now().Unix()),
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
			State:   entity.StateDecideTurn,
			Partner: *testutil.NewTestMonster(),
			Enemy:   *testutil.NewTestMonster(),
			IsError: true,
		},
		{
			Name:    "Validate State WIN",
			State:   entity.StateWin,
			Partner: *testutil.NewTestMonster(),
			Enemy:   *testutil.NewTestMonster(),
			IsError: true,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battle, _ := entity.NewBattle(entity.BattleConfig{
				GameID:  "b1c87c5c-2ac3-471d-9880-4812552ee15d",
				Partner: &testCase.Partner,
				Enemy:   &testCase.Enemy,
			})

			battle.State = testCase.State
			err := battle.PartnerAttack()
			require.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if !testCase.IsError {
				require.Equal(t, battle.Enemy.BattleStats.Health, testCase.ExpectedEnemyHealth, "enemy health is not valid")
			}
		})
	}
}

func TestPartnerAttackWin(t *testing.T) {
	battle := initNewBattle()
	// set partner attack to very high number so that enemy will be dead
	battle.Partner.BattleStats.Attack = math.MaxInt
	// set battle state to partner attack
	battle.State = entity.StatePartnerTurn
	// execute partner attack, enemy should be dead
	err := battle.PartnerAttack()
	require.NoError(t, err)
	// validate battle state
	require.Equal(t, entity.StateWin, battle.State)
}

func TestPartnerSurrender(t *testing.T) {
	battle := initNewBattle()
	// define test cases
	testCases := []struct {
		Name    string
		State   entity.State
		IsError bool
	}{
		{
			Name:    "Validate State PARTNER_TURN",
			State:   entity.StatePartnerTurn,
			IsError: false,
		},
		{
			Name:    "Validate State DECIDE_TURN",
			State:   entity.StateDecideTurn,
			IsError: true,
		},
		{
			Name:    "Validate State WIN",
			State:   entity.StateWin,
			IsError: true,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battle.State = testCase.State
			err := battle.PartnerSurrender()
			require.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if !testCase.IsError {
				require.Equal(t, entity.StateLose, battle.State)
			}
		})
	}
}

func TestEnemyAttack(t *testing.T) {
	// define test cases
	testCases := []struct {
		Name                  string
		State                 entity.State
		Partner               entity.Monster
		Enemy                 entity.Monster
		IsError               bool
		ExpectedPartnerHealth int
	}{
		{
			Name:  "Validate State ENEMY_TURN",
			State: entity.StateEnemyTurn,
			Partner: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("monster_%v", time.Now().Unix()),
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
				Name: fmt.Sprintf("monster_%v", time.Now().Unix()),
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
			State:   entity.StatePartnerTurn,
			Partner: *testutil.NewTestMonster(),
			Enemy:   *testutil.NewTestMonster(),
			IsError: true,
		},
		{
			Name:    "Validate State DECIDE_TURN",
			State:   entity.StateDecideTurn,
			Partner: *testutil.NewTestMonster(),
			Enemy:   *testutil.NewTestMonster(),
			IsError: true,
		},
		{
			Name:    "Validate State WIN",
			State:   entity.StateWin,
			Partner: *testutil.NewTestMonster(),
			Enemy:   *testutil.NewTestMonster(),
			IsError: true,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battle, _ := entity.NewBattle(entity.BattleConfig{
				GameID:  "b1c87c5c-2ac3-471d-9880-4812552ee15d",
				Partner: &testCase.Partner,
				Enemy:   &testCase.Enemy,
			})
			battle.State = testCase.State
			err := battle.EnemyAttack()
			require.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if !testCase.IsError {
				require.Equal(t, battle.Partner.BattleStats.Health, testCase.ExpectedPartnerHealth, "partner health is not valid")
			}
		})
	}
}

func TestEnemyAttackWin(t *testing.T) {
	battle := initNewBattle()
	// set enemy attack to very high number so that partner will be dead
	battle.Enemy.BattleStats.Attack = math.MaxInt
	// set battle state to enemy attack
	battle.State = entity.StateEnemyTurn
	// execute enemy attack, partner should be dead
	err := battle.EnemyAttack()
	require.NoError(t, err)
	// validate battle state
	require.Equal(t, entity.StateLose, battle.State)
}

func TestIsEnded(t *testing.T) {
	battle := initNewBattle()
	// define test cases
	testCases := []struct {
		Name     string
		State    entity.State
		Expected bool
	}{
		{
			Name:     "Validate State PARTNER_TURN",
			State:    entity.StatePartnerTurn,
			Expected: false,
		},
		{
			Name:     "Validate State DECIDE_TURN",
			State:    entity.StateDecideTurn,
			Expected: false,
		},
		{
			Name:     "Validate State ENEMY_TURN",
			State:    entity.StateEnemyTurn,
			Expected: false,
		},
		{
			Name:     "Validate State WIN",
			State:    entity.StateWin,
			Expected: true,
		},
		{
			Name:     "Validate State LOSE",
			State:    entity.StateLose,
			Expected: true,
		},
	}

	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battle.State = testCase.State
			actual := battle.IsEnded()
			require.Equal(t, testCase.Expected, actual, "unexpected dead")
		})
	}
}

func TestDecideTurn(t *testing.T) {
	// define test cases
	testCases := []struct {
		Name          string
		State         entity.State
		Partner       entity.Monster
		Enemy         entity.Monster
		IsError       bool
		ExpectedState entity.State
	}{
		{
			Name:    "Validate State PARTNER_TURN",
			State:   entity.StatePartnerTurn,
			Partner: *testutil.NewTestMonster(),
			Enemy:   *testutil.NewTestMonster(),
			IsError: true,
		},
		{
			Name:    "Validate State WIN",
			State:   entity.StateWin,
			Partner: *testutil.NewTestMonster(),
			Enemy:   *testutil.NewTestMonster(),
			IsError: true,
		},
		{
			Name:  "Validate PARTNER_TURN",
			State: entity.StateDecideTurn,
			Partner: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("monster_%v", time.Now().Unix()),
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
				Name: fmt.Sprintf("monster_%v", time.Now().Unix()),
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
			ExpectedState: entity.StatePartnerTurn,
		},
		{
			Name:  "Validate ENEMY_TURN",
			State: entity.StateDecideTurn,
			Partner: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("monster_%v", time.Now().Unix()),
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
				Name: fmt.Sprintf("monster_%v", time.Now().Unix()),
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
			ExpectedState: entity.StateEnemyTurn,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battle, _ := entity.NewBattle(entity.BattleConfig{
				GameID:  "b1c87c5c-2ac3-471d-9880-4812552ee15d",
				Partner: &testCase.Partner,
				Enemy:   &testCase.Enemy,
			})
			battle.State = testCase.State
			state, err := battle.DecideTurn()
			require.Equal(t, testCase.IsError, (err != nil), "unexpected error")
			if !testCase.IsError {
				require.Equal(t, testCase.ExpectedState, state, "expected state is not valid")
			}
		})
	}
}

func initNewBattle() *entity.Battle {
	game, _ := entity.NewBattle(entity.BattleConfig{
		GameID:  "b1c87c5c-2ac3-471d-9880-4812552ee15d",
		Partner: testutil.NewTestMonster(),
		Enemy:   testutil.NewTestMonster(),
	})
	return game
}
