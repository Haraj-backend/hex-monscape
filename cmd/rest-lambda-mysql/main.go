package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/play"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/battlestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/gamestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/pokestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driver/rest"
	"github.com/apex/gateway"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

const (
	serviceName  = "rest-lambda-mysql"
	addr         = ":9186"
	envKeySQLDSN = "SQL_DSN"
)

func main() {
	isServer := os.Getenv("SERVER_DEPLOYMENT") == "true"

	sqlDSN := os.Getenv(envKeySQLDSN)
	sqlClient, err := sqlx.Connect("mysql", sqlDSN)
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
		ServiceName:    serviceName,
	})
	if err != nil {
		log.Fatalf("unable to init rest service: %v", err)
	}

	// run server
	err = listenAndServe(isServer, addr, api.GetHandler())
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
