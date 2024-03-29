package battlestrg

import (
	"context"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gopkg.in/validator.v2"
)

type Storage struct {
	dynamoClient *dynamodb.DynamoDB
	tableName    string
}

func (s *Storage) GetBattle(ctx context.Context, gameID string) (*entity.Battle, error) {
	// construct params
	key := battleKey{GameID: gameID}
	input := &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key:       key.toDDBKey(),
	}
	// execute get item
	output, err := s.dynamoClient.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("unable to get item from %s due to: %w", s.tableName, err)
	}
	// if item is not found, returns nil as expected by battle interface
	if len(output.Item) == 0 {
		return nil, nil
	}
	// parse item
	var battleRow battleRow
	err = dynamodbattribute.UnmarshalMap(output.Item, &battleRow)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal item from %s due to: %w", s.tableName, err)
	}
	battle := battleRow.toBattle()

	return battle, nil
}

func (s *Storage) SaveBattle(ctx context.Context, b entity.Battle) error {
	// construct params
	item, _ := dynamodbattribute.MarshalMap(toBattleRow(b))
	input := &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      item,
	}
	// execute put item
	_, err := s.dynamoClient.PutItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("unable to put item to %s due to: %w", s.tableName, err)
	}

	return nil
}

type Config struct {
	DynamoClient *dynamodb.DynamoDB `validate:"nonnil"`
	TableName    string             `validate:"nonzero"`
}

func (c Config) Validate() error {
	return validator.Validate(c)
}

// New returns new instance of battlestrg dynamoDB Storage
func New(cfg Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	strg := &Storage{
		dynamoClient: cfg.DynamoClient,
		tableName:    cfg.TableName,
	}
	return strg, nil
}
