package battlestrg_test

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/battlestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/config"
)

var testBattle = battle.Battle{
	GameID: "something-id",
	State:  battle.DECIDE_TURN,
	Partner: &entity.Pokemon{
		ID:   "partner-id",
		Name: "my-partner",
	},
	Enemy: &entity.Pokemon{
		ID:   "enemy-id",
		Name: "my-enemy",
	},
	LastDamage: battle.LastDamage{
		Partner: 10,
		Enemy:   10,
	},
}

var dynamoDBConfig = config.DynamoDBStorageConfig{
	DB: config.BuildDynamoDBInstance(),
}

func TestSaveBattle(t *testing.T) {
	storage, err := battlestrg.NewDynamoDBStorage(dynamoDBConfig)
	if err != nil {
		t.Fatalf("error creating storage instance: %s", err)
	}

	err = storage.SaveBattle(context.Background(), testBattle)
	if err != nil {
		t.Errorf("function return error: %s", err)
	}
}

func TestGetBattle(t *testing.T) {
	storage, err := battlestrg.NewDynamoDBStorage(dynamoDBConfig)
	if err != nil {
		t.Fatalf("error creating storage instance: %s", err)
	}

	result, err := storage.GetBattle(context.Background(), testBattle.GameID)
	if err != nil {
		t.Errorf("function return error: %s", err)
	}

	if result.GameID != testBattle.GameID {
		t.Error("fetched object does not return same GameID")
	}

	if result.Partner.BattleStats.Attack != testBattle.Partner.BattleStats.Attack {
		t.Error("fetched object does not return same Partner")
	}
}
