package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"gopkg.in/validator.v2"

	"github.com/Haraj-backend/hex-monscape/internal/core/service/battle"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/play"
)

type APIConfig struct {
	PlayingService play.Service   `validate:"nonnil"`
	BattleService  battle.Service `validate:"nonnil"`
}

func (c APIConfig) Validate() error {
	return validator.Validate(c)
}

func NewAPI(cfg APIConfig) (*API, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	a := &API{
		playService:   cfg.PlayingService,
		battleService: cfg.BattleService,
	}
	return a, nil
}

type API struct {
	playService   play.Service
	battleService battle.Service
}

func (a *API) GetHandler() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.AllowAll().Handler)
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
	ctx := r.Context()
	partners, err := a.playService.GetAvailablePartners(ctx)
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(map[string]interface{}{
		"partners": partners,
	}))
}

func (a *API) serveNewGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var rb newGameReqBody
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
	game, err := a.playService.NewGame(ctx, rb.PlayerName, rb.PartnerID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(game))
}

func (a *API) serveGetGameDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	game, err := a.playService.GetGame(ctx, gameID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(game))
}

func (a *API) serveGetScenario(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	game, err := a.playService.GetGame(ctx, gameID)
	if err != nil {
		if errors.Is(err, play.ErrGameNotFound) {
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
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	bt, err := a.battleService.StartBattle(ctx, gameID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func (a *API) serveGetBattleInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	bt, err := a.battleService.GetBattle(ctx, gameID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func (a *API) serveDecideTurn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	bt, err := a.battleService.DecideTurn(ctx, gameID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func (a *API) serveAttack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	bt, err := a.battleService.Attack(ctx, gameID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func (a *API) serveSurrender(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	bt, err := a.battleService.Surrender(ctx, gameID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func handleServiceError(w http.ResponseWriter, r *http.Request, err error) {
	switch err {
	case battle.ErrGameNotFound:
		err = NewGameNotFoundError()
	case battle.ErrBattleNotFound:
		err = NewBattleNotFoundError()
	case battle.ErrInvalidBattleState:
		err = NewInvalidBattleStateError()
	case play.ErrGameNotFound:
		err = NewGameNotFoundError()
	case play.ErrPartnerNotFound:
		err = NewPartnerNotFoundError()
	default:
		err = NewInternalServerError(err.Error())
	}
	render.Render(w, r, NewErrorResp(err))
}
