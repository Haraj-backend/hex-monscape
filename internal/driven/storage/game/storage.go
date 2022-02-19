package game

import (
	"context"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/playing"
)

type Storage struct {
	data map[string]playing.Game
}

func (s *Storage) GetGame(ctx context.Context, gameID string) (*playing.Game, error) {
	g, ok := s.data[gameID]
	if !ok {
		return nil, nil
	}
	return &g, nil
}

func (s *Storage) SaveGame(ctx context.Context, game playing.Game) error {
	s.data[game.ID] = game
	return nil
}

func New() *Storage {
	return &Storage{data: make(map[string]playing.Game)}
}
