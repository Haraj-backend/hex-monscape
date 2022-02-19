package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battling"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/playing"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"gopkg.in/validator.v2"
)

type API struct {
	playingService  *playing.Service
	battlingService *battling.Service
}

func (a *API) GetHandler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/partners", a.serveGetAvailablePartners)
	r.Route("/games", func(r chi.Router) {
		r.Post("/", a.serveNewGame)
		r.Route("/{game_id}", func(r chi.Router) {
			r.Get("/", a.serveGetGameDetails)
			r.Get("/scenario", a.serveGetScenario)
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
	partners, err := a.playingService.GetAvailablePartners(r.Context())
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(map[string]interface{}{
		"partners": partners,
	}))
}

func (a *API) serveNewGame(w http.ResponseWriter, r *http.Request) {
	var rb newGameRespBody
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}
	err = rb.Validate()
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}
	game, err := a.playingService.NewGame(r.Context(), rb.PlayerName, rb.PartnerID)
	if err != nil {
		if errors.Is(err, playing.ErrPartnerNotFound) {
			err = NewPartnerNotFoundError()
		}
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(game))
}

func (a *API) serveGetGameDetails(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "game_id")
	game, err := a.playingService.GetGame(r.Context(), gameID)
	if err != nil {
		if errors.Is(err, playing.ErrGameNotFound) {
			err = NewGameNotFoundError()
		}
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(game))
}

func (a *API) serveGetScenario(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "game_id")
	game, err := a.playingService.GetGame(r.Context(), gameID)
	if err != nil {
		if errors.Is(err, playing.ErrGameNotFound) {
			err = NewGameNotFoundError()
		}
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(map[string]interface{}{
		"scenario": game.Scenario,
	}))
}

func (a *API) serveStartBattle(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "game_id")
	battle, err := a.battlingService.StartBattle(r.Context(), gameID)
	if err != nil {
		switch err {
		case battling.ErrGameNotFound:
			err = NewGameNotFoundError()
		case battling.ErrInvalidBattleState:
			err = NewInvalidBattleStateError()
		}
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(battle))
}

func (a *API) serveGetBattleInfo(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "game_id")
	battle, err := a.battlingService.GetBattle(r.Context(), gameID)
	if err != nil {
		switch err {
		case battling.ErrGameNotFound:
			err = NewGameNotFoundError()
		case battling.ErrBattleNotFound:
			err = NewBattleNotFoundError()
		}
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(battle))
}

func (a *API) serveDecideTurn(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "game_id")
	battle, err := a.battlingService.DecideTurn(r.Context(), gameID)
	if err != nil {
		switch err {
		case battling.ErrGameNotFound:
			err = NewGameNotFoundError()
		case battling.ErrBattleNotFound:
			err = NewBattleNotFoundError()
		}
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(battle))
}

func (a *API) serveAttack(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "game_id")
	battle, err := a.battlingService.Attack(r.Context(), gameID)
	if err != nil {
		switch err {
		case battling.ErrGameNotFound:
			err = NewGameNotFoundError()
		case battling.ErrBattleNotFound:
			err = NewBattleNotFoundError()
		}
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(battle))
}

func (a *API) serveSurrender(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "game_id")
	battle, err := a.battlingService.Surrender(r.Context(), gameID)
	if err != nil {
		switch err {
		case battling.ErrGameNotFound:
			err = NewGameNotFoundError()
		case battling.ErrBattleNotFound:
			err = NewBattleNotFoundError()
		}
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(battle))
}

type APIConfig struct {
	PlayingService  *playing.Service  `validate:"nonnil"`
	BattlingService *battling.Service `validate:"nonnil"`
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
