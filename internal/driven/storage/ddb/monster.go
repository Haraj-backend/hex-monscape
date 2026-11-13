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

type MonsterStorage struct {
	dynamoClient *dynamodb.DynamoDB
	tableName    string
}

type MonsterStorageConfig struct {
	DynamoClient *dynamodb.DynamoDB `validate:"nonnil"`
	TableName    string             `validate:"nonzero"`
}

func (c MonsterStorageConfig) Validate() error {
	return validator.Validate(c)
}

func NewMonsterStorage(cfg MonsterStorageConfig) (*MonsterStorage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	strg := &MonsterStorage{
		dynamoClient: cfg.DynamoClient,
		tableName:    cfg.TableName,
	}

	return strg, nil
}

func (s *MonsterStorage) GetAvailablePartners(ctx context.Context) ([]entity.Monster, error) {
	// query partnerable monsters
	eav, err := dynamodbattribute.MarshalMap(map[string]interface{}{
		":is_partnerable": 1,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal expression attribute values due: %w", err)
	}
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

	return toMonsters(output.Items)
}

func (s *MonsterStorage) GetPossibleEnemies(ctx context.Context) ([]entity.Monster, error) {
	// scan whole monster table, every monster can be an enemy
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

	return toMonsters(output.Items)
}

func (s *MonsterStorage) GetPartner(ctx context.Context, partnerID string) (*entity.Monster, error) {
	// get item from dynamodb
	output, err := s.dynamoClient.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key:       monsterKey{ID: partnerID}.toDDBKey(),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get item from %s due to: %w", s.tableName, err)
	}

	// if item is not found, returns nil as expected by partner interface
	if len(output.Item) == 0 {
		return nil, nil
	}

	// parse item
	var row monsterRow
	err = dynamodbattribute.UnmarshalMap(output.Item, &row)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal item from %s due to: %w", s.tableName, err)
	}

	return row.toMonster(), nil
}
