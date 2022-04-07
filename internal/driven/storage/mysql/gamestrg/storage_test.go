package gamestrg

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/jmoiron/sqlx"
)

func TestShouldGetGame(t *testing.T) {
	db0, mock, err := sqlmock.New()
	db := sqlx.NewDb(db0, "mysql")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gamestrg := New(db)

	columns := []string{"id", "player_name", "created_at", "battle_won", "scenario",
		"id", "name", "avatar_url", "max_health", "attack", "defense", "speed"}
	gameId := "b1c87c5c-2ac3-471d-9880-4812552ee15d"

	mock.ExpectQuery("^SELECT (.+) FROM games g").
		WithArgs(gameId).
		WillReturnRows(
			sqlmock.NewRows(columns).
				FromCSVString("1,player1,1646205996,1,BATTLE_3,b1c87c5c-2ac3-471d-9880-4812552ee15d,Pikachu,https://example.com/025.png,100,49,49,45"))

	ctx := context.Background()

	if _, err = gamestrg.GetGame(ctx, gameId); err != nil {
		t.Errorf("error was not expected while getting partners: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveGameInsert(t *testing.T) {
	db0, mock, err := sqlmock.New()
	db := sqlx.NewDb(db0, "mysql")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gamestrg := New(db)

	partner := entity.Pokemon{
		ID:        "b1c87c5c-2ac3-471d-9880-4812552ee15d",
		Name:      "Pikachu",
		AvatarURL: "https://example.com/025.png",
		BattleStats: entity.BattleStats{
			Health:    100,
			MaxHealth: 100,
			Attack:    49,
			Defense:   49,
			Speed:     45,
		},
	}

	game := entity.Game{
		ID:         "b1c87c5c-2ac3-471d-9880-4812552ee15d",
		PlayerName: "player1",
		CreatedAt:  1646205996,
		BattleWon:  1,
		Scenario:   "BATTLE_3",
		Partner:    &partner,
	}

	mock.ExpectQuery("^SELECT (.+) FROM games").
		WithArgs(game.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO games").
		WithArgs(game.ID, game.PlayerName, game.CreatedAt, game.BattleWon, game.Scenario, game.Partner.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.Background()

	if err = gamestrg.SaveGame(ctx, game); err != nil {
		t.Errorf("error was not expected while getting partners: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveGameUpdate(t *testing.T) {
	db0, mock, err := sqlmock.New()
	db := sqlx.NewDb(db0, "mysql")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gamestrg := New(db)

	partner := entity.Pokemon{
		ID:        "b1c87c5c-2ac3-471d-9880-4812552ee15d",
		Name:      "Pikachu",
		AvatarURL: "https://example.com/025.png",
		BattleStats: entity.BattleStats{
			Health:    100,
			MaxHealth: 100,
			Attack:    49,
			Defense:   49,
			Speed:     45,
		},
	}

	game := entity.Game{
		ID:         "b1c87c5c-2ac3-471d-9880-4812552ee15d",
		PlayerName: "player1",
		CreatedAt:  1646205996,
		BattleWon:  1,
		Scenario:   "BATTLE_3",
		Partner:    &partner,
	}

	columns := []string{"id"}

	rows := mock.NewRows(columns).
		AddRow(game.ID)

	ctx := context.Background()

	mock.ExpectQuery("^SELECT (.+) FROM games").
		WithArgs(game.ID).
		WillReturnRows(rows)

	mock.ExpectExec("^UPDATE games").
		WithArgs(game.PlayerName, game.CreatedAt, game.BattleWon, game.Scenario, game.Partner.ID, game.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err = gamestrg.SaveGame(ctx, game); err != nil {
		t.Errorf("error was not expected while getting partners: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
