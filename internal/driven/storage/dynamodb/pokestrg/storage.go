package pokestrg

import (
	"context"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Role string

const (
	FieldExtraRole  = "extra_role"
	FieldPrimaryKey = "id"

	TableName = "Pokemons"

	PARTNER Role = "PARTNER"
	ENEMY   Role = "ENEMY"
)

var (
	ErrItemNotFound = fmt.Errorf("item cannot be found within table %s", TableName)
)

type DynamoDBStorage struct {
	db dynamodbiface.DynamoDBAPI
}

func (storage *DynamoDBStorage) getPokemonsByRole(ctx context.Context, role Role) ([]entity.Pokemon, error) {
	input := dynamodb.QueryInput{
		TableName:              aws.String(TableName),
		KeyConditionExpression: aws.String("extra_role = :role"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			"role": {
				S: aws.String(string(role)),
			},
		},
	}

	output, err := storage.db.QueryWithContext(ctx, &input)
	if err != nil {
		return nil, err
	}

	results := make([]entity.Pokemon, len(output.Items))
	dynamodbattribute.UnmarshalListOfMaps(output.Items, &results)
	return results, nil
}

func (storage *DynamoDBStorage) GetAvailablePartners(ctx context.Context) ([]entity.Pokemon, error) {
	return storage.getPokemonsByRole(ctx, PARTNER)
}

func (storage *DynamoDBStorage) GetPossibleEnemies(ctx context.Context) ([]entity.Pokemon, error) {
	return storage.getPokemonsByRole(ctx, ENEMY)
}

func (storage *DynamoDBStorage) GetPartner(ctx context.Context, partnerID string) (*entity.Pokemon, error) {
	input := dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			FieldPrimaryKey: {
				S: &partnerID,
			},
			FieldExtraRole: {
				S: aws.String(string(PARTNER)),
			},
		},
	}

	output, err := storage.db.GetItemWithContext(ctx, &input)
	if err != nil {
		return nil, err
	}

	if output.Item == nil {
		return nil, ErrItemNotFound
	}

	partner := entity.Pokemon{}
	err = dynamodbattribute.UnmarshalMap(output.Item, &partner)
	return &partner, err
}

func NewDynamoDBStorage(db dynamodbiface.DynamoDBAPI) DynamoDBStorage {
	return DynamoDBStorage{
		db: db,
	}
}
