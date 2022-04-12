package shared

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

const envKeySQLDSN = "SQL_DSN"

func NewSQLClient() (*sqlx.DB, error) {
	sqlDSN := os.Getenv(envKeySQLDSN)
	sqlClient, err := sqlx.Connect("mysql", sqlDSN)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize sql client due: %w", err)
	}
	return sqlClient, nil
}
