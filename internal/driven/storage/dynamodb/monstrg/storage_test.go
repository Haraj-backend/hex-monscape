package monstrg

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/shared"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/testutil"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/require"
)

var partners, enemies []entity.Monster

func init() {
	// seed data to dynamodb
	ddbClient := shared.NewLocalTestDDBClient()
	monsterRows := []shared.MonsterRow{
		newMonsterRow(true),
		newMonsterRow(true),
		newMonsterRow(false),
		newMonsterRow(false),
	}
	for _, monsterRow := range monsterRows {
		// put monster row to dynamodb
		item, _ := dynamodbattribute.MarshalMap(monsterRow)
		_, err := ddbClient.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv(shared.TestConfig.EnvKeyPokemonTableName)),
			Item:      item,
		})
		if err != nil {
			log.Printf("[RDebug] item: %+v", item)
			panic(err)
		}
		// append monster row to partners & enemies for later tests
		if monsterRow.IsPartnerable == 1 {
			partners = append(partners, *monsterRow.ToMonster())
		}
		// we also include partnerable monsters to enemies so the enemies can be more vary
		enemies = append(enemies, *monsterRow.ToMonster())
	}
}

func TestGetAvailablePartners(t *testing.T) {
	storage := newStorage(t)

	availablePartners, err := storage.GetAvailablePartners(context.Background())
	require.NoError(t, err)
	require.ElementsMatch(t, partners, availablePartners)
}

func TestGetPossibleEnemies(t *testing.T) {
	storage := newStorage(t)

	possibleEnemies, err := storage.GetPossibleEnemies(context.Background())
	require.NoError(t, err)
	require.ElementsMatch(t, enemies, possibleEnemies)
}

func TestGetPartner(t *testing.T) {
	storage := newStorage(t)

	// define test cases
	testCases := []struct {
		Name       string
		PartnerID  string
		ExpPartner *entity.Monster
	}{
		{
			Name:       "Partner Found",
			PartnerID:  partners[0].ID,
			ExpPartner: &partners[0],
		},
		{
			Name:       "Partner Not Found",
			PartnerID:  "partner-not-found",
			ExpPartner: nil,
		},
	}

	// run test cases
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			partner, err := storage.GetPartner(context.Background(), tc.PartnerID)
			require.NoError(t, err)
			require.Equal(t, tc.ExpPartner, partner)
		})
	}
}

func newStorage(t *testing.T) *Storage {
	storage, err := New(Config{
		DynamoClient: shared.NewLocalTestDDBClient(),
		TableName:    os.Getenv(shared.TestConfig.EnvKeyPokemonTableName),
	})
	require.NoError(t, err)

	return storage
}

func newMonsterRow(isPartnerable bool) shared.MonsterRow {
	row := shared.ToMonsterRow(*(testutil.NewTestMonster()))
	if isPartnerable {
		row.IsPartnerable = 1
	}
	return row
}
