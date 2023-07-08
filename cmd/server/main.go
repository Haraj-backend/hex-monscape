package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	membattlestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/battlestrg"
	memgamestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/gamestrg"
	memmonstrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/monstrg"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jmoiron/sqlx"

	ddbbattlestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/battlestrg"
	ddbgamestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/gamestrg"
	ddbmonstrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/monstrg"

	sqlbattlestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/battlestrg"
	sqlgamestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/gamestrg"
	sqlmonstrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/monstrg"

	"github.com/Haraj-backend/hex-monscape/internal/driver/rest"
	"github.com/gosidekick/goconfig"

	"github.com/Haraj-backend/hex-monscape/internal/core/service/battle"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/play"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// initialize configs
	var cfg config
	err := goconfig.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse configs due: %v", err)
	}

	// initialize storages depending on storage type: memory, dynamodb, mysql
	var (
		bGameStorage    battle.GameStorage
		bBattleStorage  battle.BattleStorage
		bMonsterStorage battle.MonsterStorage
		pGameStorage    play.GameStorage
		pPartnerStorage play.PartnerStorage
	)
	switch cfg.Storage.Type {
	case storageTypeMemory:
		// initialize monster storage
		monsterData, err := ioutil.ReadFile(cfg.Storage.Memory.MonsterDataPath)
		if err != nil {
			log.Fatalf("unable to read monster data due: %v", err)
		}
		monsterStorage, err := memmonstrg.New(memmonstrg.Config{MonsterData: monsterData})
		if err != nil {
			log.Fatalf("unable to initialize monster storage due: %v", err)
		}

		// initialize game storage
		gameStorage := memgamestrg.New()

		// initialize battle storage
		battleStorage := membattlestrg.New()

		// set storages
		bGameStorage = gameStorage
		bBattleStorage = battleStorage
		bMonsterStorage = monsterStorage
		pGameStorage = gameStorage
		pPartnerStorage = monsterStorage

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
			log.Fatalf("unable to initialize monster storage due: %v", err)
		}
		// initialize game storage
		gameStorage, err := ddbgamestrg.New(ddbgamestrg.Config{
			DynamoClient: dynamoClient,
			TableName:    cfg.Storage.DynamoDB.GameTableName,
		})
		if err != nil {
			log.Fatalf("unable to initialize game storage due: %v", err)
		}
		// initialize battle storage
		battleStorage, err := ddbbattlestrg.New(ddbbattlestrg.Config{
			DynamoClient: dynamoClient,
			TableName:    cfg.Storage.DynamoDB.BattleTableName,
		})
		if err != nil {
			log.Fatalf("unable to initialize battle storage due: %v", err)
		}

		// set storages
		bGameStorage = gameStorage
		bBattleStorage = battleStorage
		bMonsterStorage = monsterStorage
		pGameStorage = gameStorage
		pPartnerStorage = monsterStorage

	case storageTypeMySQL:
		// initialize sql client
		sqlClient, err := sqlx.Open("mysql", cfg.Storage.MySQL.SQLDSN)
		if err != nil {
			log.Fatalf("unable to initialize sql client due: %v", err)
		}
		// initialize monster storage
		monsterStorage, err := sqlmonstrg.New(sqlmonstrg.Config{SQLClient: sqlClient})
		if err != nil {
			log.Fatalf("unable to initialize monster storage due: %v", err)
		}
		// initialize game storage
		gameStorage, err := sqlgamestrg.New(sqlgamestrg.Config{SQLClient: sqlClient})
		if err != nil {
			log.Fatalf("unable to initialize game storage due: %v", err)
		}
		// initialize battle storage
		battleStorage, err := sqlbattlestrg.New(sqlbattlestrg.Config{SQLClient: sqlClient})
		if err != nil {
			log.Fatalf("unable to initialize battle storage due: %v", err)
		}

		// set storages
		bGameStorage = gameStorage
		bBattleStorage = battleStorage
		bMonsterStorage = monsterStorage
		pGameStorage = gameStorage
		pPartnerStorage = monsterStorage

	default:
		log.Fatalf("unknown storage type: %v", cfg.Storage.Type)
	}

	// initialize play service
	playService, err := play.NewService(play.ServiceConfig{
		GameStorage:    pGameStorage,
		PartnerStorage: pPartnerStorage,
	})
	if err != nil {
		log.Fatalf("unable to initialize play service due: %v", err)
	}

	// initialize battle service
	battleService, err := battle.NewService(battle.ServiceConfig{
		GameStorage:    bGameStorage,
		BattleStorage:  bBattleStorage,
		MonsterStorage: bMonsterStorage,
	})
	if err != nil {
		log.Fatalf("unable to initialize battle service due: %v", err)
	}

	// initialize rest api
	api, err := rest.NewAPI(rest.APIConfig{
		PlayingService: playService,
		BattleService:  battleService,
	})
	if err != nil {
		log.Fatalf("unable to initialize rest api due: %v", err)
	}

	// initialize server
	server := &http.Server{
		Addr:        ":" + cfg.Port,
		Handler:     api.GetHandler(),
		ReadTimeout: 3 * time.Second,
	}

	// run server
	log.Printf("[INFO] server is listening on :%v...", cfg.Port)
	log.Printf("[INFO] please wait a moment until the game client ready in http://localhost:8161 to play the game...")
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("unable to start server due: %v", err)
	}
}
