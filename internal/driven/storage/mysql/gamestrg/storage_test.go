package gamestrg_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/gamestrg"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/shared"
	"github.com/stretchr/testify/require"

	_ "github.com/go-sql-driver/mysql"
)

func TestSaveGameGetGame(t *testing.T) {
	// initialize storage
	strg := newStorage(t)

	// save Game
	g := newGame(t)
	err := strg.SaveGame(context.Background(), g)
	require.NoError(t, err)

	// check whether Game exists on database
	savedGame, err := strg.GetGame(context.Background(), g.ID)
	require.NoError(t, err)

	// check whether Game data is match
	require.Equal(t, g, *savedGame)
}

func TestGetGameNotFound(t *testing.T) {
	// initialize storage
	strg := newStorage(t)

	// game shouldn't exists in database
	savedGame, err := strg.GetGame(context.Background(), uuid.NewString())
	require.NoError(t, err)
	require.Nil(t, savedGame)
}

func newGame(t *testing.T) entity.Game {
	// ini sql client
	sqlClient, err := shared.NewTestSQLClient()
	require.NoError(t, err)

	// insert partner to database
	partnerRow := shared.NewTestMonsterRow(true)
	err = shared.InsertMonster(sqlClient, partnerRow)
	require.NoError(t, err)

	nowTs := time.Now().Unix()
	return entity.Game{
		ID:         uuid.NewString(),
		PlayerName: fmt.Sprintf("player_%v", nowTs),
		CreatedAt:  nowTs,
		BattleWon:  2,
		Scenario:   entity.ScenarioBattle3,
		Partner:    partnerRow.ToMonster(),
	}
}

func newStorage(t *testing.T) *gamestrg.Storage {
	// initialize sql client
	sqlClient, err := shared.NewTestSQLClient()
	require.NoError(t, err)

	// initialize storage
	strg, err := gamestrg.New(gamestrg.Config{SQLClient: sqlClient})
	require.NoError(t, err)

	return strg
}
