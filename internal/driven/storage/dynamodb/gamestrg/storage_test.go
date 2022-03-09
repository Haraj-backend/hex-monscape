package gamestrg

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/shared"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveGame(t *testing.T) {
	storage, err := newStorage()
	require.NoError(t, err)

	game := entity.Game{
		ID:         uuid.NewString(),
		PlayerName: "Alfat",
		Partner: &entity.Pokemon{
			ID: "partner-id",
		},
		BattleWon: 2,
		Scenario:  entity.BATTLE_3,
	}

	err = storage.SaveGame(context.Background(), game)
	require.NoError(t, err)

	savedGame, err := getGame(storage.dynamoClient, game.ID)
	require.NoError(t, err)
	require.Equal(t, game, *savedGame)
}

func TestGetGame(t *testing.T) {
	storage, err := newStorage()
	require.NoError(t, err)

	game := entity.Game{
		ID:         uuid.NewString(),
		PlayerName: "Alfat",
		Partner: &entity.Pokemon{
			ID: "partner-id",
		},
		BattleWon: 2,
		Scenario:  entity.BATTLE_3,
	}

	err = storage.SaveGame(context.Background(), game)
	require.NoError(t, err)

	testCases := []struct {
		Name    string
		GameID  string
		ExpGame *entity.Game
	}{
		{
			Name:    "Test Game Not Found",
			GameID:  uuid.NewString(),
			ExpGame: nil,
		},
		{
			Name:    "Test Game Found",
			GameID:  game.ID,
			ExpGame: &game,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			game, err := storage.GetGame(context.Background(), testCase.GameID)
			assert.NoError(t, err)

			if testCase.ExpGame == nil {
				return
			}

			assert.Equal(t, *testCase.ExpGame, *game)
		})
	}
}

func newStorage() (*Storage, error) {
	s, err := New(Config{
		DynamoClient: shared.NewLocalTestDDBClient(),
		TableName:    os.Getenv(shared.TestConfig.EnvKeyGameTableName),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize game storage due: %w", err)
	}
	return s, nil
}

func getGame(dynamoClient *dynamodb.DynamoDB, gameID string) (*entity.Game, error) {
	output, err := dynamoClient.GetItem(&dynamodb.GetItemInput{
		Key:       (gameKey{ID: gameID}).toDDBKey(),
		TableName: aws.String(os.Getenv(shared.TestConfig.EnvKeyGameTableName)),
	})

	if err != nil {
		return nil, fmt.Errorf("unable to get item due: %w", err)
	}

	if len(output.Item) == 0 {
		return nil, fmt.Errorf("game is not found")
	}

	var g entity.Game
	err = dynamodbattribute.UnmarshalMap(output.Item, &g)
	if err != nil {
		return nil, fmt.Errorf("unable to parse item due: %w", err)
	}
	return &g, nil
}
