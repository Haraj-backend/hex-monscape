package gamestrg

import (
	"context"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type Storage struct {
	data map[string]entity.Game
}

func (s *Storage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	g := s.data[gameID]
	return &g, nil
}

func (s *Storage) SaveGame(ctx context.Context, game entity.Game) error {
	s.data[game.ID] = game
	return nil
}

func New() *Storage {
	return &Storage{data: make(map[string]entity.Game)}
}
