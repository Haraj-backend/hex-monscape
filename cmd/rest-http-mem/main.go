package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/memory/battlestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/memory/gamestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/memory/pokestrg"
	"github.com/Haraj-backend/hex-pokebattle/internal/driver/rest"
	"github.com/Haraj-backend/hex-pokebattle/internal/shared/telemetry"
	"github.com/gosidekick/goconfig"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/play"
)

func main() {
	var cfg config
	err := goconfig.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse configs due: %v", err)
	}

	// initialize tracer exporter
	traceExporter, err := telemetry.NewJaegerTracerProvider(cfg.JaegerEndpointURL, cfg.ServiceName)
	if err != nil {
		log.Fatalf("unable to initialize tracer exporter due: %v", err)
	}

	// initialize telemetry tracer
	telemetryTracer, err := telemetry.NewOpenTelemetryTracer(telemetry.OpenTelemetryConfig{
		Exporter:    *traceExporter,
		ServiceName: cfg.ServiceName,
		BaseContext: context.Background(),
	})
	if err != nil {
		log.Fatalf("unable to initialize telemetry tracer due: %v", err)
	}

	// set singleton tracer
	telemetry.SetTracer(&telemetryTracer)

	// initialize pokemon storage
	partnersData, err := ioutil.ReadFile(cfg.RootData + "/partners.json")
	if err != nil {
		log.Fatalf("unable to read partners data due: %v", err)
	}
	var partners []entity.Pokemon
	err = json.Unmarshal(partnersData, &partners)
	if err != nil {
		log.Fatalf("unable to parse partners data due: %v", err)
	}
	enemiesData, err := ioutil.ReadFile(cfg.RootData + "/enemies.json")
	if err != nil {
		log.Fatalf("unable to read enemies data due: %v", err)
	}
	var enemies []entity.Pokemon
	err = json.Unmarshal(enemiesData, &enemies)
	if err != nil {
		log.Fatalf("unable to parse enemies data due: %v", err)
	}
	pokemonStorage, err := pokestrg.New(
		pokestrg.Config{
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
		ServiceName:    cfg.ServiceName,
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

type config struct {
	Port              string `cfg:"port" cfgDefault:"9186"`
	ServiceName       string `cfg:"service_name" cfgDefault:"rest-http-mem"`
	RootData          string `cfg:"root_data" cfgDefault:"cmd/rest-http-mem"`
	JaegerEndpointURL string `cfg:"jaeger_endpoint_url" cfgDefault:"http://localhost:14268/api/traces"`
}
