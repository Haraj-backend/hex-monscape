package battlestrg

import (
	"context"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	FieldPrimaryKey = "game_id"

	TableName = "Battles"
)

var (
	ErrItemNotFound = fmt.Errorf("item cannot be found within table %s", TableName)
)

type DynamoDBStorage struct {
	db dynamodbiface.DynamoDBAPI
}

func (storage *DynamoDBStorage) GetBattle(ctx context.Context, gameID string) (*battle.Battle, error) {
	input := dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			FieldPrimaryKey: {
				S: &gameID,
			},
		},
	}

	output, err := storage.db.GetItemWithContext(ctx, &input)
	if err != nil {
		return nil, err
	}

	if output.Item == nil {
		return nil, ErrItemNotFound
	}

	battleItem := battle.Battle{}
	err = dynamodbattribute.UnmarshalMap(output.Item, &battleItem)
	return &battleItem, err
}

func (storage *DynamoDBStorage) SaveBattle(ctx context.Context, b battle.Battle) error {
	marshalledItem, err := dynamodbattribute.MarshalMap(&b)
	if err != nil {
		return err
	}

	input := dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      marshalledItem,
	}

	_, err = storage.db.PutItemWithContext(ctx, &input)
	return err
}

func NewDynamoDBStorage(db dynamodbiface.DynamoDBAPI) *DynamoDBStorage {
	return &DynamoDBStorage{
		db,
	}
}
