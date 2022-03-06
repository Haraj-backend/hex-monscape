package gamestrg_test

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/gamestrg"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockedDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	gameMap map[string]entity.Game
}

func (mock *mockedDynamoDB) GetItemWithContext(ctx context.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
	primaryKey := input.Key[gamestrg.FieldPrimaryKey]
	if primaryKey == nil {
		return nil, gamestrg.ErrItemNotFound
	}

	gameItem, found := mock.gameMap[*primaryKey.S]
	if !found {
		return nil, gamestrg.ErrItemNotFound
	}

	item, err := dynamodbattribute.MarshalMap(gameItem)
	if err != nil {
		return nil, err
	}

	result := dynamodb.GetItemOutput{
		Item: item,
	}
	return &result, nil
}

func (mock *mockedDynamoDB) PutItemWithContext(ctx context.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
	gameItem := entity.Game{}

	err := dynamodbattribute.UnmarshalMap(input.Item, &gameItem)
	if err != nil {
		return nil, err
	}

	mock.gameMap[gameItem.ID] = gameItem
	return nil, nil
}

func TestSaveGame(t *testing.T) {
	gameMap := make(map[string]entity.Game)
	mockedDynamo := mockedDynamoDB{
		gameMap: gameMap,
	}

	storage := gamestrg.NewDynamoDBStorage(&mockedDynamo)

	newGame := entity.Game{
		ID:         "game-id",
		PlayerName: "Alfat",
		Partner: &entity.Pokemon{
			ID: "partner-id",
		},
		BattleWon: 2,
		Scenario:  entity.BATTLE_3,
	}

	err := storage.SaveGame(context.Background(), newGame)
	if err != nil {
		t.Errorf("function return error: %s", err)
	}

	result := gameMap[newGame.ID]
	if result.ID != newGame.ID {
		t.Error("created object does not match test data ID")
	}

	if result.Partner.ID != newGame.Partner.ID {
		t.Error("created object does not match test data Partner's Attack BattleStats")
	}

	if result.Scenario != newGame.Scenario {
		t.Error("created object does not match test data Scenario")
	}
}

func TestGetGame(t *testing.T) {
	mockedGame := entity.Game{
		ID:         "game-id",
		PlayerName: "Alfat",
		Partner: &entity.Pokemon{
			ID: "partner-id",
		},
		BattleWon: 2,
		Scenario:  entity.BATTLE_3,
	}

	gameMap := map[string]entity.Game{
		mockedGame.ID: mockedGame,
	}

	mockedDynamo := mockedDynamoDB{
		gameMap: gameMap,
	}

	storage := gamestrg.NewDynamoDBStorage(&mockedDynamo)
	result, err := storage.GetGame(context.Background(), "game-id")
	if err != nil {
		t.Errorf("function return error: %s", err)
	}

	if result.ID != mockedGame.ID {
		t.Errorf("fetched object does not return same ID")
	}

	if result.Partner.ID != mockedGame.Partner.ID {
		t.Error("fetched object does not match test data Partner's Attack BattleStats")
	}

	if result.Scenario != mockedGame.Scenario {
		t.Error("fetched object does not match test data Scenario")
	}

}
