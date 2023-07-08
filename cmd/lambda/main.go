package main

import (
	"log"
	"net/http"

	"github.com/Haraj-backend/hex-monscape/internal/core/service/battle"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/play"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/battlestrg"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/gamestrg"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/monstrg"
	"github.com/Haraj-backend/hex-monscape/internal/driver/rest"
	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gosidekick/goconfig"
)

func main() {
	// load configs
	var cfg config
	err := goconfig.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse configs due: %v", err)
	}

	// initialize aws session
	awsSess := session.Must(session.NewSession())
	if cfg.IsLocalDeployment() {
		awsSess.Config.Endpoint = &cfg.LocalstackEndpoint
	}

	// initialize dynamodb client
	dynamoClient := dynamodb.New(awsSess)

	// initialize battle storage
	battleStrg, err := battlestrg.New(battlestrg.Config{
		DynamoClient: dynamoClient,
		TableName:    cfg.DynamoDB.BattleTable,
	})
	if err != nil {
		log.Fatalf("unable to initialize battle storage due: %v", err)
	}

	// initialize game storage
	gameStrg, err := gamestrg.New(gamestrg.Config{
		DynamoClient: dynamoClient,
		TableName:    cfg.DynamoDB.GameTable,
	})
	if err != nil {
		log.Fatalf("unable to initialize game storage due: %v", err)
	}

	// initialize monster storage
	monStrg, err := monstrg.New(monstrg.Config{
		DynamoClient: dynamoClient,
		TableName:    cfg.DynamoDB.MonsterTable,
	})
	if err != nil {
		log.Fatalf("unable to initialize monster storage due: %v", err)
	}

	// initialize play service
	playSvc, err := play.NewService(play.ServiceConfig{
		GameStorage:    gameStrg,
		PartnerStorage: monStrg,
	})
	if err != nil {
		log.Fatalf("unable to initialize play service due: %v", err)
	}

	// initialize battle service
	battleSvc, err := battle.NewService(battle.ServiceConfig{
		GameStorage:    gameStrg,
		BattleStorage:  battleStrg,
		MonsterStorage: monStrg,
	})
	if err != nil {
		log.Fatalf("unable to initialize battle service due: %v", err)
	}

	// initialize rest api
	api, err := rest.NewAPI(rest.APIConfig{
		PlayingService: playSvc,
		BattleService:  battleSvc,
		IsWebEnabled:   true,
	})
	if err != nil {
		log.Fatalf("unable to initialize rest api due: %v", err)
	}

	// listen and serve
	err = listenAndServe(cfg.IsLocalDeployment(), cfg.GetLocalAddr(), api.GetHandler())
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("unable to listen and serve due: %v", err)
	}
}

func listenAndServe(serverMode bool, serverAddr string, handler http.Handler) error {
	if serverMode {
		log.Printf("running in server mode at %v", serverAddr)
		return http.ListenAndServe(serverAddr, handler)
	}

	log.Println("running in serverless mode")
	return gateway.ListenAndServe("", handler)
}

type config struct {
	LocalstackEndpoint string       `cfg:"localstack_endpoint"`
	DynamoDB           dynamoConfig `cfg:"dynamodb"`
}

func (c config) IsLocalDeployment() bool {
	return len(c.LocalstackEndpoint) > 0
}

func (c config) GetLocalAddr() string {
	return ":9186"
}

type dynamoConfig struct {
	BattleTable  string `cfg:"battle_table" cfgRequired:"true"`
	GameTable    string `cfg:"game_table" cfgRequired:"true"`
	MonsterTable string `cfg:"monster_table" cfgRequired:"true"`
}
