package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/battlestrg"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/gamestrg"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/monstrg"
	"github.com/Haraj-backend/hex-monscape/internal/driver/rest"
	"github.com/gosidekick/goconfig"

	"github.com/Haraj-backend/hex-monscape/internal/core/service/battle"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/play"
)

func main() {
	var cfg config
	err := goconfig.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse configs due: %v", err)
	}

	// initialize pokemon storage
	partnersData, err := ioutil.ReadFile("/partners.json")
	if err != nil {
		log.Fatalf("unable to read partners data due: %v", err)
	}
	var partners []entity.Monster
	err = json.Unmarshal(partnersData, &partners)
	if err != nil {
		log.Fatalf("unable to parse partners data due: %v", err)
	}
	enemiesData, err := ioutil.ReadFile("/enemies.json")
	if err != nil {
		log.Fatalf("unable to read enemies data due: %v", err)
	}
	var enemies []entity.Monster
	err = json.Unmarshal(enemiesData, &enemies)
	if err != nil {
		log.Fatalf("unable to parse enemies data due: %v", err)
	}
	pokemonStorage, err := monstrg.New(
		monstrg.Config{
			Partners: partners,
			Enemies:  enemies,
		},
	)
	if err != nil {
		log.Fatalf("unable to initialize pokemon storage due: %v", err)
	}
	// initialize game storage
	gameStorage := gamestrg.New()
	// initialize battle storage
	battleStorage := battlestrg.New()
	// initialize play service
	playService, err := play.NewService(play.ServiceConfig{
		GameStorage:    gameStorage,
		PartnerStorage: pokemonStorage,
	})
	if err != nil {
		log.Fatalf("unable to initialize play service due: %v", err)
	}
	// initialize battle service
	battleService, err := battle.NewService(battle.ServiceConfig{
		GameStorage:    gameStorage,
		BattleStorage:  battleStorage,
		PokemonStorage: pokemonStorage,
	})
	if err != nil {
		log.Fatalf("unable to initialize battle service due: %v", err)
	}
	// initialize rest api
	api, err := rest.NewAPI(rest.APIConfig{
		PlayingService: playService,
		BattleService:  battleService,
		ServiceName:    "",
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