package gamestrg

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

func (s *Storage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	key := gameKey{ID: gameID}
	input := dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key:       key.toDDBKey(),
	}

	output, err := s.dynamoClient.GetItemWithContext(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("unable to get item from %s due to: %w", s.tableName, err)
	}

	if len(output.Item) == 0 {
		return nil, nil
	}

	var gameRow gameRow
	err = dynamodbattribute.UnmarshalMap(output.Item, &gameRow)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal item from %s due to: %w", s.tableName, err)
	}

	return gameRow.toGame(), nil
}

func (s *Storage) SaveGame(ctx context.Context, game entity.Game) error {
	item, _ := dynamodbattribute.MarshalMap(toGameRow(game))
	input := dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      item,
	}

	_, err := s.dynamoClient.PutItemWithContext(ctx, &input)
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

// New returns new instance of gamestrg dynamoDB Storage
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
