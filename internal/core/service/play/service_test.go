package play_test

/*
	The purpose of testing the Service component is to ensure it has correct
	implementation of business logic.

	The common pitfall when creating test for Service component is we tend to use
	concrete implementation for the dependency components (e.g actual GameStorage
	for MySQL). Not only this will increase the test complexity but also it will
	increase the possibility of getting false test result. The reason is simply
	because service such as MySQL has its own constraints & has much higher chance
	of failing rather than its mock counterpart (e.g network failure).

	So to avoid this pitfall, our first go to choice is to use mock implementation
	for the dependency when testing the Service component. This way we can control
	more the behavior of the dependency components to fit our test scenarios.
*/

import (
	"context"
	"errors"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/play"
	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/google/uuid"
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
	require.ElementsMatch(t, output.Partners, retPartners, "mismatch partners")

	// set error on get available partners
	output.PartnerStorage.SetRetErrOnGetAvailablePartners(true)

	// get available partners, should return error
	_, err = output.Service.GetAvailablePartners(context.Background())
	require.Error(t, err, "expected error")
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

	// create new game with invalid partner, should return error
	game, err = output.Service.NewGame(context.Background(), "Riandy R.N", uuid.NewString())
	require.Equal(t, play.ErrPartnerNotFound, err, "mismatch error")
	require.Nil(t, game, "unexpected game")

	// create new game with empty player name, should return error
	game, err = output.Service.NewGame(context.Background(), "", partner.ID)
	require.Error(t, err, "expected error")
	require.Nil(t, game, "unexpected game")

	// set error on save game, should return error
	output.GameStorage.SetRetErrOnSaveGame(true)
	game, err = output.Service.NewGame(context.Background(), "Riandy R.N", partner.ID)
	output.GameStorage.SetRetErrOnSaveGame(false)
	require.Error(t, err, "expected error")
	require.Nil(t, game, "unexpected game")

	// set error on get partner, should return error
	output.PartnerStorage.SetRetErrOnGetPartner(true)
	game, err = output.Service.NewGame(context.Background(), "Riandy R.N", partner.ID)
	output.PartnerStorage.SetRetErrOnGetPartner(false)
	require.Error(t, err, "expected error")
	require.Nil(t, game, "unexpected game")
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
			require.Equal(t, testCase.ExpErr, err, "mismatch error")
			if retGame == nil {
				return
			}
			require.Equal(t, game, retGame, "mismatch game")
		})
	}

	// set error on get game, should return error
	output.GameStorage.SetRetErrOnGetGame(true)
	game, err = output.Service.GetGame(context.Background(), game.ID)
	require.Error(t, err, "expected error")
	require.Nil(t, game, "unexpected game")
}

func newService() *newServiceOutput {
	// generate partners
	partners := []entity.Monster{
		*(testutil.NewTestMonster()),
		*(testutil.NewTestMonster()),
		*(testutil.NewTestMonster()),
		*(testutil.NewTestMonster()),
	}

	// initialize dependencies
	gameStorage := newMockGameStorage()
	partnerStorage := newMockPartnerStorage(partners)

	// initialize service
	cfg := play.ServiceConfig{
		GameStorage:    gameStorage,
		PartnerStorage: partnerStorage,
	}
	svc, _ := play.NewService(cfg)

	return &newServiceOutput{
		Service:        svc,
		GameStorage:    gameStorage,
		PartnerStorage: partnerStorage,
		Partners:       partners,
	}
}

type newServiceOutput struct {
	Service        play.Service
	GameStorage    *mockGameStorage
	PartnerStorage *mockPartnerStorage
	Partners       []entity.Monster
}

type mockGameStorage struct {
	data             map[string]entity.Game
	retErrOnGetGame  bool
	retErrOnSaveGame bool
}

func (gs *mockGameStorage) SetRetErrOnGetGame(retErr bool) {
	gs.retErrOnGetGame = retErr
}

func (gs *mockGameStorage) SetRetErrOnSaveGame(retErr bool) {
	gs.retErrOnSaveGame = retErr
}

func (gs *mockGameStorage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	if gs.retErrOnGetGame {
		return nil, ErrIntentionalError
	}

	game, ok := gs.data[gameID]
	if !ok {
		return nil, nil
	}
	return &game, nil
}

func (gs *mockGameStorage) SaveGame(ctx context.Context, game entity.Game) error {
	if gs.retErrOnSaveGame {
		return ErrIntentionalError
	}

	gs.data[game.ID] = game
	return nil
}

func newMockGameStorage() *mockGameStorage {
	return &mockGameStorage{
		data: map[string]entity.Game{},
	}
}

type mockPartnerStorage struct {
	partnerMap                   map[string]entity.Monster
	retErrOnGetAvailablePartners bool
	retErrOnGetPartner           bool
}

func (gs *mockPartnerStorage) SetRetErrOnGetAvailablePartners(retErr bool) {
	gs.retErrOnGetAvailablePartners = retErr
}

func (gs *mockPartnerStorage) SetRetErrOnGetPartner(retErr bool) {
	gs.retErrOnGetPartner = retErr
}

func (gs *mockPartnerStorage) GetAvailablePartners(ctx context.Context) ([]entity.Monster, error) {
	if gs.retErrOnGetAvailablePartners {
		return nil, ErrIntentionalError
	}
	var partners []entity.Monster
	for _, v := range gs.partnerMap {
		partners = append(partners, v)
	}
	return partners, nil
}

func (gs *mockPartnerStorage) GetPartner(ctx context.Context, partnerID string) (*entity.Monster, error) {
	if gs.retErrOnGetPartner {
		return nil, ErrIntentionalError
	}
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

var ErrIntentionalError = errors.New("intentional error")
