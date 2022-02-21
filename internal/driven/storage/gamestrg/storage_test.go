package gamestrg

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetGame(t *testing.T) {
	strg := New()
	game := initNewGame()
	err := strg.SaveGame(context.Background(), *game)
	if err != nil {
		t.Fatalf("unable to save game, due: %v", err)
	}
	newGame, err := strg.GetGame(context.Background(), game.ID)
	if err != nil {
		t.Fatalf("unable to get game, due: %v", err)
	}
	require.Equal(t, newGame, game, "game is not equal")
}

func TestSaveGame(t *testing.T) {
	strg := New()
	game := initNewGame()
	err := strg.SaveGame(context.Background(), *game)
	if err != nil {
		t.Fatalf("unable to save game, due: %v", err)
	}
	newGame, err := strg.GetGame(context.Background(), game.ID)
	if err != nil {
		t.Fatalf("unable to get game, due: %v", err)
	}
	require.Equal(t, newGame, game, "game is not equal")
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

func initNewGame() *entity.Game {
	game, _ := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    newSamplePokemon(),
		CreatedAt:  time.Now().Unix(),
	})
	return game
}
