package battle

import (
	"context"

	"github.com/Haraj-backend/hex-pokebattle/internal/domain/battling"
)

type BattleStorage struct{}

func (bs *BattleStorage) GetBattle(ctx context.Context, gameID string) (*battling.Battle, error) {
	// TODO
	return nil, nil
}

func (bs *BattleStorage) SaveBattle(ctx context.Context, b battling.Battle) error {
	// TODO
	return nil
}
