package battle

import (
	"context"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type BattleStorage interface {
	// GetBattle returns battle instance for given gameID from storage. Returns nil
	// when battle instance is not found.
	GetBattle(ctx context.Context, gameID string) (*entity.Battle, error)

	// SaveBattle is used for saving given battle instance into storage. If battle
	// instance is already exists in the storage, it will be overwritten.
	SaveBattle(ctx context.Context, b entity.Battle) error
}

type GameStorage interface {
	// GetGame returns game instance for given gameID from storage. Returns nil
	// when given gameID is not found in database.
	GetGame(ctx context.Context, gameID string) (*entity.Game, error)

	// Save is used for saving game instance in storage.
	SaveGame(ctx context.Context, game entity.Game) error
}

type MonsterStorage interface {
	// GetPossibleEnemies returns all possible enemies in the game. Returns nil
	// when there is no possible enemies.
	GetPossibleEnemies(ctx context.Context) ([]entity.Monster, error)
}
