package gamestrg

import (
	"context"
	"sync"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type Storage struct {
	mu   sync.RWMutex
	data map[string]entity.Game
}

func (s *Storage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	g, exist := s.data[gameID]
	if !exist {
		// if item is not found, returns nil as expected by game interface
		return nil, nil
	}

	return &g, nil
}

func (s *Storage) SaveGame(ctx context.Context, game entity.Game) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[game.ID] = game
	return nil
}

func New() *Storage {
	return &Storage{data: make(map[string]entity.Game)}
}
