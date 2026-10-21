# Core: Service, Entity, Storage ports

Canonical sources in this repo:
- `internal/core/service/play/{service.go,storage.go,service_test.go}`
- `internal/core/service/battle/` (a service with several ports and entity state transitions)
- `internal/core/entity/{game.go,battle.go,monster.go}`
- `internal/core/testutil/util.go`

The templates below use a running example: module `example.com/todoapp`, context `todo`, entity `Task`. Rename to fit.

Core may import only the stdlib, `internal/core/...`, and small pure libraries (`google/uuid`, `gopkg.in/validator.v2`). It must **never** import `net/http`, `database/sql`, `sqlx`, `aws-sdk-go`, `chi`, or anything under `internal/driver`, `internal/driven`, or `cmd`.

## Entity: `internal/core/entity/task.go`

```go
package entity

import (
	"errors"

	"github.com/google/uuid"
	"gopkg.in/validator.v2"
)

var ErrInvalidStatus = errors.New("invalid status for given action")

type Status string

const (
	StatusOpen Status = "OPEN"
	StatusDone Status = "DONE"
)

// no json/db tags: the REST layer converts keys with conjson, storage uses its own row models
type Task struct {
	ID        string
	Title     string
	Status    Status
	CreatedAt int64
}

// Complete marks the task as done. The task must be OPEN, otherwise
// the action will be rejected.
func (t *Task) Complete() error {
	if t.Status != StatusOpen {
		return ErrInvalidStatus
	}
	t.Status = StatusDone
	return nil
}

type TaskConfig struct {
	Title     string `validate:"nonzero"`
	CreatedAt int64  `validate:"nonzero"`
}

func (c TaskConfig) Validate() error {
	return validator.Validate(c)
}

// NewTask returns new instance of Task based on given config.
func NewTask(cfg TaskConfig) (*Task, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	t := &Task{
		ID:        uuid.NewString(),
		Title:     cfg.Title,
		Status:    StatusOpen,
		CreatedAt: cfg.CreatedAt,
	}
	return t, nil
}
```

Rules:
- Put behavior on the entity when it only concerns that entity's state (see `Battle.PartnerAttack`, `Game.IncBattleWon`). Coordination across entities and storage belongs in the service.
- Pure data structs with no invariants (like `entity.Monster`) may skip the constructor.
- Time is passed in through config (`CreatedAt`), not read inside the entity, which keeps entities deterministic in tests.

## Driven ports: `internal/core/service/todo/storage.go`

```go
package todo

import (
	"context"

	"example.com/todoapp/internal/core/entity"
)

type TaskStorage interface {
	// GetTask returns task instance for given taskID from storage. Returns nil
	// when given taskID is not found.
	GetTask(ctx context.Context, taskID string) (*entity.Task, error)

	// SaveTask is used for saving (insert or replace) task instance in storage.
	SaveTask(ctx context.Context, task entity.Task) error
}
```

Rules:
- The port lives in the **consumer's** package and contains only the methods that service uses. Another service needing tasks declares its own `TaskStorage`.
- Doc comments must state the not-found contract (`Returns nil when ... not found`).
- Pass entities by value on save and return pointers on get, as in the existing ports.

## Service: `internal/core/service/todo/service.go`

```go
package todo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"example.com/todoapp/internal/core/entity"
	"gopkg.in/validator.v2"
)

var (
	ErrTaskNotFound      = errors.New("task is not found")
	ErrInvalidTaskStatus = errors.New("invalid task status")
)

type Service interface {
	// CreateTask creates new open task with given title and saves it to storage.
	// Upon success it returns the saved task instance.
	CreateTask(ctx context.Context, title string) (*entity.Task, error)

	// CompleteTask marks task with given id as done. Returns `ErrTaskNotFound`
	// when the task is not found & `ErrInvalidTaskStatus` when the task is
	// already done.
	CompleteTask(ctx context.Context, taskID string) (*entity.Task, error)
}

type service struct {
	taskStorage TaskStorage
}

func (s *service) CreateTask(ctx context.Context, title string) (*entity.Task, error) {
	// initiate new task instance
	task, err := entity.NewTask(entity.TaskConfig{
		Title:     title,
		CreatedAt: time.Now().Unix(),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize task instance due: %w", err)
	}
	// store the task instance on storage
	err = s.taskStorage.SaveTask(ctx, *task)
	if err != nil {
		return nil, fmt.Errorf("unable to save task instance due: %w", err)
	}
	return task, nil
}

func (s *service) CompleteTask(ctx context.Context, taskID string) (*entity.Task, error) {
	// get task instance, check err first then nil
	task, err := s.taskStorage.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("unable to get task instance due: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	// complete the task, map entity error into service error
	err = task.Complete()
	if err != nil {
		return nil, ErrInvalidTaskStatus
	}
	// save the updated task
	err = s.taskStorage.SaveTask(ctx, *task)
	if err != nil {
		return nil, fmt.Errorf("unable to save task instance due: %w", err)
	}
	return task, nil
}

type ServiceConfig struct {
	TaskStorage TaskStorage `validate:"nonnil"`
}

func (c ServiceConfig) Validate() error {
	return validator.Validate(c)
}

// NewService returns new instance of service.
func NewService(cfg ServiceConfig) (Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &service{
		taskStorage: cfg.TaskStorage,
	}
	return s, nil
}
```

Rules:
- One method per use case. Its parameters are primitives or entities, never transport types (no `*http.Request`, no request body structs).
- Return sentinel errors **unwrapped** for expected outcomes (not found, invalid state) so drivers can map them. Wrap unexpected infra errors with `%w`.
- Validation of input shape (non-empty strings) happens in the driver; invariants happen in entity constructors and methods.
- If the use cases are hard to reason about, draw a class diagram (`docs/diagrams/`).

## Service tests: `internal/core/service/todo/service_test.go`

Key points (see the header comment in `internal/core/service/play/service_test.go`, which explains why):
- Use the external package `package todo_test`.
- Use **hand-written, map-backed mocks** with `retErrOnXxx` toggles to inject failures. Do not use real storage and do not add a mock generator.
- Write a table-driven `TestNewService` covering each nil dependency, then one `TestServiceXxx` per use case covering the success path, each sentinel error, and each injected storage error.
- A `newService()` helper returns an output struct holding the service, the mocks, and the fixtures.

```go
package todo_test

import (
	"context"
	"errors"
	"testing"

	"example.com/todoapp/internal/core/entity"
	"example.com/todoapp/internal/core/service/todo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	testCases := []struct {
		Name    string
		Config  todo.ServiceConfig
		IsError bool
	}{
		{
			Name:    "Test Missing Task Storage",
			Config:  todo.ServiceConfig{TaskStorage: nil},
			IsError: true,
		},
		{
			Name:    "Test Valid Config",
			Config:  todo.ServiceConfig{TaskStorage: newMockTaskStorage()},
			IsError: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := todo.NewService(testCase.Config)
			require.Equal(t, testCase.IsError, (err != nil), "unexpected error")
		})
	}
}

func TestServiceCompleteTask(t *testing.T) {
	// initialize new service
	output := newService()

	// create task then complete it
	task, err := output.Service.CreateTask(context.Background(), "write tests")
	require.NoError(t, err)
	task, err = output.Service.CompleteTask(context.Background(), task.ID)
	require.NoError(t, err)
	require.Equal(t, entity.StatusDone, task.Status)

	// complete again, should return invalid status
	_, err = output.Service.CompleteTask(context.Background(), task.ID)
	require.Equal(t, todo.ErrInvalidTaskStatus, err)

	// complete unknown task, should return not found
	_, err = output.Service.CompleteTask(context.Background(), uuid.NewString())
	require.Equal(t, todo.ErrTaskNotFound, err)

	// set error on get task, should return error
	output.TaskStorage.SetRetErrOnGetTask(true)
	_, err = output.Service.CompleteTask(context.Background(), task.ID)
	require.ErrorIs(t, err, ErrIntentionalError)
}

func newService() *newServiceOutput {
	taskStorage := newMockTaskStorage()
	svc, _ := todo.NewService(todo.ServiceConfig{TaskStorage: taskStorage})
	return &newServiceOutput{Service: svc, TaskStorage: taskStorage}
}

type newServiceOutput struct {
	Service     todo.Service
	TaskStorage *mockTaskStorage
}

type mockTaskStorage struct {
	data             map[string]entity.Task
	retErrOnGetTask  bool
	retErrOnSaveTask bool
}

func (s *mockTaskStorage) SetRetErrOnGetTask(retErr bool)  { s.retErrOnGetTask = retErr }
func (s *mockTaskStorage) SetRetErrOnSaveTask(retErr bool) { s.retErrOnSaveTask = retErr }

func (s *mockTaskStorage) GetTask(ctx context.Context, taskID string) (*entity.Task, error) {
	if s.retErrOnGetTask {
		return nil, ErrIntentionalError
	}
	task, ok := s.data[taskID]
	if !ok {
		return nil, nil
	}
	return &task, nil
}

func (s *mockTaskStorage) SaveTask(ctx context.Context, task entity.Task) error {
	if s.retErrOnSaveTask {
		return ErrIntentionalError
	}
	s.data[task.ID] = task
	return nil
}

func newMockTaskStorage() *mockTaskStorage {
	return &mockTaskStorage{data: map[string]entity.Task{}}
}

var ErrIntentionalError = errors.New("intentional error")
```

## Entity tests: `internal/core/entity/task_test.go`

Use `package entity_test` with table-driven tests per constructor and per behavior method (see `internal/core/entity/battle_test.go`). Put random fixture generators shared across packages in `internal/core/testutil` (e.g. `NewTestMonster()`).
