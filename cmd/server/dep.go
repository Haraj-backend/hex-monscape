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

	membattlestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/battlestrg"
	memgamestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/gamestrg"
	memmonstrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/monstrg"

	ddbbattlestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/battlestrg"
	ddbgamestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/gamestrg"
	ddbmonstrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/monstrg"

	sqlbattlestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/battlestrg"
	sqlgamestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/gamestrg"
	sqlmonstrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/monstrg"
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
			return nil, fmt.Errorf("unable to read monster data due: %v", err)
		}
		monsterStorage, err := memmonstrg.New(memmonstrg.Config{MonsterData: monsterData})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize monster storage due: %v", err)
		}

		// initialize game storage
		gameStorage := memgamestrg.New()

		// initialize battle storage
		battleStorage := membattlestrg.New()

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
		monsterStorage, err := ddbmonstrg.New(ddbmonstrg.Config{
			DynamoClient: dynamoClient,
			TableName:    cfg.Storage.DynamoDB.MonsterTableName,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize monster storage due: %v", err)
		}
		// initialize game storage
		gameStorage, err := ddbgamestrg.New(ddbgamestrg.Config{
			DynamoClient: dynamoClient,
			TableName:    cfg.Storage.DynamoDB.GameTableName,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize game storage due: %v", err)
		}
		// initialize battle storage
		battleStorage, err := ddbbattlestrg.New(ddbbattlestrg.Config{
			DynamoClient: dynamoClient,
			TableName:    cfg.Storage.DynamoDB.BattleTableName,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize battle storage due: %v", err)
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
			return nil, fmt.Errorf("unable to initialize sql client due: %v", err)
		}
		// initialize monster storage
		monsterStorage, err := sqlmonstrg.New(sqlmonstrg.Config{SQLClient: sqlClient})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize monster storage due: %v", err)
		}
		// initialize game storage
		gameStorage, err := sqlgamestrg.New(sqlgamestrg.Config{SQLClient: sqlClient})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize game storage due: %v", err)
		}
		// initialize battle storage
		battleStorage, err := sqlbattlestrg.New(sqlbattlestrg.Config{SQLClient: sqlClient})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize battle storage due: %v", err)
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
