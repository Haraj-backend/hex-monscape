package main

import (
	"errors"
	"log"
	"net/http"
	"time"

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
	deps, err := initStorageDeps(cfg)
	if err != nil {
		log.Fatalf("unable to initialize storages due: %v", err)
	}

	// initialize play service
	playService, err := play.NewService(play.ServiceConfig{
		GameStorage:    deps.PlayGameStorage,
		PartnerStorage: deps.PlayPartnerStorage,
	})
	if err != nil {
		log.Fatalf("unable to initialize play service due: %v", err)
	}

	// initialize battle service
	battleService, err := battle.NewService(battle.ServiceConfig{
		GameStorage:    deps.BattleGameStorage,
		BattleStorage:  deps.BattleBattleStorage,
		MonsterStorage: deps.BattleMonsterStorage,
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
