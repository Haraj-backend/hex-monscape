package main

import (
	"fmt"
	"io/ioutil"

	"github.com/Haraj-backend/hex-monscape/internal/core/service/battle"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/play"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jmoiron/sqlx"

	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/ddb"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql"
)

type storageDeps struct {
	BattleGameStorage    battle.GameStorage
	BattleBattleStorage  battle.BattleStorage
	BattleMonsterStorage battle.MonsterStorage
	PlayGameStorage      play.GameStorage
	PlayPartnerStorage   play.PartnerStorage
}

func initStorageDeps(cfg config) (*storageDeps, error) {
	var deps storageDeps

	switch cfg.Storage.Type {
	case storageTypeMemory:
		// initialize monster storage
		monsterData, err := ioutil.ReadFile(cfg.Storage.Memory.MonsterDataPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read monster data due: %w", err)
		}
		monsterStorage, err := memory.NewMonsterStorage(memory.MonsterStorageConfig{MonsterData: monsterData})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize monster storage due: %w", err)
		}

		// initialize game storage
		gameStorage := memory.NewGameStorage()

		// initialize battle storage
		battleStorage := memory.NewBattleStorage()

		// set storages
		deps.BattleGameStorage = gameStorage
		deps.BattleBattleStorage = battleStorage
		deps.BattleMonsterStorage = monsterStorage
		deps.PlayGameStorage = gameStorage
		deps.PlayPartnerStorage = monsterStorage

	case storageTypeDynamoDB:
		// initialize aws session
		awsSess := session.Must(session.NewSessionWithOptions(
			session.Options{
				Config: aws.Config{Endpoint: aws.String(cfg.Storage.DynamoDB.LocalstackEndpoint)},
			},
		))

		// initialize dynamodb client
		dynamoClient := dynamodb.New(awsSess)

		// initialize monster storage
		monsterStorage, err := ddb.NewMonsterStorage(ddb.MonsterStorageConfig{
			DynamoClient: dynamoClient,
			TableName:    cfg.Storage.DynamoDB.MonsterTableName,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize monster storage due: %w", err)
		}

		// initialize game storage
		gameStorage, err := ddb.NewGameStorage(ddb.GameStorageConfig{
			DynamoClient: dynamoClient,
			TableName:    cfg.Storage.DynamoDB.GameTableName,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize game storage due: %w", err)
		}

		// initialize battle storage
		battleStorage, err := ddb.NewBattleStorage(ddb.BattleStorageConfig{
			DynamoClient: dynamoClient,
			TableName:    cfg.Storage.DynamoDB.BattleTableName,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize battle storage due: %w", err)
		}

		// set storages
		deps.BattleGameStorage = gameStorage
		deps.BattleBattleStorage = battleStorage
		deps.BattleMonsterStorage = monsterStorage
		deps.PlayGameStorage = gameStorage
		deps.PlayPartnerStorage = monsterStorage

	case storageTypeMySQL:
		// initialize sql client
		sqlClient, err := sqlx.Open("mysql", cfg.Storage.MySQL.SQLDSN)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize sql client due: %w", err)
		}

		// initialize monster storage
		monsterStorage, err := mysql.NewMonsterStorage(mysql.MonsterStorageConfig{SQLClient: sqlClient})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize monster storage due: %w", err)
		}

		// initialize game storage
		gameStorage, err := mysql.NewGameStorage(mysql.GameStorageConfig{SQLClient: sqlClient})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize game storage due: %w", err)
		}

		// initialize battle storage
		battleStorage, err := mysql.NewBattleStorage(mysql.BattleStorageConfig{SQLClient: sqlClient})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize battle storage due: %w", err)
		}

		// set storages
		deps.BattleGameStorage = gameStorage
		deps.BattleBattleStorage = battleStorage
		deps.BattleMonsterStorage = monsterStorage
		deps.PlayGameStorage = gameStorage
		deps.PlayPartnerStorage = monsterStorage

	default:
		return nil, fmt.Errorf("unknown storage type: %v", cfg.Storage.Type)
	}

	return &deps, nil
}
