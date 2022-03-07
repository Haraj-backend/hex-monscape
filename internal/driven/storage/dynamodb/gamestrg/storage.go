package gamestrg

import (
	"context"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	FieldPrimaryKey = "id"

	TableName = "Games"
)

var (
	ErrItemNotFound = fmt.Errorf("item cannot be found within table %s", TableName)
)

type DynamoDBStorage struct {
	db dynamodbiface.DynamoDBAPI
}

func (storage *DynamoDBStorage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
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

	gameItem := entity.Game{}
	err = dynamodbattribute.UnmarshalMap(output.Item, &gameItem)
	return &gameItem, err
}

func (storage *DynamoDBStorage) SaveGame(ctx context.Context, game entity.Game) error {
	marshalledItem, err := dynamodbattribute.MarshalMap(&game)
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
		db: db,
	}
}
