package battlestrg

import (
	"context"

	"github.com/Haraj-backend/hex-monscape/internal/core/battle"
)

type Storage struct {
	data map[string]battle.Battle
}

func (s *Storage) GetBattle(ctx context.Context, gameID string) (*battle.Battle, error) {
	b, ok := s.data[gameID]
	if !ok {
		return nil, nil
	}
	return &b, nil
}

func (s *Storage) SaveBattle(ctx context.Context, b battle.Battle) error {
	s.data[b.GameID] = b
	return nil
}

func New() *Storage {
	return &Storage{data: make(map[string]battle.Battle)}
}
