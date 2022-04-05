package battlestrg

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetBattle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	battlestrg := New(db)

	columns := []string{"game_id", "partner_last_damage", "enemy_last_damage", "state",
		"partner_id", "partner_name", "partner_max_health", "partner_health",
		"partner_attack", "partner_defense", "partner_speed", "partner_avatar_url",
		"enemy_id", "enemy_name", "enemy_max_health", "enemy_health",
		"enemy_attack", "enemy_defense", "enemy_speed", "enemy_avatar_url",
	}
	gameId := "b1c87c5c-2ac3-471d-9880-4812552ee15d"

	mock.ExpectQuery(`
	^SELECT (.+)
	FROM battles b
	LEFT JOIN pokemon_battle_states p (.+)
	LEFT JOIN pokemon_battle_states e (.+)`).
		WithArgs(gameId).
		WillReturnRows(
			sqlmock.NewRows(columns).
				FromCSVString("1,1,1,BATTLE_3,b1c87c5c-2ac3-471d-9880-4812552ee15d,Pikachu,100,100,49,49,45,https://example.com/025.png,1,Bulbasaur,100,100,49,49,45,https://example.com/001.png"))

	ctx := context.Background()

	if _, err = battlestrg.GetBattle(ctx, gameId); err != nil {
		t.Errorf("error was not expected while getting partners: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetBattleNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	battlestrg := New(db)

	gameId := "b1c87c5c-2ac3-471d-9880-4812552ee15d"

	mock.ExpectQuery(`
	^SELECT (.+)
	FROM battles b
	LEFT JOIN pokemon_battle_states p (.+)
	LEFT JOIN pokemon_battle_states e (.+)`).
		WithArgs(gameId).
		WillReturnRows(sqlmock.NewRows([]string{})).
		WillReturnError(nil)

	ctx := context.Background()

	if _, err = battlestrg.GetBattle(ctx, gameId); err != nil {
		t.Errorf("error was not expected while getting partners: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
