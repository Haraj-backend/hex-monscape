package pokestrg

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestShouldGetPartners(t *testing.T) {
	db0, mock, err := sqlmock.New()
	db := sqlx.NewDb(db0, "mysql")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pokestrg := New(db)

	columns := []string{"id", "name", "health", "max_health", "attack", "defense", "speed", "avatar_url"}

	mock.ExpectQuery("^SELECT (.+) FROM pokemons").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("b1c87c5c-2ac3-471d-9880-4812552ee15d,Pikachu,100,100,49,49,45,https://example.com/025.png"))

	ctx := context.Background()

	if _, err = pokestrg.GetAvailablePartners(ctx); err != nil {
		t.Errorf("error was not expected while getting partners: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldGetEnemies(t *testing.T) {
	db0, mock, err := sqlmock.New()
	db := sqlx.NewDb(db0, "mysql")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pokestrg := New(db)

	columns := []string{"id", "name", "health", "max_health", "attack", "defense", "speed", "avatar_url"}

	mock.ExpectQuery("^SELECT (.+) FROM pokemons").
		WithArgs(0).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("b1c87c5c-2ac3-471d-9880-4812552ee15d,Pikachu,100,100,49,49,45,https://example.com/025.png"))

	ctx := context.Background()

	if _, err = pokestrg.GetPossibleEnemies(ctx); err != nil {
		t.Errorf("error was not expected while getting partners: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldGetPartner(t *testing.T) {
	db0, mock, err := sqlmock.New()
	db := sqlx.NewDb(db0, "mysql")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pokestrg := New(db)

	columns := []string{"id", "name", "health", "max_health", "attack", "defense", "speed", "avatar_url"}
	id := "b1c87c5c-2ac3-471d-9880-4812552ee15d"

	mock.ExpectQuery("^SELECT (.+) FROM pokemons").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("b1c87c5c-2ac3-471d-9880-4812552ee15d,Pikachu,100,100,49,49,45,https://example.com/025.png"))

	ctx := context.Background()

	if _, err = pokestrg.GetPartner(ctx, id); err != nil {
		t.Errorf("error was not expected while getting partners: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldGetNoPartner(t *testing.T) {
	db0, mock, err := sqlmock.New()
	db := sqlx.NewDb(db0, "mysql")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pokestrg := New(db)

	id := "b1c87c5c-2ac3-471d-9880-4812552ee15d"

	mock.ExpectQuery("^SELECT (.+) FROM pokemons").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()

	if _, err = pokestrg.GetPartner(ctx, id); err == nil {
		t.Errorf("error was not expected while getting partners: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
