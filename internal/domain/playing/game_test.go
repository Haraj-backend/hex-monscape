package playing

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riandyrn/pokebattle/internal/domain/entity"
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
			if (err != nil) != testCase.IsError {
				t.Fail()
			}
		})
	}
}

func TestNewGame(t *testing.T) {
	// TODO
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
