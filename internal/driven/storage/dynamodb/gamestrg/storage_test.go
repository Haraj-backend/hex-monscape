package gamestrg_test

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/config"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/gamestrg"
)

var (
	testGame = entity.Game{
		ID:         "game-id",
		PlayerName: "Alfat",
		Partner: &entity.Pokemon{
			ID: "partner-id",
		},
		BattleWon: 2,
		Scenario:  entity.BATTLE_3,
	}

	dynamoDBConfig = config.DynamoDBStorageConfig{
		DB: config.BuildDynamoDBInstance(),
	}
)

func TestSaveGame(t *testing.T) {
	storage, err := gamestrg.NewDynamoDBStorage(dynamoDBConfig)
	if err != nil {
		t.Fatalf("error creating storage instance %s", err)
	}

	err = storage.SaveGame(context.Background(), testGame)
	if err != nil {
		t.Errorf("function return error: %s", err)
	}
}

func TestGetGame(t *testing.T) {
	storage, err := gamestrg.NewDynamoDBStorage(dynamoDBConfig)
	if err != nil {
		t.Fatalf("error creating storage instance %s", err)
	}

	result, err := storage.GetGame(context.Background(), testGame.ID)
	if err != nil {
		t.Errorf("function return error: %s", err)
	}

	if result.ID != testGame.ID {
		t.Errorf("fetched object does not return same ID")
	}

	if result.Partner.ID != testGame.Partner.ID {
		t.Error("fetched object does not match test data Partner's ID")
	}

	if result.Scenario != testGame.Scenario {
		t.Error("fetched object does not match test data Scenario")
	}

}
