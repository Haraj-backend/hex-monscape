package gamestrg

import (
	"context"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/shared/telemetry"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/validator.v2"
)

type Storage struct {
	dynamoClient *dynamodb.DynamoDB
	tableName    string
}

func (s *Storage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "GameStorage: GetGame", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(attribute.Key("game-id").String(gameID))

	key := gameKey{ID: gameID}
	input := dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key:       key.toDDBKey(),
	}

	output, err := s.dynamoClient.GetItemWithContext(ctx, &input)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Key("error").Bool(true))

		return nil, fmt.Errorf("unable to get item from %s due to: %w", s.tableName, err)
	}

	if len(output.Item) == 0 {
		return nil, nil
	}

	gameItem := entity.Game{}
	err = dynamodbattribute.UnmarshalMap(output.Item, &gameItem)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Key("error").Bool(true))

		return nil, fmt.Errorf("unable to unmarshal item from %s due to: %w", s.tableName, err)
	}

	return &gameItem, nil
}

func (s *Storage) SaveGame(ctx context.Context, game entity.Game) error {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "GameStorage: SaveGame", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(attribute.Key("game-id").String(game.ID))

	item, _ := dynamodbattribute.MarshalMap(&game)
	input := dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      item,
	}

	_, err := s.dynamoClient.PutItemWithContext(ctx, &input)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Key("error").Bool(true))

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
