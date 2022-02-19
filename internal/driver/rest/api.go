package rest

import (
	"net/http"

	"github.com/Haraj-backend/hex-pokebattle/internal/domain/battling"
	"github.com/Haraj-backend/hex-pokebattle/internal/domain/playing"
	"github.com/go-chi/chi/v5"
	"gopkg.in/validator.v2"
)

type API struct {
	playingService  playing.Service
	battlingService battling.Service
}

func (a *API) GetHandler() http.Handler {
	r := chi.NewRouter()
	r.Get("/partners", a.serveGetAvailablePartners)
	r.Route("/games", func(r chi.Router) {
		r.Post("/", a.serveNewGame)
		r.Route("/{game_id}", func(r chi.Router) {
			r.Get("/", a.serveGetGameDetails)
			r.Get("/scenario", a.serveGetNextScenario)
			r.Route("/battle", func(r chi.Router) {
				r.Put("/", a.serveStartBattle)
				r.Get("/", a.serveGetBattleInfo)
				r.Put("/turn", a.serveDecideTurn)
				r.Put("/attack", a.serveAttack)
				r.Put("/surrender", a.serveSurrender)
			})
		})
	})
	return r
}

func (a *API) serveGetAvailablePartners(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (a *API) serveNewGame(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (a *API) serveGetGameDetails(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (a *API) serveGetNextScenario(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (a *API) serveStartBattle(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (a *API) serveGetBattleInfo(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (a *API) serveDecideTurn(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (a *API) serveAttack(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (a *API) serveSurrender(w http.ResponseWriter, r *http.Request) {
	// TODO
}

type APIConfig struct {
	PlayingService  playing.Service  `validate:"nonnil"`
	BattlingService battling.Service `validate:"nonnil"`
}

func (c APIConfig) Validate() error {
	return validator.Validate(c)
}

func NewAPI(cfg APIConfig) (*API, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	a := &API{
		playingService:  cfg.PlayingService,
		battlingService: cfg.BattlingService,
	}
	return a, nil
}
