package battlestrg

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

const envKeySQLDSN = "SQL_DSN"

func TestSaveBattle(t *testing.T) {
	// initialize sql client
	sqlClient, err := newSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// save battle
	b := battle.Battle{
		GameID: uuid.NewString(),
		State:  battle.DECIDE_TURN,
		Partner: &entity.Pokemon{
			ID:   uuid.NewString(),
			Name: "Pokemon 1",
			BattleStats: entity.BattleStats{
				Health:    10,
				MaxHealth: 10,
				Attack:    10,
				Defense:   10,
				Speed:     10,
			},
			AvatarURL: "https://someurl.com/1.jpg",
		},
		Enemy: &entity.Pokemon{
			ID:   uuid.NewString(),
			Name: "Pokemon 2",
			BattleStats: entity.BattleStats{
				Health:    20,
				MaxHealth: 20,
				Attack:    20,
				Defense:   20,
				Speed:     20,
			},
			AvatarURL: "https://someurl.com/2.jpg",
		},
		LastDamage: battle.LastDamage{
			Partner: 0,
			Enemy:   10,
		},
	}
	err = strg.SaveBattle(context.Background(), b)
	require.NoError(t, err)
	// check whether battle exists on database
	savedBattle, err := getBattle(sqlClient, b.GameID)
	require.NoError(t, err)
	// check whether battle data is match
	require.Equal(t, b, savedBattle)
}

func TestSaveBattleExistingBattle(t *testing.T) {
	// TODO
}

func TestGetBattle(t *testing.T) {
	// TODO
}

func TestGetBattleNotFound(t *testing.T) {
	// TODO
}

func newSQLClient() (*sqlx.DB, error) {
	sqlClient, err := sqlx.Connect("mysql", os.Getenv(envKeySQLDSN))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize sql client due: %w", err)
	}
	return sqlClient, nil
}

func getBattle(sqlClient *sqlx.DB, gameID string) (*battle.Battle, error) {
	var row battleRow
	query := `SELECT * FROM battles WHERE game_id = ?`
	err := sqlClient.Get(&row, query, gameID)
	if err != nil {
		return nil, fmt.Errorf("unable to execute query due: %w", err)
	}
	return row.ToBattle(), nil
}
