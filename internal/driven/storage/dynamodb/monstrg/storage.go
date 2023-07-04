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

// New returns new instance of pokestrg dynamoDB Storage
func New(cfg Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	strg := &Storage{
		dynamoClient: cfg.DynamoClient,
		tableName:    cfg.TableName,
	}
	return strg, nil
}

type Storage struct {
	dynamoClient *dynamodb.DynamoDB
	tableName    string
}

func (s *Storage) GetAvailablePartners(ctx context.Context) ([]entity.Monster, error) {
	// query partnerable monsters
	eav, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
		":is_partnerable": 1,
	})
	output, err := s.dynamoClient.QueryWithContext(ctx, &dynamodb.QueryInput{
		KeyConditionExpression:    aws.String("is_partnerable = :is_partnerable"),
		ExpressionAttributeValues: eav,
		IndexName:                 aws.String("is_partnerable"),
		TableName:                 aws.String(s.tableName),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to query table %s due to: %w", s.tableName, err)
	}
	// just return if empty
	if len(output.Items) == 0 {
		return nil, nil
	}
	// parse monster rows
	return toMonsters(output.Items)
}

func (s *Storage) GetPossibleEnemies(ctx context.Context) ([]entity.Monster, error) {
	// scan whole monster table
	output, err := s.dynamoClient.ScanWithContext(ctx, &dynamodb.ScanInput{
		TableName: aws.String(s.tableName),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to scan table %s due to: %w", s.tableName, err)
	}
	// just return if empty
	if len(output.Items) == 0 {
		return nil, nil
	}
	// parse monster rows
	return toMonsters(output.Items)
}

func (s *Storage) GetPartner(ctx context.Context, partnerID string) (*entity.Monster, error) {
	key := monsterKey{ID: partnerID}.toDDBKey()
	output, err := s.dynamoClient.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key:       key,
	})
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

type Config struct {
	DynamoClient *dynamodb.DynamoDB `validate:"nonnil"`
	TableName    string             `validate:"nonzero"`
}

func (c Config) Validate() error {
	return validator.Validate(c)
}
