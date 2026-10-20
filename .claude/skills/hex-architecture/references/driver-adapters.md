# Driver adapters

A driver adapter translates a driver actor (an HTTP client, a Lambda event, a CLI invocation, a queue message) into calls on core `Service` interfaces, and translates results and errors back. It depends **only** on `internal/core/service/<ctx>` (plus `internal/core/entity` for types). It never imports driven adapters.

Canonical sources: `internal/driver/rest/{api.go,error.go,req_body.go,resp_body.go}`, `cmd/lambda/main.go`.

The templates continue the `todo` example from [core.md](core.md).

## REST: `internal/driver/rest/api.go`

```go
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

	"example.com/todoapp/internal/core/service/todo"
)

type APIConfig struct {
	TodoService todo.Service `validate:"nonnil"`
}

func (c APIConfig) Validate() error {
	return validator.Validate(c)
}

func NewAPI(cfg APIConfig) (*API, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	a := &API{todoService: cfg.TodoService}
	return a, nil
}

type API struct {
	todoService todo.Service
}

func (a *API) GetHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/health", a.serveHealthCheck)
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", a.serveCreateTask)
		r.Put("/{task_id}/complete", a.serveCompleteTask)
	})

	return r
}

func (a *API) serveHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *API) serveCreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// decode & validate request body
	var rb createTaskReqBody
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
	// call the service
	task, err := a.todoService.CreateTask(ctx, rb.Title)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(task))
}

func (a *API) serveCompleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskID := chi.URLParam(r, "task_id")
	task, err := a.todoService.CompleteTask(ctx, taskID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(task))
}

// handleServiceError maps service sentinel errors into REST errors. Use
// errors.Is (not ==) so wrapped errors are still mapped correctly.
func handleServiceError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, todo.ErrTaskNotFound):
		err = NewTaskNotFoundError()
	case errors.Is(err, todo.ErrInvalidTaskStatus):
		err = NewInvalidTaskStatusError()
	default:
		err = NewInternalServerError(err.Error())
	}
	render.Render(w, r, NewErrorResp(err))
}
```

Rules:
- Handlers are named `serveXxx` and follow the flow **decode, validate, call service, map error, render**. Keep them thin: no application logic in handlers.
- All service error mapping goes through one `handleServiceError`. Do not map errors inline in handlers.
- Route params use snake_case (`{task_id}`), matching JSON keys.

## Request bodies: `internal/driver/rest/req_body.go`

```go
package rest

import "gopkg.in/validator.v2"

type createTaskReqBody struct {
	Title string `json:"title" validate:"nonzero"`
}

func (rb createTaskReqBody) Validate() error {
	err := validator.Validate(rb)
	if err != nil {
		return NewBadRequestError(err.Error())
	}
	return nil
}
```

## Errors: `internal/driver/rest/error.go`

Copy `Error`, `Error()`, `Is()`, `NewInternalServerError` and `NewBadRequestError` from this repo's `internal/driver/rest/error.go` verbatim. Then add one constructor per service error:

```go
func NewTaskNotFoundError() *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Err:        "ERR_TASK_NOT_FOUND",
		Message:    "task is not found",
	}
}

func NewInvalidTaskStatusError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Err:        "ERR_INVALID_TASK_STATUS",
		Message:    "invalid task status",
	}
}
```

Status mapping: not found → 404, invalid state → 409, bad input → 400, everything else → 500. Document every `ERR_*` code in `docs/api/rest-api.md`.

## Response envelope: `internal/driver/rest/resp_body.go`

Copy `RespBody`, `Render`, `NewSuccessResp` and `NewErrorResp` from this repo verbatim. Every response looks like:

```json
{ "ok": true, "data": { "id": "...", "title": "...", "status": "OPEN", "created_at": 1700000000 }, "ts": 1700000000 }
{ "ok": false, "err": "ERR_TASK_NOT_FOUND", "msg": "task is not found", "ts": 1700000000 }
```

`NewSuccessResp` wraps data with `conjson.NewMarshaler(data, transform.ConventionalKeys())`, which converts `CreatedAt` to `created_at`. This is why entities need no JSON tags. To return a list, wrap it in a map: `NewSuccessResp(map[string]interface{}{"tasks": tasks})`.

## Other drivers

- **AWS Lambda (HTTP)**: reuse the REST handler. In `cmd/lambda/main.go`, serve `api.GetHandler()` with `gateway.ListenAndServe` (`github.com/apex/gateway`) on Lambda and `http.ListenAndServe` locally. See this repo's `cmd/lambda/main.go`.
- **CLI / queue consumer / cron**: create `internal/driver/<tech>/` with a `Config` holding the Service interfaces and a `New(Config)` constructor. Parse the input into primitives, call the service, and map sentinel errors to the transport's own semantics (exit code, ack/nack, retry). Add a `cmd/<deliverable>/` for it.
- **Driver tests**: this repo has none for REST yet. When adding them, use `httptest` with a fake `Service` implementation (hand-written, like the core mocks) and assert status code plus envelope.
