package shared

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type RowResultInterface interface {
	Scan(dest ...interface{}) error
}

type Database struct {
	Db *sqlx.DB
}

// New creates a new Database
func New(dataSourceName string) (*Database, error) {
	// connect
	db, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("db connection failure: %v", err)
	}

	// test db connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db ping failure: %v", err)
	}

	return &Database{Db: db}, nil
}

// CloseDbConnection closes the db  connection
func (da Database) CloseDbConnection() {
	err := da.Db.Close()
	if err != nil {
		log.Fatalf("db close failure: %v", err)
	}
}
