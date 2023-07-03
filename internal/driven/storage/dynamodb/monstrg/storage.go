package monstrg

import (
	"context"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/shared"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gopkg.in/validator.v2"
)

type Storage struct {
	dynamoClient *dynamodb.DynamoDB
	tableName    string
}

func (s *Storage) getPokemonsByRole(ctx context.Context, extraRole extraRole) ([]entity.Monster, error) {
	query := monsterExtraRoleQuery{
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
		return nil, fmt.Errorf("unable to query from %s due to: %w", s.tableName, err)
	}

	if len(output.Items) == 0 {
		return nil, nil
	}

	rows := make([]shared.MonsterRow, len(output.Items))
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &rows)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal items from %s due to: %w", s.tableName, err)
	}
	monsters := make([]entity.Monster, len(rows))
	for i, row := range rows {
		m := row.ToMonster()
		monsters[i] = *m
	}
	return monsters, nil
}

func (s *Storage) GetAvailablePartners(ctx context.Context) ([]entity.Monster, error) {
	return s.getPokemonsByRole(ctx, partnerRole)
}

func (s *Storage) GetPossibleEnemies(ctx context.Context) ([]entity.Monster, error) {
	return s.getPokemonsByRole(ctx, enemyRole)
}

func (s *Storage) GetPartner(ctx context.Context, partnerID string) (*entity.Monster, error) {
	key := monsterKey{
		ID: partnerID,
	}

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

	var monsterRow shared.MonsterRow
	err = dynamodbattribute.UnmarshalMap(output.Item, &monsterRow)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal item from %s due to: %w", s.tableName, err)
	}

	return monsterRow.ToMonster(), nil
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
