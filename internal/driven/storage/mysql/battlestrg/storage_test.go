package battlestrg

import (
	"context"
	"fmt"
	"testing"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	_ "github.com/go-sql-driver/mysql"
)

func TestSaveBattle(t *testing.T) {
	// initialize sql client
	sqlClient, err := shared.NewTestSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// create partner pokemon
	partner := shared.NewTestPokemon()
	err = shared.InsertTestPokemon(sqlClient, partner, 1)
	require.NoError(t, err)
	// create enemy pokemon
	enemy := shared.NewTestPokemon()
	err = shared.InsertTestPokemon(sqlClient, enemy, 0)
	require.NoError(t, err)
	// save battle
	b := newBattle(partner, enemy)
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
	sqlClient, err := shared.NewTestSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// create partner pokemon
	partner := shared.NewTestPokemon()
	err = shared.InsertTestPokemon(sqlClient, partner, 1)
	require.NoError(t, err)
	// create enemy pokemon
	enemy := shared.NewTestPokemon()
	err = shared.InsertTestPokemon(sqlClient, enemy, 0)
	require.NoError(t, err)
	// save battle
	b := newBattle(partner, enemy)
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
	sqlClient, err := shared.NewTestSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// create partner pokemon
	partner := shared.NewTestPokemon()
	err = shared.InsertTestPokemon(sqlClient, partner, 1)
	require.NoError(t, err)
	// create enemy pokemon
	enemy := shared.NewTestPokemon()
	err = shared.InsertTestPokemon(sqlClient, enemy, 0)
	require.NoError(t, err)
	// save battle
	b := newBattle(partner, enemy)
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
	sqlClient, err := shared.NewTestSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// check whether battle exists on database
	savedBattle, err := strg.GetBattle(context.Background(), uuid.NewString())
	require.NoError(t, err)
	require.Nil(t, savedBattle)
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

func newBattle(partner entity.Pokemon, enemy entity.Pokemon) battle.Battle {
	return battle.Battle{
		GameID:  uuid.NewString(),
		State:   battle.DECIDE_TURN,
		Partner: &partner,
		Enemy:   &enemy,
		LastDamage: battle.LastDamage{
			Partner: 0,
			Enemy:   10,
		},
	}
}
