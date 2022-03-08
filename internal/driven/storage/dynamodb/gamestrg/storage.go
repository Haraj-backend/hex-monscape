package gamestrg

import (
	"context"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	FieldPrimaryKey = "id"

	TableName = "Games"
)

type dynamoDBStorage struct {
	db dynamodbiface.DynamoDBAPI
}

func (storage *dynamoDBStorage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	key, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
		FieldPrimaryKey: gameID,
	})

	input := dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key:       key,
	}

	output, err := storage.db.GetItemWithContext(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("unable to get item from %s due to: %w", TableName, err)
	}

	gameItem := entity.Game{}
	err = dynamodbattribute.UnmarshalMap(output.Item, &gameItem)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal item from %s due to: %w", TableName, err)
	}

	return &gameItem, nil
}

func (storage *dynamoDBStorage) SaveGame(ctx context.Context, game entity.Game) error {
	marshalledItem, _ := dynamodbattribute.MarshalMap(&game)

	input := dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      marshalledItem,
	}

	_, err := storage.db.PutItemWithContext(ctx, &input)
	if err != nil {
		return fmt.Errorf("unable to put item to %s due to: %w", TableName, err)
	}

	return nil
}

// NewDynamoDBStorage returns new instance of gamestrg dynamoDBStorage
func NewDynamoDBStorage(cfg config.DynamoDBStorageConfig) (*dynamoDBStorage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	strg := &dynamoDBStorage{
		db: cfg.DB,
	}

	return strg, nil
}
