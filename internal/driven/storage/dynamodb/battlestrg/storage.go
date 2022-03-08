package battlestrg

import (
	"context"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	FieldPrimaryKey = "game_id"

	TableName = "Battles"
)

type dynamoDBStorage struct {
	db dynamodbiface.DynamoDBAPI
}

func (storage *dynamoDBStorage) GetBattle(ctx context.Context, gameID string) (*battle.Battle, error) {
	key, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
		FieldPrimaryKey: gameID,
	})

	input := dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key:       key,
	}

	output, err := storage.db.GetItemWithContext(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("unable tp get item from %s due to: %w", TableName, err)
	}

	battleItem := battle.Battle{}
	err = dynamodbattribute.UnmarshalMap(output.Item, &battleItem)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal item from %s due to: %w", TableName, err)
	}

	return &battleItem, nil
}

func (storage *dynamoDBStorage) SaveBattle(ctx context.Context, b battle.Battle) error {
	marshalledItem, _ := dynamodbattribute.MarshalMap(&b)

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

// NewDynamoDBStorage returns new instance of battlestrg dynamoDBStorage
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
