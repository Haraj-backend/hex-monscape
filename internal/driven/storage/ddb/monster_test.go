package ddb

import (
	"context"
	"os"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/require"
)

var partners, enemies []entity.Monster

func init() {
	// seed data to dynamodb
	ddbClient := newTestDDBClient()
	rows := []monsterRow{
		newTestMonsterRow(true),
		newTestMonsterRow(true),
		newTestMonsterRow(false),
		newTestMonsterRow(false),
	}
	for _, row := range rows {
		// put monster row to dynamodb
		item, _ := dynamodbattribute.MarshalMap(row)
		_, err := ddbClient.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv(envKeyMonsterTableName)),
			Item:      item,
		})
		if err != nil {
			panic(err)
		}
		// append monster row to partners & enemies for later tests
		if row.IsPartnerable == 1 {
			partners = append(partners, *row.toMonster())
		}
		// we also include partnerable monsters to enemies so the enemies can be more vary
		enemies = append(enemies, *row.toMonster())
	}
}

func TestGetAvailablePartners(t *testing.T) {
	storage := newMonsterStorage(t)

	availablePartners, err := storage.GetAvailablePartners(context.Background())
	require.NoError(t, err)
	require.ElementsMatch(t, partners, availablePartners)
}

func TestGetPossibleEnemies(t *testing.T) {
	storage := newMonsterStorage(t)

	possibleEnemies, err := storage.GetPossibleEnemies(context.Background())
	require.NoError(t, err)
	require.ElementsMatch(t, enemies, possibleEnemies)
}

func TestGetPartner(t *testing.T) {
	storage := newMonsterStorage(t)

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

func newMonsterStorage(t *testing.T) *MonsterStorage {
	storage, err := NewMonsterStorage(MonsterStorageConfig{
		DynamoClient: newTestDDBClient(),
		TableName:    os.Getenv(envKeyMonsterTableName),
	})
	require.NoError(t, err)

	return storage
}

func newTestMonsterRow(isPartnerable bool) monsterRow {
	row := newMonsterRow(*(testutil.NewTestMonster()))
	if isPartnerable {
		row.IsPartnerable = 1
	}
	return row
}
