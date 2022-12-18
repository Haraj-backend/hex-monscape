package gamestrg

import (
	"context"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/shared/telemetry"
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
	tr := telemetry.GetTracer()
	_, span := tr.Trace(ctx, "GetGame GameStorage")
	defer span.End()

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

	gameItem := entity.Game{}
	err = dynamodbattribute.UnmarshalMap(output.Item, &gameItem)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal item from %s due to: %w", s.tableName, err)
	}

	return &gameItem, nil
}

func (s *Storage) SaveGame(ctx context.Context, game entity.Game) error {
	tr := telemetry.GetTracer()
	_, span := tr.Trace(ctx, "SaveGame GameStorage")
	defer span.End()

	item, _ := dynamodbattribute.MarshalMap(&game)
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
