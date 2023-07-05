package shared

import (
	"fmt"
	"os"

	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/jmoiron/sqlx"
)

const envKeySQLDSN = "TEST_SQL_DSN"

func NewTestSQLClient() (*sqlx.DB, error) {
	sqlDSN := os.Getenv(envKeySQLDSN)
	sqlClient, err := sqlx.Connect("mysql", sqlDSN)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize sql client due: %w", err)
	}
	return sqlClient, nil
}

func NewTestMonsterRow(isPartnerable bool) MonsterRow {
	row := ToMonsterRow(testutil.NewTestMonster())
	if isPartnerable {
		row.IsPartnerable = 1
	}
	return *row
}

func InsertMonster(sqlClient *sqlx.DB, monsterRow MonsterRow) error {
	query := `
		INSERT INTO monster (
			id,
			name,
			health,
			max_health,
			attack,
			defense,
			speed,
			avatar_url,
			is_partnerable
		) VALUES (
			:id,
			:name,
			:health,
			:max_health,
			:attack,
			:defense,
			:speed,
			:avatar_url,
			:is_partnerable
		)
	`
	_, err := sqlClient.NamedExec(query, monsterRow)
	if err != nil {
		return fmt.Errorf("unable to insert monster due: %w", err)
	}
	return nil
}
