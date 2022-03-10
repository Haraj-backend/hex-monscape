package pokestrg

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/shared"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeedData(t *testing.T) {
	storage, err := newStorage()
	require.NoError(t, err)

	partner := newSamplePokemon()
	enemy := newSamplePokemon()

	seeder := PokemonSeeder{
		partners: []entity.Pokemon{*partner},
		enemies:  []entity.Pokemon{*enemy},
	}

	err = storage.SeedData(context.Background(), &seeder)
	require.NoError(t, err)

	seededPartner, err := getPokemon(storage.dynamoClient, partner.ID)
	require.NoError(t, err)
	require.Equal(t, *partner, *seededPartner)

	seededEnemy, err := getPokemon(storage.dynamoClient, enemy.ID)
	require.NoError(t, err)
	require.Equal(t, *enemy, *seededEnemy)

	err = deletePokemon(storage.dynamoClient, partner.ID)
	require.NoError(t, err)

	err = deletePokemon(storage.dynamoClient, enemy.ID)
	require.NoError(t, err)
}

func TestGetPartner(t *testing.T) {
	storage, err := newStorage()
	require.NoError(t, err)

	partner := newSamplePokemon()
	seeder := PokemonSeeder{
		partners: []entity.Pokemon{*partner},
	}

	err = storage.SeedData(context.Background(), &seeder)
	require.NoError(t, err)

	testCases := []struct {
		Name       string
		PartnerID  string
		ExpPartner *entity.Pokemon
	}{
		{
			Name:       "Test Partner Not Found",
			PartnerID:  uuid.NewString(),
			ExpPartner: nil,
		},
		{
			Name:       "Test Partner Found",
			PartnerID:  partner.ID,
			ExpPartner: partner,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			pkmn, err := storage.GetPartner(context.Background(), testCase.PartnerID)
			assert.NoError(t, err)

			if testCase.ExpPartner == nil {
				return
			}

			assert.Equal(t, *pkmn, *partner)
		})
	}

	err = deletePokemon(storage.dynamoClient, partner.ID)
	require.NoError(t, err)
}

func TestGetPossibleEnemies(t *testing.T) {
	storage, err := newStorage()
	require.NoError(t, err)

	enemy := newSamplePokemon()

	testCases := []struct {
		Name    string
		Enemies []entity.Pokemon
	}{
		{
			Name:    "Test Empty Possible Enemies",
			Enemies: []entity.Pokemon{},
		},
		{
			Name:    "Test Exists Possible Enemies",
			Enemies: []entity.Pokemon{*enemy},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			seeder := PokemonSeeder{
				enemies: testCase.Enemies,
			}

			err := storage.SeedData(context.Background(), &seeder)
			assert.NoError(t, err)

			enemies, err := storage.GetPossibleEnemies(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, len(testCase.Enemies), len(enemies))

			if len(testCase.Enemies) == 1 && len(enemies) == 1 {
				assert.Equal(t, testCase.Enemies[0], enemies[0])
				deletePokemon(storage.dynamoClient, enemy.ID)
			}
		})
	}
}

func TestGetAvailablePartners(t *testing.T) {
	storage, err := newStorage()
	require.NoError(t, err)

	partner := newSamplePokemon()

	testCases := []struct {
		Name     string
		Partners []entity.Pokemon
	}{
		{
			Name:     "Test Empty Available Partners",
			Partners: []entity.Pokemon{},
		},
		{
			Name:     "Test Exists Availabel Partners",
			Partners: []entity.Pokemon{*partner},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			seeder := PokemonSeeder{
				partners: testCase.Partners,
			}

			err := storage.SeedData(context.Background(), &seeder)
			assert.NoError(t, err)

			partners, err := storage.GetAvailablePartners(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, len(testCase.Partners), len(partners))

			if len(testCase.Partners) == 1 && len(partners) == 1 {
				assert.Equal(t, testCase.Partners[0], partners[0])
				deletePokemon(storage.dynamoClient, partner.ID)
			}
		})
	}
}

func newStorage() (*Storage, error) {
	storage, err := New(Config{
		DynamoClient: shared.NewLocalTestDDBClient(),
		TableName:    os.Getenv(shared.TestConfig.EnvKeyPokemonTableName),
	})

	if err != nil {
		return nil, fmt.Errorf("unable to initialize pokemon storage due: %w", err)
	}

	return storage, nil
}

func newSamplePokemon() *entity.Pokemon {
	currentTs := time.Now().Unix()
	return &entity.Pokemon{
		ID:   uuid.NewString(),
		Name: fmt.Sprintf("pokemon_%v", currentTs),
		BattleStats: entity.BattleStats{
			Health:    100,
			MaxHealth: 100,
			Attack:    25,
			Defense:   10,
			Speed:     20,
		},
		AvatarURL: fmt.Sprintf("http://example.com/%v", currentTs),
	}
}

func getPokemon(dynamoClient *dynamodb.DynamoDB, ID string) (*entity.Pokemon, error) {
	output, err := dynamoClient.GetItem(&dynamodb.GetItemInput{
		Key:       (pokemonKey{ID: ID}).toDDBKey(),
		TableName: aws.String(os.Getenv(shared.TestConfig.EnvKeyPokemonTableName)),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get pokemon due: %w", err)
	}

	if len(output.Item) == 0 {
		return nil, fmt.Errorf("pokemon is not found")
	}

	var p entity.Pokemon
	err = dynamodbattribute.UnmarshalMap(output.Item, &p)
	if err != nil {
		return nil, fmt.Errorf("unable to parse item due: %w", err)
	}

	return &p, nil
}

func deletePokemon(dynamoClient *dynamodb.DynamoDB, ID string) error {
	output, err := dynamoClient.DeleteItem(&dynamodb.DeleteItemInput{
		Key:       (pokemonKey{ID: ID}).toDDBKey(),
		TableName: aws.String(os.Getenv(shared.TestConfig.EnvKeyPokemonTableName)),
		// using this to make sure if the item is surely deleted
		// reference: https://stackoverflow.com/questions/46464303/how-to-determine-if-a-dynamodb-item-was-indeed-deleted
		ReturnValues: aws.String("ALL_OLD"),
	})
	if err != nil {
		return fmt.Errorf("unable to remove pokemon due: %w", err)
	}

	if len(output.Attributes) == 0 {
		return fmt.Errorf("pokemon is not deleted")
	}

	return nil
}
