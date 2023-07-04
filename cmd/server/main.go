package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/battlestrg"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/gamestrg"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/monstrg"
	"github.com/Haraj-backend/hex-monscape/internal/driver/rest"
	"github.com/gosidekick/goconfig"

	"github.com/Haraj-backend/hex-monscape/internal/core/service/battle"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/play"
)

func main() {
	// initialize configs
	var cfg config
	err := goconfig.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse configs due: %v", err)
	}

	// initialize monster storage
	monsterData, err := ioutil.ReadFile(cfg.Storage.Memory.MonsterDataPath)
	if err != nil {
		log.Fatalf("unable to read monster data due: %v", err)
	}
	monsterStorage, err := monstrg.New(monstrg.Config{MonsterData: monsterData})
	if err != nil {
		log.Fatalf("unable to initialize monster storage due: %v", err)
	}

	// initialize game storage
	gameStorage := gamestrg.New()

	// initialize battle storage
	battleStorage := battlestrg.New()

	// initialize play service
	playService, err := play.NewService(play.ServiceConfig{
		GameStorage:    gameStorage,
		PartnerStorage: monsterStorage,
	})
	if err != nil {
		log.Fatalf("unable to initialize play service due: %v", err)
	}

	// initialize battle service
	battleService, err := battle.NewService(battle.ServiceConfig{
		GameStorage:    gameStorage,
		BattleStorage:  battleStorage,
		MonsterStorage: monsterStorage,
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
	log.Printf("server is listening on :%v...", cfg.Port)
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("unable to start server due: %v", err)
	}
}
