package monstrg_test

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/monstrg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		Name       string
		Config     monstrg.Config
		ExpStorage bool
		ExpErr     bool
	}{
		{
			Name:       "Empty Monster Data in Config",
			Config:     monstrg.Config{},
			ExpStorage: false,
			ExpErr:     true,
		},
		{
			Name: "Invalid Monster Data in Config",
			Config: monstrg.Config{
				MonsterData: []byte(`invalid`),
			},
			ExpStorage: false,
			ExpErr:     true,
		},
		{
			Name: "Valid Config",
			Config: monstrg.Config{
				MonsterData: monsterData,
			},
			ExpStorage: true,
			ExpErr:     false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			storage, err := monstrg.New(testCase.Config)
			assert.Equal(t, testCase.ExpStorage, storage != nil, "unexpected storage")
			assert.Equal(t, testCase.ExpErr, err != nil, "unexpected error")
		})
	}
}

func TestGetAvailablePartners(t *testing.T) {
	// init storage
	strg := initStorage(t)

	// get available partners
	partners, err := strg.GetAvailablePartners(context.Background())
	require.NoError(t, err)

	// we got the expected partners from monsterData variable
	expPartners := []entity.Monster{
		{
			ID:   "b1c87c5c-2ac3-471d-9880-4812552ee15d",
			Name: "Waneye",
			BattleStats: entity.BattleStats{
				Health:    100,
				MaxHealth: 100,
				Attack:    25,
				Defense:   5,
				Speed:     15,
			},
			AvatarURL: "https://monster.com/waneye.png",
		},
	}

	// ensure partners is equal to expected partners
	require.ElementsMatch(t, partners, expPartners, "unexpected partners")
}

func TestGetPartner(t *testing.T) {
	// init storage
	strg := initStorage(t)

	// define test cases
	testCases := []struct {
		Name       string
		PartnerID  string
		ExpPartner *entity.Monster
	}{
		{
			Name:      "Partner Found",
			PartnerID: "b1c87c5c-2ac3-471d-9880-4812552ee15d",
			ExpPartner: &entity.Monster{
				ID:   "b1c87c5c-2ac3-471d-9880-4812552ee15d",
				Name: "Waneye",
				BattleStats: entity.BattleStats{
					Health:    100,
					MaxHealth: 100,
					Attack:    25,
					Defense:   5,
					Speed:     15,
				},
				AvatarURL: "https://monster.com/waneye.png",
			},
		},
		{
			Name:       "Partner Not Found",
			PartnerID:  "88a98dee-ce84-4afb-b5a8-7cc07535f73f",
			ExpPartner: nil,
		},
	}

	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			partner, err := strg.GetPartner(context.Background(), testCase.PartnerID)
			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpPartner, partner, "unexpected partner")
		})
	}
}

func TestGetPossibleEnemies(t *testing.T) {
	// init storage
	strg := initStorage(t)

	// get possible enemies
	possibleEnemies, err := strg.GetPossibleEnemies(context.Background())
	if err != nil {
		t.Fatalf("unable to get possible enemies from storage, due: %v", err)
	}

	// we got the expected enemies from monsterData variable
	expPossibleEnemies := []entity.Monster{
		{
			ID:   "88a98dee-ce84-4afb-b5a8-7cc07535f73f",
			Name: "Bluebub",
			BattleStats: entity.BattleStats{
				Health:    100,
				MaxHealth: 100,
				Attack:    20,
				Defense:   10,
				Speed:     15,
			},
			AvatarURL: "https://monster.com/bluebub.png",
		},
		{
			ID:   "b1c87c5c-2ac3-471d-9880-4812552ee15d",
			Name: "Waneye",
			BattleStats: entity.BattleStats{
				Health:    100,
				MaxHealth: 100,
				Attack:    25,
				Defense:   5,
				Speed:     15,
			},
			AvatarURL: "https://monster.com/waneye.png",
		},
	}

	// ensure possible enemies is equal to expected possible enemies
	require.ElementsMatch(t, expPossibleEnemies, possibleEnemies, "enemies is not equal")
}

var monsterData = []byte(`
	[
		{
			"id": "88a98dee-ce84-4afb-b5a8-7cc07535f73f",
			"name": "Bluebub",
			"battle_stats": {
				"health": 100,
				"max_health": 100,
				"attack": 20,
				"defense": 10,
				"speed": 15
			},
			"avatar_url": "https://monster.com/bluebub.png",
			"is_partnerable": false
		},
		{
			"id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
			"name": "Waneye",
			"battle_stats": {
				"health": 100,
				"max_health": 100,
				"attack": 25,
				"defense": 5,
				"speed": 15
			},
			"avatar_url": "https://monster.com/waneye.png",
			"is_partnerable": true
		}
	]
`)

func initStorage(t *testing.T) *monstrg.Storage {
	strg, err := monstrg.New(monstrg.Config{
		MonsterData: monsterData,
	})
	require.NoError(t, err)

	return strg
}
