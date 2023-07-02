package monstrg

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetAvailablePartners(t *testing.T) {
	partners := []entity.Monster{
		*newSamplePokemon(),
	}
	strg, err := New(Config{
		Partners: partners,
		Enemies: []entity.Monster{
			*newSamplePokemon(),
		},
	})
	if err != nil {
		t.Fatalf("unable to init new storage, due: %v", err)
	}
	newPartners, err := strg.GetAvailablePartners(context.Background())
	if err != nil {
		t.Fatalf("unable to get available partners from storage, due: %v", err)
	}
	require.Equal(t, partners, newPartners, "partners is not equal")
}

func TestGetPartner(t *testing.T) {
	partner := newSamplePokemon()
	strg, err := New(Config{
		Partners: []entity.Monster{*partner},
		Enemies: []entity.Monster{
			*newSamplePokemon(),
		},
	})
	if err != nil {
		t.Fatalf("unable to init new storage, due: %v", err)
	}
	newPartner, err := strg.GetPartner(context.Background(), partner.ID)
	if err != nil {
		t.Fatalf("unable to get partner from storage, due: %v", err)
	}
	require.Equal(t, partner, newPartner, "partner is not equal")
}

func TestGetPossibleEnemies(t *testing.T) {
	enemies := []entity.Monster{
		*newSamplePokemon(),
	}
	strg, err := New(Config{
		Partners: []entity.Monster{
			*newSamplePokemon(),
		},
		Enemies: enemies,
	})
	if err != nil {
		t.Fatalf("unable to init new storage, due: %v", err)
	}
	newEnemies, err := strg.GetPossibleEnemies(context.Background())
	if err != nil {
		t.Fatalf("unable to get possible enemies from storage, due: %v", err)
	}
	require.Equal(t, enemies, newEnemies, "enemies is not equal")
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
