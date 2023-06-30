package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/play"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/battlestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/gamestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/dynamodb/pokestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driver/rest"
	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gosidekick/goconfig"
	_ "github.com/gosidekick/goconfig/yaml"
)

const (
	serviceName = "rest-lambda-ddb"
)

type config struct {
	LocalDeployment localDeploymentConfig `yaml:"local_deployment" cfg:"local_deployment"`
	Dynamo          dynamoConfig          `yaml:"ddb" cfg:"ddb"`
}

type localDeploymentConfig struct {
	Enabled  bool   `cfg:"enabled"`
	Endpoint string `cfg:"endpoint"`
	Port     int    `cfg:"port" cfgDefault:"9186"`
}

type dynamoConfig struct {
	BattleTable  string `cfg:"battle_table" cfgDefault:"Battles"`
	GameTable    string `cfg:"game_table" cfgDefault:"Games"`
	PokemonTable string `cfg:"pokemon_table" cfgDefault:"Pokemons"`
}

func main() {
	// read config
	var cfg config
	err := goconfig.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse config due: %v", err)
	}

	// initialize aws session
	awsSess := session.Must(session.NewSession())
	if cfg.LocalDeployment.Enabled {
		awsSess.Config.Endpoint = &cfg.LocalDeployment.Endpoint
	}

	// initialize dynamodb client
	ddbClient := dynamodb.New(awsSess)

	// initialize battle storage
	battleStrg, err := battlestrg.New(battlestrg.Config{
		DynamoClient: ddbClient,
		TableName:    cfg.Dynamo.BattleTable,
	})
	if err != nil {
		log.Fatalf("unable to initialize battleStrg due: %v", err)
	}

	// initialize game storage
	gameStrg, err := gamestrg.New(gamestrg.Config{
		DynamoClient: ddbClient,
		TableName:    cfg.Dynamo.GameTable,
	})
	if err != nil {
		log.Fatalf("unable to initialize gameStrg due: %v", err)
	}

	// initialize pokemon storage
	pokeStrg, err := pokestrg.New(pokestrg.Config{
		DynamoClient: ddbClient,
		TableName:    cfg.Dynamo.PokemonTable,
	})
	if err != nil {
		log.Fatalf("unable to initialize pokeStrg due: %v", err)
	}

	// initialize play service
	playSvc, err := play.NewService(play.ServiceConfig{
		GameStorage:    gameStrg,
		PartnerStorage: pokeStrg,
	})
	if err != nil {
		log.Fatalf("unable to initialize play service due: %v", err)
	}

	// initialize battle service
	battleSvc, err := battle.NewService(battle.ServiceConfig{
		BattleStorage:  battleStrg,
		GameStorage:    gameStrg,
		PokemonStorage: pokeStrg,
	})
	if err != nil {
		log.Fatalf("unable to initialize battle service due: %v", err)
	}

	// initialize rest API
	api, err := rest.NewAPI(rest.APIConfig{
		PlayingService: playSvc,
		BattleService:  battleSvc,
		ServiceName:    serviceName,
	})
	if err != nil {
		log.Fatalf("unable to initialize rest api due: %v", err)
	}

	err = listenAndServe(cfg.LocalDeployment.Enabled, fmt.Sprintf(":%d", cfg.LocalDeployment.Port), api.GetHandler())
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("unable to start server due: %v", err)
	}
}

func listenAndServe(serverMode bool, addr string, handler http.Handler) error {
	if serverMode {
		log.Printf("Running in server mode at %s", addr)
		return http.ListenAndServe(addr, handler)
	}

	log.Println("Running in serverless mode")
	return gateway.ListenAndServe("", handler)
}
