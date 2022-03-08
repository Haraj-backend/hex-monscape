package pokestrg_test

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/config"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/pokestrg"
)

var pikachu = entity.Pokemon{
	ID:   "pikachu-id",
	Name: "Pikachu",
	BattleStats: entity.BattleStats{
		MaxHealth: 100,
		Attack:    25,
		Defense:   10,
		Speed:     20,
	},
	AvatarURL: "http://avatar.com",
}

var bulbasaur = entity.Pokemon{
	ID:   "bulbasaur-id",
	Name: "Bulbasaur",
	BattleStats: entity.BattleStats{
		MaxHealth: 100,
		Attack:    27,
		Defense:   20,
		Speed:     10,
	},
	AvatarURL: "http://avatar.com",
}

var pokeStrgConfig = pokestrg.Config{
	DynamoDBStorageConfig: config.DynamoDBStorageConfig{
		DB: config.BuildDynamoDBInstance(),
	},
	Partners: []entity.Pokemon{pikachu},
	Enemies:  []entity.Pokemon{bulbasaur},
}

func TestGetPartner(t *testing.T) {
	storage, err := pokestrg.NewDynamoDBStorage(pokeStrgConfig)
	if err != nil {
		t.Fatalf("error creating storage instance %s", err)
	}

	result, err := storage.GetPartner(context.Background(), pikachu.ID)
	if err != nil {
		t.Errorf("function return error: %s", err)
	}

	if result.ID != pikachu.ID {
		t.Error("fetched object does not return same ID")
	}

	if result.BattleStats.Attack != pikachu.BattleStats.Attack {
		t.Error("fetched object does not return same BattleStats")
	}
}

func TestGetPossibleEnemies(t *testing.T) {
	storage, err := pokestrg.NewDynamoDBStorage(pokeStrgConfig)
	if err != nil {
		t.Fatalf("error creating storage instance %s", err)
	}

	enemies, err := storage.GetPossibleEnemies(context.Background())
	if err != nil {
		t.Errorf("function return error: %s", err)
	}

	if len(enemies) != 1 {
		t.Error("fetched enemies count does not match test data")
	}

	if enemies[0].ID != bulbasaur.ID {
		t.Errorf("fetched enemy ID does not match test data ")
	}
}

func TestGetAvailableEnemies(t *testing.T) {
	storage, err := pokestrg.NewDynamoDBStorage(pokeStrgConfig)
	if err != nil {
		t.Fatalf("error creating storage instance %s", err)
	}

	partners, err := storage.GetAvailablePartners(context.Background())
	if err != nil {
		t.Errorf("function return error: %s", err)
	}

	if len(partners) != 1 {
		t.Error("fetched partners count does not match test data")
	}

	if partners[0].ID != pikachu.ID {
		t.Errorf("fetched partner ID does not match test data ")
	}
}
