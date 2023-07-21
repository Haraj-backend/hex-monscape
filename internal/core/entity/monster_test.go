package entity_test

import (
	"fmt"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
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
				Name: generateMonsterName(false),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: generateAvatarURL(),
			},
			Expected: false,
		},
		{
			Name: "Monster Has 0 Health",
			Monster: entity.Monster{
				ID:   uuid.NewString(),
				Name: generateMonsterName(false),
				BattleStats: entity.BattleStats{
					Health:    0,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: generateAvatarURL(),
			},
			Expected: true,
		},
		{
			Name: "Monster Has Negative Health",
			Monster: entity.Monster{
				ID:   uuid.NewString(),
				Name: generateMonsterName(false),
				BattleStats: entity.BattleStats{
					Health:    -100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: generateAvatarURL(),
			},
			Expected: true,
		},
	}

	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := testCase.Monster.IsDead()
			require.Equal(t, testCase.Expected, actual, "unexpected dead")
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
				Name: generateMonsterName(false),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   0,
					Speed:     100,
				},
				AvatarURL: generateAvatarURL(),
			},
			Enemy: entity.Monster{
				ID:   uuid.NewString(),
				Name: generateMonsterName(true),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: generateAvatarURL(),
			},
			ExpectedHealthAmount: 0,
		},
		{
			Name: "Monster Get 50 Damage",
			Monster: entity.Monster{
				ID:   uuid.NewString(),
				Name: generateMonsterName(false),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   50,
					Speed:     100,
				},
				AvatarURL: generateAvatarURL(),
			},
			Enemy: entity.Monster{
				ID:   uuid.NewString(),
				Name: generateMonsterName(true),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: generateAvatarURL(),
			},
			ExpectedHealthAmount: 50,
		},
		{
			Name: "Enemy Attack is Lower Than Monster Defense",
			Monster: entity.Monster{
				ID:   uuid.NewString(),
				Name: generateMonsterName(false),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    100,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: generateAvatarURL(),
			},
			Enemy: entity.Monster{
				ID:   uuid.NewString(),
				Name: generateMonsterName(true),
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    10,
					Defense:   100,
					Speed:     100,
				},
				AvatarURL: generateAvatarURL(),
			},
			ExpectedHealthAmount: 100 - entity.MinDamage,
		},
	}

	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := testCase.Monster.InflictDamage(testCase.Enemy)
			if err != nil {
				t.Errorf("unable to inflict damage, due: %v", err)
			}
			require.Equal(t, testCase.ExpectedHealthAmount, testCase.Monster.BattleStats.Health, "unexpected health amount")
		})
	}
}

func TestResetBattleStats(t *testing.T) {
	// create new monster
	m := testutil.NewTestMonster()

	// set the monster health to 0
	m.BattleStats.Health = 0

	// reset the battle stats
	m.ResetBattleStats()

	// the monster health should be equal to max health
	require.Equal(t, m.BattleStats.MaxHealth, m.BattleStats.Health, "unexpected health amount")
}

func generateMonsterName(isEnemy bool) string {
	prefix := "monster"
	if isEnemy {
		prefix = "enemy"
	}
	return fmt.Sprintf("%v_%v", prefix, uuid.NewString())
}

func generateAvatarURL() string {
	return fmt.Sprintf("https://example.com/%v.jpg", uuid.NewString())
}
