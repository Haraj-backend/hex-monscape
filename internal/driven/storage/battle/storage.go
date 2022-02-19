package battle

import (
	"context"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battling"
)

type Storage struct {
	data map[string]battling.Battle
}

func (s *Storage) GetBattle(ctx context.Context, gameID string) (*battling.Battle, error) {
	b, ok := s.data[gameID]
	if !ok {
		return nil, nil
	}
	return &b, nil
}

func (s *Storage) SaveBattle(ctx context.Context, b battling.Battle) error {
	s.data[b.GameID] = b
	return nil
}

func New() *Storage {
	return &Storage{data: make(map[string]battling.Battle)}
}
