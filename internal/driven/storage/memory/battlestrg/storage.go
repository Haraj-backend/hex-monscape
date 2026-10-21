package battlestrg

import (
	"context"
	"sync"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type Storage struct {
	mu   sync.RWMutex
	data map[string]entity.Battle
}

func (s *Storage) GetBattle(ctx context.Context, gameID string) (*entity.Battle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, exist := s.data[gameID]
	if !exist {
		// if item is not found, returns nil as expected by battle interface
		return nil, nil
	}

	return &b, nil
}

func (s *Storage) SaveBattle(ctx context.Context, b entity.Battle) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[b.GameID] = b
	return nil
}

func New() *Storage {
	return &Storage{data: make(map[string]entity.Battle)}
}
