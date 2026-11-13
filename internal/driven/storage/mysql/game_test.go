package mysql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/stretchr/testify/require"
)

func TestSaveGameGetGame(t *testing.T) {
	// initialize storage
	strg := newGameStorage(t)

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
	strg := newGameStorage(t)

	// game shouldn't exists in database
	savedGame, err := strg.GetGame(context.Background(), uuid.NewString())
	require.NoError(t, err)
	require.Nil(t, savedGame)
}

func newGame(t *testing.T) entity.Game {
	// ini sql client
	sqlClient, err := newTestSQLClient()
	require.NoError(t, err)

	// insert partner to database
	partnerRow := newTestMonsterRow(true)
	err = insertMonster(sqlClient, partnerRow)
	require.NoError(t, err)

	nowTs := time.Now().Unix()
	return entity.Game{
		ID:         uuid.NewString(),
		PlayerName: fmt.Sprintf("player_%v", nowTs),
		CreatedAt:  nowTs,
		BattleWon:  2,
		Scenario:   entity.ScenarioBattle3,
		Partner:    partnerRow.toMonster(),
	}
}

func newGameStorage(t *testing.T) *GameStorage {
	// initialize sql client
	sqlClient, err := newTestSQLClient()
	require.NoError(t, err)

	// initialize storage
	strg, err := NewGameStorage(GameStorageConfig{SQLClient: sqlClient})
	require.NoError(t, err)

	return strg
}
