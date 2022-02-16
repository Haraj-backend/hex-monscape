package playing

import (
	"context"

	"github.com/riandyrn/pokebattle/internal/domain/entity"
)

func NewService(storage Storage) (*Service, error) {
	// TODO
	return nil, nil
}

type Service struct {
	storage Storage
}

func (s *Service) GetAvailablePartners(ctx context.Context) ([]entity.Pokemon, error)
func (s *Service) NewGame(ctx context.Context, playerName string, partnerID string) (*Game, error)
func (s *Service) GetGame(ctx context.Context, gameID string) (*Game, error)
func (s *Service) GetNextScenario(ctx context.Context, gameID string) (*Game, error)
