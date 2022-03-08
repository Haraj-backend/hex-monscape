package pokestrg

import (
	"context"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Role string

const (
	FieldExtraRole  = "extra_role"
	FieldPrimaryKey = "id"
	IndexExtraRole  = "index_extra_role"

	TableName = "Pokemons"

	PARTNER Role = "PARTNER"
	ENEMY   Role = "ENEMY"
)

type dynamoDBStorage struct {
	db dynamodbiface.DynamoDBAPI
}

func (storage *dynamoDBStorage) getPokemonsByRole(ctx context.Context, role Role) ([]entity.Pokemon, error) {
	exprAttrVals, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
		":role": aws.String(string(role)),
	})

	input := dynamodb.QueryInput{
		TableName:                 aws.String(TableName),
		KeyConditionExpression:    aws.String("extra_role = :role"),
		ExpressionAttributeValues: exprAttrVals,
		IndexName:                 aws.String(IndexExtraRole),
	}

	output, err := storage.db.QueryWithContext(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("unable to query from %s due to: %w", TableName, err)
	}

	results := make([]entity.Pokemon, len(output.Items))
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &results)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal items from %s due to: %w", TableName, err)
	}

	return results, nil
}

func (storage *dynamoDBStorage) GetAvailablePartners(ctx context.Context) ([]entity.Pokemon, error) {
	return storage.getPokemonsByRole(ctx, PARTNER)
}

func (storage *dynamoDBStorage) GetPossibleEnemies(ctx context.Context) ([]entity.Pokemon, error) {
	return storage.getPokemonsByRole(ctx, ENEMY)
}

func (storage *dynamoDBStorage) GetPartner(ctx context.Context, partnerID string) (*entity.Pokemon, error) {
	key, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
		FieldPrimaryKey: partnerID,
	})

	input := dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key:       key,
	}

	output, err := storage.db.GetItemWithContext(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("unable to get item from %s due to: %w", TableName, err)
	}

	partner := entity.Pokemon{}
	err = dynamodbattribute.UnmarshalMap(output.Item, &partner)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal item from %s due to: %w", TableName, err)
	}

	return &partner, nil
}

func pokemonListToWriteRequests(pokemons []entity.Pokemon, role Role) []*dynamodb.WriteRequest {
	writeRequests := make([]*dynamodb.WriteRequest, 0)

	for _, pokemon := range pokemons {
		item, _ := dynamodbattribute.MarshalMap(pokemon)
		item["extra_role"], _ = dynamodbattribute.Marshal(string(role))

		putRequest := dynamodb.PutRequest{
			Item: item,
		}

		writeRequest := dynamodb.WriteRequest{
			PutRequest: &putRequest,
		}

		writeRequests = append(writeRequests, &writeRequest)
	}

	return writeRequests
}

func (storage *dynamoDBStorage) seedData(ctx context.Context, partners []entity.Pokemon, enemies []entity.Pokemon) error {
	writeRequests := make([]*dynamodb.WriteRequest, 0)
	writeRequests = append(writeRequests, pokemonListToWriteRequests(partners, PARTNER)...)
	writeRequests = append(writeRequests, pokemonListToWriteRequests(enemies, ENEMY)...)

	input := dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			TableName: writeRequests,
		},
	}

	_, err := storage.db.BatchWriteItemWithContext(ctx, &input)
	if err != nil {
		return fmt.Errorf("unable to batch write item to %s due to: %w", TableName, err)
	}

	return nil
}

type Config struct {
	Partners []entity.Pokemon `validate:"min=1"`
	Enemies  []entity.Pokemon `validate:"min=1"`

	config.DynamoDBStorageConfig
}

// NewDynamoDBStorage returns new instance of pokestrg dynamoDBStorage
func NewDynamoDBStorage(cfg Config) (*dynamoDBStorage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	strg := &dynamoDBStorage{
		db: cfg.DB,
	}

	err = strg.seedData(context.Background(), cfg.Partners, cfg.Enemies)
	if err != nil {
		return nil, err
	}

	return strg, nil
}
