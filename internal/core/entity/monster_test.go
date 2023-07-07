package entity_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIsDead(t *testing.T) {
	// define test cases
	testCases := []struct {
		Name     string
		Monster  entity.Monster
		Expected bool
	}{
		{
			Name: "Monster is Not Dead",
			Monster: entity.Monster{
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
			Expected: false,
		},
		{
			Name: "Monster Has 0 Health",
			Monster: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("monster_%v", time.Now().Unix()),
				BattleStats: entity.BattleStats{
					Health:    0,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			Expected: true,
		},
		{
			Name: "Monster Has Negative Health",
			Monster: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("monster_%v", time.Now().Unix()),
				BattleStats: entity.BattleStats{
					Health:    -100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			Expected: true,
		},
	}

	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := testCase.Monster.IsDead()
			assert.Equal(t, testCase.Expected, actual, "unexpected dead")
		})
	}
}

func TestInflictDamage(t *testing.T) {
	// define test cases
	testCases := []struct {
		Name                 string
		Monster              entity.Monster
		Enemy                entity.Monster
		ExpectedHealthAmount int
	}{
		{
			Name: "Monster Get Zero Damage",
			Monster: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("monster_%v", time.Now().Unix()),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   0,
					Speed:     100,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			Enemy: entity.Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("enemy_%v", time.Now().Unix()),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			ExpectedHealthAmount: 0,
		},
		{
			Name: "Monster Get 50 Damage",
			Monster: entity.Monster{
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
				Name: fmt.Sprintf("enemy_%v", time.Now().Unix()),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			ExpectedHealthAmount: 50,
		},
	}

	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := testCase.Monster.InflictDamage(testCase.Enemy)
			if err != nil {
				t.Errorf("unable to inflict damage, due: %v", err)
			}
			assert.Equal(t, testCase.ExpectedHealthAmount, testCase.Monster.BattleStats.Health, "unexpected health amount")
		})
	}
}
