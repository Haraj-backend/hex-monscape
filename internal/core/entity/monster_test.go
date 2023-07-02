package entity

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIsDead(t *testing.T) {
	// define test cases
	testCases := []struct {
		Name     string
		Pokemon  Monster
		Expected bool
	}{
		{
			Name: "Pokemon is Not Dead",
			Pokemon: Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("pokemon_%v", time.Now().Unix()),
				BattleStats: BattleStats{
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
			Name: "Pokemon Has 0 Health",
			Pokemon: Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("pokemon_%v", time.Now().Unix()),
				BattleStats: BattleStats{
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
			Name: "Pokemon Has Negative Health",
			Pokemon: Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("pokemon_%v", time.Now().Unix()),
				BattleStats: BattleStats{
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
			actual := testCase.Pokemon.IsDead()
			assert.Equal(t, testCase.Expected, actual, "unexpected dead")
		})
	}
}

func TestInflictDamage(t *testing.T) {
	// define test cases
	testCases := []struct {
		Name                 string
		Pokemon              Monster
		Enemy                Monster
		ExpectedHealthAmount int
	}{
		{
			Name: "Pokemon Get Zero Damage",
			Pokemon: Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("pokemon_%v", time.Now().Unix()),
				BattleStats: BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   0,
					Speed:     100,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			Enemy: Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("enemy_%v", time.Now().Unix()),
				BattleStats: BattleStats{
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
			Name: "Pokemon Get 50 Damage",
			Pokemon: Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("pokemon_%v", time.Now().Unix()),
				BattleStats: BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   50,
					Speed:     100,
				},
				AvatarURL: fmt.Sprintf("https://example.com/%v", time.Now().Unix()),
			},
			Enemy: Monster{
				ID:   uuid.NewString(),
				Name: fmt.Sprintf("enemy_%v", time.Now().Unix()),
				BattleStats: BattleStats{
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
			_, err := testCase.Pokemon.InflictDamage(testCase.Enemy)
			if err != nil {
				t.Errorf("unable to inflict damage, due: %v", err)
			}
			assert.Equal(t, testCase.ExpectedHealthAmount, testCase.Pokemon.BattleStats.Health, "unexpected health amount")
		})
	}
}
