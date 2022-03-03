package battlestrg_test

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/battlestrg"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockedDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	battleMap map[string]battle.Battle
}

func (mock *mockedDynamoDB) GetItemWithContext(ctx context.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
	primaryKey := input.Key[battlestrg.PrimaryKey]
	if primaryKey == nil {
		return nil, battlestrg.ErrItemNotFound
	}

	battleItem, found := mock.battleMap[*primaryKey.S]
	if !found {
		return nil, battlestrg.ErrItemNotFound
	}

	item, err := dynamodbattribute.MarshalMap(battleItem)
	if err != nil {
		return nil, err
	}

	result := dynamodb.GetItemOutput{
		Item: item,
	}
	return &result, nil
}

func (mock *mockedDynamoDB) PutItemWithContext(ctx context.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
	battleItem := battle.Battle{}

	err := dynamodbattribute.UnmarshalMap(input.Item, &battleItem)
	if err != nil {
		return nil, err
	}

	mock.battleMap[battleItem.GameID] = battleItem
	return nil, nil
}

func TestSaveBattle(t *testing.T) {
	battleMap := make(map[string]battle.Battle)
	mockedDynamo := mockedDynamoDB{
		battleMap: battleMap,
	}

	storage := battlestrg.NewDynamoStorage(&mockedDynamo)

	newBattle := battle.Battle{
		GameID: "something-id",
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

	err := storage.SaveBattle(context.Background(), newBattle)
	if err != nil {
		t.Errorf("function return error: %s", err)
	}

	result := battleMap[newBattle.GameID]
	if result.GameID != newBattle.GameID {
		t.Error("did not store the entity")
	}
}

func TestGetBattle(t *testing.T) {
	mockedBattle := battle.Battle{
		GameID: "something-id",
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

	battleMap := map[string]battle.Battle{
		mockedBattle.GameID: mockedBattle,
	}

	mockedDynamo := mockedDynamoDB{
		battleMap: battleMap,
	}

	storage := battlestrg.NewDynamoStorage(&mockedDynamo)
	result, err := storage.GetBattle(context.Background(), "something-id")
	if err != nil {
		t.Errorf("function return error: %s", err)
	}

	if result.GameID != mockedBattle.GameID {
		t.Error("fetched object does not return same GameID")
	}

	if result.State != mockedBattle.State {
		t.Error("fetched object does not return same State")
	}
}
