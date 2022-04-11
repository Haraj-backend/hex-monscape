package battlestrg

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	_ "github.com/go-sql-driver/mysql"
)

const envKeySQLDSN = "SQL_DSN"

func TestSaveBattle(t *testing.T) {
	// initialize sql client
	sqlClient, err := newSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(shared.Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// save battle
	b := newBattle()
	err = strg.SaveBattle(context.Background(), b)
	require.NoError(t, err)
	// check whether battle exists on database
	savedBattle, err := getBattle(sqlClient, b.GameID)
	require.NoError(t, err)
	// check whether battle data is match
	require.Equal(t, b, *savedBattle)
}

func TestSaveBattleExistingBattle(t *testing.T) {
	// initialize sql client
	sqlClient, err := newSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(shared.Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// save battle
	b := newBattle()
	err = strg.SaveBattle(context.Background(), b)
	require.NoError(t, err)
	// override battle state
	b.State = battle.ENEMY_TURN
	// save again
	err = strg.SaveBattle(context.Background(), b)
	require.NoError(t, err)
	// check whether battle exists on database
	savedBattle, err := getBattle(sqlClient, b.GameID)
	require.NoError(t, err)
	// check whether battle data is match
	require.Equal(t, b, *savedBattle)
}

func TestGetBattle(t *testing.T) {
	// initialize sql client
	sqlClient, err := newSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(shared.Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// save battle
	b := newBattle()
	err = strg.SaveBattle(context.Background(), b)
	require.NoError(t, err)
	// override battle state
	b.State = battle.ENEMY_TURN
	// save again
	err = strg.SaveBattle(context.Background(), b)
	require.NoError(t, err)
	// check whether battle exists on database
	savedBattle, err := strg.GetBattle(context.Background(), b.GameID)
	require.NoError(t, err)
	// check whether battle data is match
	require.Equal(t, b, *savedBattle)
}

func TestGetBattleNotFound(t *testing.T) {
	// initialize sql client
	sqlClient, err := newSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(shared.Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// check whether battle exists on database
	savedBattle, err := strg.GetBattle(context.Background(), uuid.NewString())
	require.NoError(t, err)
	require.Nil(t, savedBattle)
}

func newSQLClient() (*sqlx.DB, error) {
	sqlDSN := os.Getenv(envKeySQLDSN)
	sqlClient, err := sqlx.Connect("mysql", sqlDSN)
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

func newBattle() battle.Battle {
	return battle.Battle{
		GameID: uuid.NewString(),
		State:  battle.DECIDE_TURN,
		Partner: &entity.Pokemon{
			ID:   "b1c87c5c-2ac3-471d-9880-4812552ee15d",
			Name: "Pikachu",
			BattleStats: entity.BattleStats{
				Health:    100,
				MaxHealth: 100,
				Attack:    49,
				Defense:   49,
				Speed:     45,
			},
			AvatarURL: "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png",
		},
		Enemy: &entity.Pokemon{
			ID:   "0f9b84b6-a768-4ba9-8800-207740fc993d",
			Name: "Bulbasaur",
			BattleStats: entity.BattleStats{
				Health:    100,
				MaxHealth: 100,
				Attack:    49,
				Defense:   49,
				Speed:     45,
			},
			AvatarURL: "https://assets.pokemon.com/assets/cms2/img/pokedex/full/001.png",
		},
		LastDamage: battle.LastDamage{
			Partner: 0,
			Enemy:   10,
		},
	}
}
