package pokestrg_test

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/pokestrg"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockedDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	pokemonMap map[pokestrg.Role]map[string]entity.Pokemon
}

var (
	pikachu = entity.Pokemon{
		ID:   "pikachu-id",
		Name: "Pikachu",
		BattleStats: entity.BattleStats{
			MaxHealth: 100,
			Attack:    25,
			Defense:   10,
			Speed:     10,
		},
		AvatarURL: "http://avatar.co",
	}

	bulbasaur = entity.Pokemon{
		ID:   "bulbasaur-id",
		Name: "Bulbasaur",
		BattleStats: entity.BattleStats{
			MaxHealth: 100,
			Attack:    25,
			Defense:   10,
			Speed:     10,
		},
		AvatarURL: "http://avatar.co",
	}

	databaseMap = map[pokestrg.Role]map[string]entity.Pokemon{
		pokestrg.PARTNER: {
			pikachu.ID: pikachu,
		},
		pokestrg.ENEMY: {
			bulbasaur.ID: bulbasaur,
		},
	}
)

func (mock *mockedDynamoDB) QueryWithContext(ctx context.Context, input *dynamodb.QueryInput, opts ...request.Option) (*dynamodb.QueryOutput, error) {
	roleVal := input.ExpressionAttributeValues["role"]
	pokemonsMapByRole := mock.pokemonMap[pokestrg.Role(*roleVal.S)]

	resultList := make([]map[string]*dynamodb.AttributeValue, 0)
	for _, pokemon := range pokemonsMapByRole {
		item, err := dynamodbattribute.MarshalMap(pokemon)
		if err != nil {
			return nil, err
		}

		resultList = append(resultList, item)
	}

	output := dynamodb.QueryOutput{
		Items: resultList,
	}

	return &output, nil
}

func (mock *mockedDynamoDB) GetItemWithContext(ctx context.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
	pokemonID := input.Key[pokestrg.FieldPrimaryKey]
	extraRole := input.Key[pokestrg.FieldExtraRole]

	pokemon, found := mock.pokemonMap[pokestrg.Role(*extraRole.S)][*pokemonID.S]
	if !found {
		return nil, pokestrg.ErrItemNotFound
	}

	item, err := dynamodbattribute.MarshalMap(pokemon)
	if err != nil {
		return nil, err
	}

	result := dynamodb.GetItemOutput{
		Item: item,
	}

	return &result, nil
}

func TestGetPartner(t *testing.T) {
	mockedDynamo := mockedDynamoDB{
		pokemonMap: databaseMap,
	}

	mockedPartner := databaseMap[pokestrg.PARTNER][pikachu.ID]

	storage := pokestrg.NewDynamoDBStorage(&mockedDynamo)
	result, err := storage.GetPartner(context.Background(), pikachu.ID)
	if err != nil {
		t.Errorf("function return error: %s", err)
	}

	if result.ID != mockedPartner.ID {
		t.Error("fetched object does not return same ID")
	}

	if result.BattleStats.MaxHealth != mockedPartner.BattleStats.MaxHealth {
		t.Error("fetched object does not return same BattleStats")
	}
}

func TestGetPossibleEnemies(t *testing.T) {
	mockedDynamo := mockedDynamoDB{
		pokemonMap: databaseMap,
	}

	storage := pokestrg.NewDynamoDBStorage(&mockedDynamo)
	enemies, err := storage.GetPossibleEnemies(context.Background())
	if err != nil {
		t.Errorf("function return error: %s", err)
	}

	if len(enemies) != len(databaseMap[pokestrg.ENEMY]) {
		t.Error("fetched enemies count does not match test data")
	}

	if enemies[0].ID != bulbasaur.ID {
		t.Errorf("fetched enemy ID does not match test data ")
	}
}

func TestGetAvailableEnemies(t *testing.T) {
	mockedDynamo := mockedDynamoDB{
		pokemonMap: databaseMap,
	}

	storage := pokestrg.NewDynamoDBStorage(&mockedDynamo)
	partners, err := storage.GetAvailablePartners(context.Background())
	if err != nil {
		t.Errorf("function return error: %s", err)
	}

	if len(partners) != len(databaseMap[pokestrg.PARTNER]) {
		t.Error("fetched partners count does not match test data")
	}

	if partners[0].ID != pikachu.ID {
		t.Errorf("fetched partner ID does not match test data ")
	}
}
