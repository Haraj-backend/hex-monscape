package mysql

import (
	"database/sql"
	"log"
)

type Database struct {
	db *sql.DB
}

// New creates a new Database
func New(driverName, dataSourceName string) (*Database, error) {
	// connect
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("db connection failure: %v", err)
	}

	// test db connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("db ping failure: %v", err)
	}

	return &Database{db: db}, nil
}

// CloseDbConnection closes the db  connection
func (da Database) CloseDbConnection() {
	err := da.db.Close()
	if err != nil {
		log.Fatalf("db close failure: %v", err)
	}
}
