package monstrg

import (
	"context"
	"os"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/shared"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSeedDataGetPartnersGetEnemies(t *testing.T) {
	// initialize storage
	storage := newStorage(t)

	// seed data
	partner := testutil.NewTestMonster()
	enemy := testutil.NewTestMonster()

	seeder := PokemonSeeder{
		partners: []entity.Monster{*partner},
		enemies:  []entity.Monster{*enemy},
	}

	err := storage.SeedData(context.Background(), &seeder)
	require.NoError(t, err)

	// get available partners
	partners, err := storage.GetAvailablePartners(context.Background())
	require.NoError(t, err)
	require.Contains(t, partners, *partner)

	// get possible enemies
	enemies, err := storage.GetPossibleEnemies(context.Background())
	require.NoError(t, err)
	require.Contains(t, enemies, *enemy)
}

func TestGetPartner(t *testing.T) {
	storage := newStorage(t)

	partner := testutil.NewTestMonster()
	seeder := PokemonSeeder{
		partners: []entity.Monster{*partner},
	}

	err := storage.SeedData(context.Background(), &seeder)
	require.NoError(t, err)

	testCases := []struct {
		Name       string
		PartnerID  string
		ExpPartner *entity.Monster
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
			partner, err := storage.GetPartner(context.Background(), testCase.PartnerID)
			require.NoError(t, err)
			require.Equal(t, partner, testCase.ExpPartner)
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
