package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/play"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/battlestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/gamestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/pokestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driver/rest"
	"github.com/apex/gateway"
	"github.com/gosidekick/goconfig"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	ServiceName              string `cfg:"service_name" cfgRequired:"true" cfgDefault:"rest-lambda-mysql"`
	Port                     string `cfg:"port" cfgDefault:"9186"`
	SQLDSN                   string `cfg:"sql_dsn" cfgRequired:"true"`
	OtelExporterOTLPEndpoint string `cfg:"otel_exporter_otlp_endpoint" cfgRequired:"true"`
	IsServer                 bool   `cfg:"server_deployment" cfgDefault:"false"`
}

func main() {
	var cfg config
	err := goconfig.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse configs due: %v", err)
	}

	sqlClient, err := sqlx.Connect("mysql", cfg.SQLDSN)
	if err != nil {
		log.Fatalf("unable to init db connection: %v", err)
	}
	defer sqlClient.Close()

	// init pokemon storage
	configPokeDB := pokestrg.Config{SQLClient: sqlClient}
	pokeStrg, err := pokestrg.New(configPokeDB)
	if err != nil {
		log.Fatalf("unable to initialize pokemon storage due: %v", err)
	}

	// init game storage
	configGameDB := gamestrg.Config{SQLClient: sqlClient}
	gameStrg, err := gamestrg.New(configGameDB)
	if err != nil {
		log.Fatalf("unable to initialize game storage due: %v", err)
	}

	// init battle storage
	configBattleDB := battlestrg.Config{SQLClient: sqlClient}
	battleStrg, err := battlestrg.New(configBattleDB)
	if err != nil {
		log.Fatalf("unable to initialize battle storage due: %v", err)
	}

	// init play service
	playService, err := play.NewService(play.ServiceConfig{
		PartnerStorage: pokeStrg,
		GameStorage:    gameStrg,
	})
	if err != nil {
		log.Fatalf("unable to init play service: %v", err)
	}

	// init battle service
	battleService, err := battle.NewService(battle.ServiceConfig{
		GameStorage:    gameStrg,
		BattleStorage:  battleStrg,
		PokemonStorage: pokeStrg,
	})
	if err != nil {
		log.Fatalf("unable to init battle service: %v", err)
	}

	// init rest service
	api, err := rest.NewAPI(rest.APIConfig{
		PlayingService: playService,
		BattleService:  battleService,
		ServiceName:    cfg.ServiceName,
	})
	if err != nil {
		log.Fatalf("unable to init rest service: %v", err)
	}

	// run server
	err = listenAndServe(cfg.IsServer, ":"+cfg.Port, api.GetHandler())
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
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
