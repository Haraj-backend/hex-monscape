package game

import (
	"context"

	"github.com/Haraj-backend/hex-pokebattle/internal/domain/playing"
)

type GameStorage struct{}

func (gs *GameStorage) GetGame(ctx context.Context, gameID string) (*playing.Game, error) {
	// TODO
	return nil, nil
}

func (gs *GameStorage) SaveGame(ctx context.Context, game playing.Game) error {
	// TODO
	return nil
}
