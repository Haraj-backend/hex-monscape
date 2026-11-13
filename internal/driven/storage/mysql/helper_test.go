package mysql

import (
	"fmt"
	"os"

	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

const envKeySQLDSN = "TEST_SQL_DSN"

func newTestSQLClient() (*sqlx.DB, error) {
	sqlClient, err := sqlx.Connect("mysql", os.Getenv(envKeySQLDSN))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize sql client due: %w", err)
	}

	return sqlClient, nil
}

func newTestMonsterRow(isPartnerable bool) monsterRow {
	row := newMonsterRow(testutil.NewTestMonster())
	if isPartnerable {
		row.IsPartnerable = 1
	}

	return *row
}

func insertMonster(sqlClient *sqlx.DB, row monsterRow) error {
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
	_, err := sqlClient.NamedExec(query, row)
	if err != nil {
		return fmt.Errorf("unable to insert monster due: %w", err)
	}

	return nil
}
