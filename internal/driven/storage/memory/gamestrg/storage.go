package gamestrg

import (
	"context"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/shared/telemetry"
)

type Storage struct {
	data map[string]entity.Game
}

func (s *Storage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	tr := telemetry.GetTracer()
	span := tr.Trace(ctx, "GetGame GameStorage")
	defer span.End()

	g, ok := s.data[gameID]
	if !ok {
		return nil, nil
	}
	return &g, nil
}

func (s *Storage) SaveGame(ctx context.Context, game entity.Game) error {
	tr := telemetry.GetTracer()
	span := tr.Trace(ctx, "SaveGame GameStorage")
	defer span.End()

	s.data[game.ID] = game
	return nil
}

func New() *Storage {
	return &Storage{data: make(map[string]entity.Game)}
}
