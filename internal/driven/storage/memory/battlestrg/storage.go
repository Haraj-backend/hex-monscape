package battlestrg

import (
	"context"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type Storage struct {
	data map[string]entity.Battle
}

func (s *Storage) GetBattle(ctx context.Context, gameID string) (*entity.Battle, error) {
	b, ok := s.data[gameID]
	if !ok {
		return nil, nil
	}
	return &b, nil
}

func (s *Storage) SaveBattle(ctx context.Context, b entity.Battle) error {
	s.data[b.GameID] = b
	return nil
}

func New() *Storage {
	return &Storage{data: make(map[string]entity.Battle)}
}
