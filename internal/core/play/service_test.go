package play

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	// define mock dependencies
	gameStorage := newMockGameStorage()
	partnerStorage := newMockPartnerStorage(nil)
	// define function for validating new game instance
	validateGame := func(t *testing.T, svc Service, cfg ServiceConfig) {
		assert.Equal(t, cfg.GameStorage, svc.(*service).gameStorage)
		assert.Equal(t, cfg.PartnerStorage, svc.(*service).partnerStorage)
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
			validateGame(t, svc, testcase.Config)
		})
	}
}

func TestServiceGetAvailablePartners(t *testing.T) {
	// initialize new service
	svc, partners := newNewService()
	// get available partners
	retPartners, err := svc.GetAvailablePartners(context.Background())
	require.NoError(t, err, "unexpected error")
	// check returned partners
	require.Equal(t, partners, retPartners, "mismatch partners")
}

func TestServiceNewGame(t *testing.T) {
	// initialize new service
	svc, partners := newNewService()
	// create new game
	partner := partners[0]
	game, err := svc.NewGame(context.Background(), "Riandy R.N", partner.ID)
	require.NoError(t, err, "unexpected error")
	// validate returned game with stored game, this is to make sure the game
	// is also stored on storage
	storedGame, err := svc.(*service).gameStorage.GetGame(context.Background(), game.ID)
	require.NoError(t, err, "unexpected error")
	require.Equal(t, *game, *storedGame, "mismatch game")
}

func TestServiceGetGame(t *testing.T) {
	// initialize new service
	svc, partners := newNewService()
	// create new game
	partner := partners[0]
	game, err := svc.NewGame(context.Background(), "Riandy R.N", partner.ID)
	require.NoError(t, err, "unexpected error")
	// define test cases
	testCases := []struct {
		Name   string
		GameID string
		ExpErr error
	}{
		{
			Name:   "Test Game Not Found",
			GameID: game.ID + "abc",
			ExpErr: ErrGameNotFound,
		},
		{
			Name:   "Test Game Found",
			GameID: game.ID,
			ExpErr: nil,
		},
	}
	// execute test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			retGame, err := svc.GetGame(context.Background(), testCase.GameID)
			assert.Equal(t, testCase.ExpErr, err, "mismatch error")
			if retGame == nil {
				return
			}
			assert.Equal(t, game, retGame, "mismatch game")
		})
	}
}

type mockGameStorage struct {
	data map[string]entity.Game
}

func (gs *mockGameStorage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	game, ok := gs.data[gameID]
	if !ok {
		return nil, nil
	}
	return &game, nil
}

func (gs *mockGameStorage) SaveGame(ctx context.Context, game entity.Game) error {
	gs.data[game.ID] = game
	return nil
}

func newMockGameStorage() *mockGameStorage {
	return &mockGameStorage{
		data: map[string]entity.Game{},
	}
}

type mockPartnerStorage struct {
	partnerMap map[string]entity.Monster
}

func (gs *mockPartnerStorage) GetAvailablePartners(ctx context.Context) ([]entity.Monster, error) {
	var partners []entity.Monster
	for _, v := range gs.partnerMap {
		partners = append(partners, v)
	}
	return partners, nil
}

func (gs *mockPartnerStorage) GetPartner(ctx context.Context, partnerID string) (*entity.Monster, error) {
	partner, ok := gs.partnerMap[partnerID]
	if !ok {
		return nil, nil
	}
	return &partner, nil
}

func newMockPartnerStorage(partners []entity.Monster) *mockPartnerStorage {
	data := map[string]entity.Monster{}
	for _, partner := range partners {
		data[partner.ID] = partner
	}
	return &mockPartnerStorage{partnerMap: data}
}

func newNewService() (Service, []entity.Monster) {
	// generate partners
	partners := []entity.Monster{
		{
			ID:   "b1c87c5c-2ac3-471d-9880-4812552ee15d",
			Name: "Pikachu",
			BattleStats: entity.BattleStats{
				Health:    100,
				MaxHealth: 100,
				Attack:    25,
				Defense:   5,
				Speed:     10,
			},
			AvatarURL: "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png",
		},
	}
	// initialize service
	cfg := ServiceConfig{
		GameStorage:    newMockGameStorage(),
		PartnerStorage: newMockPartnerStorage(partners),
	}
	svc, _ := NewService(cfg)

	return svc, partners
}
