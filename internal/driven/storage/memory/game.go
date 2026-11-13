package memory

import (
	"context"
	"sync"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type GameStorage struct {
	mu   sync.RWMutex
	data map[string]entity.Game
}

func NewGameStorage() *GameStorage {
	return &GameStorage{data: make(map[string]entity.Game)}
}

func (s *GameStorage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	g, exist := s.data[gameID]
	if !exist {
		// if item is not found, returns nil as expected by game interface
		return nil, nil
	}

	return &g, nil
}

func (s *GameStorage) SaveGame(ctx context.Context, game entity.Game) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[game.ID] = game

	return nil
}
