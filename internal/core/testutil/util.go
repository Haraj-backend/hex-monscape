package testutil

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/google/uuid"
)

func NewTestMonster() *entity.Monster {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	uid := uuid.NewString()
	return &entity.Monster{
		ID:   uuid.NewString(),
		Name: fmt.Sprintf("monster_%v", uid),
		BattleStats: entity.BattleStats{
			Health:    r.Intn(100) + 1,
			MaxHealth: r.Intn(100) + 1,
			Attack:    r.Intn(100) + 1,
			Defense:   r.Intn(100) + 1,
			Speed:     r.Intn(100) + 1,
		},
		AvatarURL: fmt.Sprintf("https://example.com/%v.png", uid),
	}
}
