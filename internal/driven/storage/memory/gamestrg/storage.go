package gamestrg

import (
	"context"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/shared/telemetry"
	"go.opentelemetry.io/otel/attribute"
)

type Storage struct {
	data map[string]entity.Game
}

func (s *Storage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "GameStorage: GetGame")
	defer span.End()

	span.SetAttributes(attribute.Key("game-id").String(gameID))

	g, ok := s.data[gameID]
	if !ok {
		return nil, nil
	}
	return &g, nil
}

func (s *Storage) SaveGame(ctx context.Context, game entity.Game) error {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "GameStorage: SaveGame")
	defer span.End()

	span.SetAttributes(attribute.Key("game-id").String(game.ID))

	s.data[game.ID] = game
	return nil
}

func New() *Storage {
	return &Storage{data: make(map[string]entity.Game)}
}
