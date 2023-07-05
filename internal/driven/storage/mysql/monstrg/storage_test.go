package monstrg_test

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/monstrg"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/shared"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var partners, enemies []entity.Monster

func init() {
	// init mysql client
	sqlClient, err := shared.NewTestSQLClient()
	if err != nil {
		panic(err)
	}
	// seed data to mysql
	rows := []shared.MonsterRow{
		shared.NewTestMonsterRow(true),
		shared.NewTestMonsterRow(true),
		shared.NewTestMonsterRow(false),
		shared.NewTestMonsterRow(false),
	}
	for _, row := range rows {
		// insert monster row to mysql
		err := shared.InsertMonster(sqlClient, row)
		if err != nil {
			panic(err)
		}
		// assign partners & enemies for later tests
		if row.IsPartnerable == 1 {
			partners = append(partners, *row.ToMonster())
		}
		// we also include partnerable monsters to enemies so the enemies can be more vary
		enemies = append(enemies, *row.ToMonster())
	}

}

func TestGetAvailablePartners(t *testing.T) {
	storage := newStorage(t)

	availablePartners, err := storage.GetAvailablePartners(context.Background())
	require.NoError(t, err)
	require.Subset(t, availablePartners, partners, "availablePartners should contains partners") // we use subset here because in other tests we might need to insert monsters as well
}

func TestGetPossibleEnemies(t *testing.T) {
	storage := newStorage(t)

	availableEnemies, err := storage.GetPossibleEnemies(context.Background())
	require.NoError(t, err)
	require.Subset(t, availableEnemies, enemies, "availableEnemies should contains enemies") // we use subset here because in other tests we might need to insert monsters as well
}

func TestGetPartner(t *testing.T) {
	storage := newStorage(t)

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

func newStorage(t *testing.T) *monstrg.Storage {
	// init mysql client
	sqlClient, err := shared.NewTestSQLClient()
	require.NoError(t, err)

	// init storage
	storage, err := monstrg.New(monstrg.Config{SQLClient: sqlClient})
	require.NoError(t, err)

	return storage
}
