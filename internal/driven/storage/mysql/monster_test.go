package mysql

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var partners, enemies []entity.Monster

func init() {
	// init mysql client
	sqlClient, err := newTestSQLClient()
	if err != nil {
		panic(err)
	}
	// seed data to mysql
	rows := []monsterRow{
		newTestMonsterRow(true),
		newTestMonsterRow(true),
		newTestMonsterRow(false),
		newTestMonsterRow(false),
	}
	for _, row := range rows {
		// insert monster row to mysql
		err := insertMonster(sqlClient, row)
		if err != nil {
			panic(err)
		}
		// assign partners & enemies for later tests
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
	require.Subset(t, availablePartners, partners, "availablePartners should contains partners") // we use subset here because in other tests we might need to insert monsters as well
}

func TestGetPossibleEnemies(t *testing.T) {
	storage := newMonsterStorage(t)

	availableEnemies, err := storage.GetPossibleEnemies(context.Background())
	require.NoError(t, err)
	require.Subset(t, availableEnemies, enemies, "availableEnemies should contains enemies") // we use subset here because in other tests we might need to insert monsters as well
}

func TestGetPartner(t *testing.T) {
	storage := newMonsterStorage(t)

	testCases := []struct {
		Name       string
		PartnerID  string
		ExpPartner *entity.Monster
	}{
		{
			Name:       "Partner Exists",
			PartnerID:  partners[0].ID,
			ExpPartner: &partners[0],
		},
		{
			Name:       "Partner Not Exists",
			PartnerID:  uuid.NewString(),
			ExpPartner: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			partner, err := storage.GetPartner(context.Background(), testCase.PartnerID)
			require.NoError(t, err)
			require.Equal(t, testCase.ExpPartner, partner)
		})
	}
}

func newMonsterStorage(t *testing.T) *MonsterStorage {
	// init mysql client
	sqlClient, err := newTestSQLClient()
	require.NoError(t, err)

	// init storage
	storage, err := NewMonsterStorage(MonsterStorageConfig{SQLClient: sqlClient})
	require.NoError(t, err)

	return storage
}
