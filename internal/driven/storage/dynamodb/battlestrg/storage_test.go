package battlestrg

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/shared"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveBattle(t *testing.T) {
	// initialize storage
	storage, err := newStorage()
	require.NoError(t, err)

	// save battle
	battle := battle.Battle{
		GameID: uuid.NewString(),
		State:  battle.DECIDE_TURN,
		Partner: &entity.Pokemon{
			ID:   "partner-id",
			Name: "my-partner",
		},
		Enemy: &entity.Pokemon{
			ID:   "enemy-id",
			Name: "my-enemy",
		},
		LastDamage: battle.LastDamage{
			Partner: 10,
			Enemy:   10,
		},
	}
	err = storage.SaveBattle(context.Background(), battle)
	require.NoError(t, err)

	// check if battle is truly saved, notice that here we are using independently
	// defined `getBattle()` instead `storage.GetBattle()` the reason is because
	// the implementation of `storage.GetBattle()` is not yet tested so we cannot
	// sure if it is implemented correctly, thus we use our own defined `getBattle()`
	savedBattle, err := getBattle(storage.dynamoClient, battle.GameID)
	require.NoError(t, err)
	require.Equal(t, battle, *savedBattle)
}

func TestGetBattle(t *testing.T) {
	// initialize storage
	storage, err := newStorage()
	require.NoError(t, err)

	// create new battle, notice that we use `storage.SaveBattle()` here since
	// it is already tested, so we can sure it is already correctly implemented
	bt := battle.Battle{
		GameID: uuid.NewString(),
		State:  battle.DECIDE_TURN,
		Partner: &entity.Pokemon{
			ID:   "partner-id",
			Name: "my-partner",
		},
		Enemy: &entity.Pokemon{
			ID:   "enemy-id",
			Name: "my-enemy",
		},
		LastDamage: battle.LastDamage{
			Partner: 10,
			Enemy:   10,
		},
	}
	err = storage.SaveBattle(context.Background(), bt)
	require.NoError(t, err)

	// define testcases, here we want to test condition when battle
	// is found & not found
	testCases := []struct {
		Name      string
		GameID    string
		ExpBattle *battle.Battle
	}{
		{
			Name:      "Test Battle Not Found",
			GameID:    uuid.NewString(),
			ExpBattle: nil,
		},
		{
			Name:      "Test Battle Found",
			GameID:    bt.GameID,
			ExpBattle: &bt,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// get battle from storage
			bt, err := storage.GetBattle(context.Background(), testCase.GameID)
			assert.NoError(t, err)
			// if expected battle is nil no further check required
			if testCase.ExpBattle == nil {
				return
			}
			// check value of battle from storage
			assert.Equal(t, *testCase.ExpBattle, *bt)
		})
	}
}

func newStorage() (*Storage, error) {
	awsSess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String(os.Getenv(shared.TestConfig.EnvKeyLocalstackEndpoint)),
		},
	}))
	s, err := New(Config{
		DynamoClient: dynamodb.New(awsSess),
		TableName:    os.Getenv(shared.TestConfig.EnvKeyBattleTableName),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize battle storage due: %w", err)
	}
	return s, nil
}

func getBattle(dynamoClient *dynamodb.DynamoDB, gameID string) (*battle.Battle, error) {
	output, err := dynamoClient.GetItem(&dynamodb.GetItemInput{
		Key:       (battleKey{GameID: gameID}).toDDBKey(),
		TableName: aws.String(os.Getenv(shared.TestConfig.EnvKeyBattleTableName)),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get item due: %w", err)
	}
	if len(output.Item) == 0 {
		// notice that this behavior is different from `battlestrg.Storage.GetBattle()`
		// implementation, in here when item is not found we return error
		return nil, fmt.Errorf("battle is not found")
	}
	var b battle.Battle
	err = dynamodbattribute.UnmarshalMap(output.Item, &b)
	if err != nil {
		return nil, fmt.Errorf("unable to parse item due: %w", err)
	}
	return &b, nil
}
