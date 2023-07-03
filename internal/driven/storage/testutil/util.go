package testutil

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/google/uuid"
)

func NewTestMonster() *entity.Monster {
	now := time.Now()
	r := rand.New(rand.NewSource(now.UnixNano()))

	return &entity.Monster{
		ID:   uuid.NewString(),
		Name: fmt.Sprintf("monster_%v", now.Unix()),
		BattleStats: entity.BattleStats{
			Health:    r.Intn(100) + 1,
			MaxHealth: r.Intn(100) + 1,
			Attack:    r.Intn(100) + 1,
			Defense:   r.Intn(100) + 1,
			Speed:     r.Intn(100) + 1,
		},
		AvatarURL: fmt.Sprintf("https://example.com/%v", now.Unix()),
	}
}
