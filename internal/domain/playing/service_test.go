package playing

import (
	"context"
	"testing"

	"github.com/riandyrn/pokebattle/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	// define mock dependencies
	gameStorage := newMockGameStorage()
	partnerStorage := newMockPartnerStorage(nil)
	// define function for validating new game instance
	validateGame := func(t *testing.T, svc Service, cfg ServiceConfig) {
		assert.Equal(t, cfg.GameStorage, svc.gameStorage)
		assert.Equal(t, cfg.PartnerStorage, svc.partnerStorage)
	}
	// define test cases
	testCases := []struct {
		Name    string
		Config  ServiceConfig
		IsError bool
	}{
		{
			Name: "Test Missing Game Storage",
			Config: ServiceConfig{
				GameStorage:    nil,
				PartnerStorage: partnerStorage,
			},
			IsError: true,
		},
		{
			Name: "Test Missing Partner Storage",
			Config: ServiceConfig{
				GameStorage:    gameStorage,
				PartnerStorage: nil,
			},
			IsError: true,
		},
		{
			Name: "Test Valid Config",
			Config: ServiceConfig{
				GameStorage:    gameStorage,
				PartnerStorage: partnerStorage,
			},
			IsError: false,
		},
	}
	// execute test cases
	for _, testcase := range testCases {
		t.Run(testcase.Name, func(t *testing.T) {
			svc, err := NewService(testcase.Config)
			assert.Equal(t, testcase.IsError, (err != nil), "unexpected error")
			if svc == nil {
				return
			}
			validateGame(t, *svc, testcase.Config)
		})
	}
}

func TestServiceGetAvailablePartners(t *testing.T) {}

func TestServiceNewGame(t *testing.T) {}

func TestServiceGetGame(t *testing.T) {}

func TestServiceAdvanceScenario(t *testing.T) {}

type mockGameStorage struct {
	data map[string]Game
}

func (gs *mockGameStorage) GetGame(ctx context.Context, gameID string) (*Game, error) {
	game, ok := gs.data[gameID]
	if !ok {
		return nil, nil
	}
	return &game, nil
}

func (gs *mockGameStorage) SaveGame(ctx context.Context, game Game) error {
	gs.data[game.ID] = game
	return nil
}

func newMockGameStorage() *mockGameStorage {
	return &mockGameStorage{
		data: map[string]Game{},
	}
}

type mockPartnerStorage struct {
	partnerMap map[string]entity.Pokemon
}

func (gs *mockPartnerStorage) GetAvailablePartners(ctx context.Context) ([]entity.Pokemon, error) {
	var partners []entity.Pokemon
	for _, v := range gs.partnerMap {
		partners = append(partners, v)
	}
	return partners, nil
}

func (gs *mockPartnerStorage) GetPartner(ctx context.Context, partnerID string) (*entity.Pokemon, error) {
	partner, ok := gs.partnerMap[partnerID]
	if !ok {
		return nil, nil
	}
	return &partner, nil
}

func newMockPartnerStorage(partners []entity.Pokemon) *mockPartnerStorage {
	data := map[string]entity.Pokemon{}
	for _, partner := range partners {
		data[partner.ID] = partner
	}
	return &mockPartnerStorage{partnerMap: data}
}
