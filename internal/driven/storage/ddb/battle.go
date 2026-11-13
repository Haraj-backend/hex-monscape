package ddb

import (
	"context"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gopkg.in/validator.v2"
)

type BattleStorage struct {
	dynamoClient *dynamodb.DynamoDB
	tableName    string
}

type BattleStorageConfig struct {
	DynamoClient *dynamodb.DynamoDB `validate:"nonnil"`
	TableName    string             `validate:"nonzero"`
}

func (c BattleStorageConfig) Validate() error {
	return validator.Validate(c)
}

func NewBattleStorage(cfg BattleStorageConfig) (*BattleStorage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	strg := &BattleStorage{
		dynamoClient: cfg.DynamoClient,
		tableName:    cfg.TableName,
	}

	return strg, nil
}

func (s *BattleStorage) GetBattle(ctx context.Context, gameID string) (*entity.Battle, error) {
	// get item from dynamodb
	input := &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key:       battleKey{GameID: gameID}.toDDBKey(),
	}
	output, err := s.dynamoClient.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("unable to get item from %s due to: %w", s.tableName, err)
	}

	// if item is not found, returns nil as expected by battle interface
	if len(output.Item) == 0 {
		return nil, nil
	}

	// parse item
	var row battleRow
	err = dynamodbattribute.UnmarshalMap(output.Item, &row)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal item from %s due to: %w", s.tableName, err)
	}

	return row.toBattle(), nil
}

func (s *BattleStorage) SaveBattle(ctx context.Context, b entity.Battle) error {
	// construct item
	item, err := dynamodbattribute.MarshalMap(newBattleRow(b))
	if err != nil {
		return fmt.Errorf("unable to marshal battle row due: %w", err)
	}

	// put item to dynamodb
	input := &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      item,
	}
	_, err = s.dynamoClient.PutItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("unable to put item to %s due to: %w", s.tableName, err)
	}

	return nil
}
