package play_test

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/play"
	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	// define mock dependencies
	gameStorage := newMockGameStorage()
	partnerStorage := newMockPartnerStorage(nil)

	// define test cases
	testCases := []struct {
		Name    string
		Config  play.ServiceConfig
		IsError bool
	}{
		{
			Name: "Test Missing Game Storage",
			Config: play.ServiceConfig{
				GameStorage:    nil,
				PartnerStorage: partnerStorage,
			},
			IsError: true,
		},
		{
			Name: "Test Missing Partner Storage",
			Config: play.ServiceConfig{
				GameStorage:    gameStorage,
				PartnerStorage: nil,
			},
			IsError: true,
		},
		{
			Name: "Test Valid Config",
			Config: play.ServiceConfig{
				GameStorage:    gameStorage,
				PartnerStorage: partnerStorage,
			},
			IsError: false,
		},
	}
	// execute test cases
	for _, testcase := range testCases {
		t.Run(testcase.Name, func(t *testing.T) {
			_, err := play.NewService(testcase.Config)
			require.Equal(t, testcase.IsError, (err != nil), "unexpected error")
		})
	}
}

func TestServiceGetAvailablePartners(t *testing.T) {
	// initialize new service
	output := newService()
	// get available partners
	retPartners, err := output.Service.GetAvailablePartners(context.Background())
	require.NoError(t, err, "unexpected error")
	// check returned partners
	require.Equal(t, output.Partners, retPartners, "mismatch partners")
}

func TestServiceNewGame(t *testing.T) {
	// initialize new service
	output := newService()
	// create new game
	partner := output.Partners[0]
	game, err := output.Service.NewGame(context.Background(), "Riandy R.N", partner.ID)
	require.NoError(t, err, "unexpected error")
	// validate returned game with stored game, this is to make sure the game
	// is also stored on storage
	storedGame, err := output.GameStorage.GetGame(context.Background(), game.ID)
	require.NoError(t, err, "unexpected error")
	require.Equal(t, *game, *storedGame, "mismatch game")
}

func TestServiceGetGame(t *testing.T) {
	// initialize new service
	output := newService()
	// create new game
	partner := output.Partners[0]
	game, err := output.Service.NewGame(context.Background(), "Riandy R.N", partner.ID)
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
			ExpErr: play.ErrGameNotFound,
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
			retGame, err := output.Service.GetGame(context.Background(), testCase.GameID)
			assert.Equal(t, testCase.ExpErr, err, "mismatch error")
			if retGame == nil {
				return
			}
			assert.Equal(t, game, retGame, "mismatch game")
		})
	}
}

func newService() *newServiceOutput {
	// generate partners
	partners := []entity.Monster{
		*(testutil.NewTestMonster()),
		*(testutil.NewTestMonster()),
		*(testutil.NewTestMonster()),
		*(testutil.NewTestMonster()),
	}
	// initialize service
	cfg := play.ServiceConfig{
		GameStorage:    newMockGameStorage(),
		PartnerStorage: newMockPartnerStorage(partners),
	}
	svc, _ := play.NewService(cfg)

	return &newServiceOutput{
		Service:        svc,
		GameStorage:    cfg.GameStorage,
		PartnerStorage: cfg.PartnerStorage,
		Partners:       partners,
	}
}

type newServiceOutput struct {
	Service        play.Service
	GameStorage    play.GameStorage
	PartnerStorage play.PartnerStorage
	Partners       []entity.Monster
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
