package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/play"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/battlestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/gamestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/pokestrg"
	mysql "github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
	"github.com/Haraj-backend/hex-pokebattle/internal/driver/rest"

	_ "github.com/go-sql-driver/mysql"
)

const addr = ":9186"

func main() {
	log.Printf("Running service...")

	Db, err := mysql.NewSQLClient()
	if err != nil {
		log.Fatalf("unable to init db connection: %v", err)
	}

	// init pokemon storage
	configPokeDB := pokestrg.Config{SQLClient: Db}
	pokeStrg, err := pokestrg.New(configPokeDB)
	if err != nil {
		log.Fatalf("unable to initialize pokemon storage due: %v", err)
	}

	// init game storage
	configGameDB := gamestrg.Config{SQLClient: Db}
	gameStrg, err := gamestrg.New(configGameDB)
	if err != nil {
		log.Fatalf("unable to initialize game storage due: %v", err)
	}

	// init battle storage
	configBattleDB := battlestrg.Config{SQLClient: Db}
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
	})
	if err != nil {
		log.Fatalf("unable to init rest service: %v", err)
	}

	// initialize server
	server := &http.Server{
		Addr:        addr,
		Handler:     api.GetHandler(),
		ReadTimeout: 3 * time.Second,
	}
	// run server
	log.Printf("server is listening on %v...", addr)
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("unable to start server due: %v", err)
	}
}
