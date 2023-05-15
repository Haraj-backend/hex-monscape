package pokestrg

import (
	"context"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/shared/telemetry"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/validator.v2"
)

type Storage struct {
	dynamoClient *dynamodb.DynamoDB
	tableName    string
}

func (s *Storage) getPokemonsByRole(ctx context.Context, extraRole extraRole) ([]entity.Pokemon, error) {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "PokeStorage: getPokemonsByRole", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	query := pokemonExtraRoleQuery{
		ExtraRole: extraRole,
	}

	input := dynamodb.QueryInput{
		TableName:                 aws.String(s.tableName),
		KeyConditionExpression:    query.toQueryExpression(),
		ExpressionAttributeValues: query.toQueryExpressionValue(),
		IndexName:                 aws.String(indexExtraRole),
	}

	output, err := s.dynamoClient.QueryWithContext(ctx, &input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, fmt.Errorf("unable to query from %s due to: %w", s.tableName, err)
	}

	if len(output.Items) == 0 {
		return nil, nil
	}

	results := make([]entity.Pokemon, len(output.Items))
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &results)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, fmt.Errorf("unable to unmarshal items from %s due to: %w", s.tableName, err)
	}

	return results, nil
}

func (s *Storage) GetAvailablePartners(ctx context.Context) ([]entity.Pokemon, error) {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "PokeStorage: GetAvailablePartners", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	return s.getPokemonsByRole(ctx, partnerRole)
}

func (s *Storage) GetPossibleEnemies(ctx context.Context) ([]entity.Pokemon, error) {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "PokeStorage: GetPossibleEnemies", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	return s.getPokemonsByRole(ctx, enemyRole)
}

func (s *Storage) GetPartner(ctx context.Context, partnerID string) (*entity.Pokemon, error) {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "PokeStorage: GetPartner", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(attribute.Key("partner-id").String(partnerID))

	key := pokemonKey{
		ID: partnerID,
	}

	input := dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key:       key.toDDBKey(),
	}

	output, err := s.dynamoClient.GetItemWithContext(ctx, &input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, fmt.Errorf("unable to get item from %s due to: %w", s.tableName, err)
	}

	if len(output.Item) == 0 {
		return nil, nil
	}

	partner := entity.Pokemon{}
	err = dynamodbattribute.UnmarshalMap(output.Item, &partner)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, fmt.Errorf("unable to unmarshal item from %s due to: %w", s.tableName, err)
	}

	return &partner, nil
}

func (s *Storage) SeedData(ctx context.Context, seeder *PokemonSeeder) error {
	if seeder.isEmpty() {
		return nil
	}

	input := seeder.toBatchWriteInput(s.tableName)
	_, err := s.dynamoClient.BatchWriteItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("unable to batch write item to %s due to: %w", s.tableName, err)
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

// New returns new instance of pokestrg dynamoDB Storage
func New(cfg Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	strg := &Storage{
		dynamoClient: cfg.DynamoClient,
		tableName:    cfg.TableName,
	}

	// here I'm not mimicking the memory storage for seeding the data when constructing the instance
	// instead, I would like to call the SeedData method that I provide on line 94 in the main function, WDYT?

	return strg, nil
}
