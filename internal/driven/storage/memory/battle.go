package memory

import (
	"context"
	"sync"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type BattleStorage struct {
	mu   sync.RWMutex
	data map[string]entity.Battle
}

func NewBattleStorage() *BattleStorage {
	return &BattleStorage{data: make(map[string]entity.Battle)}
}

func (s *BattleStorage) GetBattle(ctx context.Context, gameID string) (*entity.Battle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, exist := s.data[gameID]
	if !exist {
		// if item is not found, returns nil as expected by battle interface
		return nil, nil
	}

	return &b, nil
}

func (s *BattleStorage) SaveBattle(ctx context.Context, b entity.Battle) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[b.GameID] = b

	return nil
}
