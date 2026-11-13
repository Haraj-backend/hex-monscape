# Driven adapters

A driven adapter implements one or more **driven ports** (`XxxStorage` interfaces in `internal/core/service/<ctx>/storage.go`) using a concrete technology. It imports `internal/core/entity` for types and may import the service package only to assert that it satisfies the port. It never contains application logic.

Canonical sources (one package per tech, one file per entity, plus `model.go`):
- `internal/driven/storage/memory/{game,battle,monster,model}.go`
- `internal/driven/storage/mysql/{game,battle,monster,model}.go`
- `internal/driven/storage/ddb/{game,battle,monster,model}.go` (named `ddb` so it does not clash with the AWS SDK `dynamodb` package)

The templates continue the `todo` example from [core.md](core.md).

## Package rules
- **Pick the granularity that keeps the code easy to navigate** (see "Package granularity" in SKILL.md). The default is one package per kind, e.g. `internal/driven/storage` with one file per concern (`task.go`, `progress.go`) and a shared `model.go`. Types are named after what they store (`TaskStorage`, `NewTaskStorage(TaskStorageConfig)`).
- Split by tech (`storage/memory`, `storage/mysql`) only when there are alternative implementations of the same ports that `cmd/` switches between. hex-monscape does this because it has three storage techs, and the templates below follow its layout. Split a tech package by entity only when it becomes too big to navigate.
- A memory adapter with no deps may use a constructor without config (`NewTaskStorage() *TaskStorage`).
- One adapter may satisfy several ports from different services. For example, `memory.MonsterStorage` satisfies both `play.PartnerStorage` and `battle.MonsterStorage`.
- **Not found returns `nil, nil`.** All other failures are wrapped with context.
- Persistence models (`xxxRow`, `xxxKey`) live in `model.go` with mappers both ways. Keep them unexported: nothing outside the adapter package needs them. Never put `db:`/`dynamodbav:` tags on entities.
- Do not add `var _ Port = (*Storage)(nil)` assertions. The team finds them confusing, especially for junior developers. Wiring in `cmd/` already catches a missing method at compile time.

## Memory: `internal/driven/storage/memory/task.go`

```go
package memory

import (
	"context"
	"sync"

	"example.com/todoapp/internal/core/entity"
)

type TaskStorage struct {
	mu   sync.RWMutex
	data map[string]entity.Task
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{data: make(map[string]entity.Task)}
}

func (s *TaskStorage) GetTask(ctx context.Context, taskID string) (*entity.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, exist := s.data[taskID]
	if !exist {
		// if item is not found, returns nil as expected by the port
		return nil, nil
	}

	return &t, nil
}

func (s *TaskStorage) SaveTask(ctx context.Context, task entity.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[task.ID] = task

	return nil
}
```

Keep the mutex: the adapter is served concurrently over HTTP. Seed data is passed as `[]byte` through config (see `memory.NewMonsterStorage`), not read from disk inside the adapter.

## MySQL (sqlx)

`internal/driven/storage/mysql/task.go`:

```go
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/todoapp/internal/core/entity"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
)

type TaskStorage struct {
	sqlClient *sqlx.DB
}

type TaskStorageConfig struct {
	SQLClient *sqlx.DB `validate:"nonnil"`
}

func (c TaskStorageConfig) Validate() error {
	return validator.Validate(c)
}

func NewTaskStorage(cfg TaskStorageConfig) (*TaskStorage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	return &TaskStorage{sqlClient: cfg.SQLClient}, nil
}

func (s *TaskStorage) GetTask(ctx context.Context, taskID string) (*entity.Task, error) {
	var row taskRow
	query := `SELECT id, title, status, created_at FROM task WHERE id = ?`
	err := s.sqlClient.GetContext(ctx, &row, query, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to get task %s due: %w", taskID, err)
	}

	return row.toTask(), nil
}

func (s *TaskStorage) SaveTask(ctx context.Context, task entity.Task) error {
	query := `
		REPLACE INTO task (id, title, status, created_at)
		VALUES (:id, :title, :status, :created_at)
	`
	_, err := s.sqlClient.NamedExecContext(ctx, query, newTaskRow(&task))
	if err != nil {
		return fmt.Errorf("unable to execute query due: %w", err)
	}

	return nil
}
```

`internal/driven/storage/mysql/model.go`:

```go
type taskRow struct {
	ID        string        `db:"id"`
	Title     string        `db:"title"`
	Status    entity.Status `db:"status"`
	CreatedAt int64         `db:"created_at"`
}

func newTaskRow(t *entity.Task) *taskRow {
	return &taskRow{ID: t.ID, Title: t.Title, Status: t.Status, CreatedAt: t.CreatedAt}
}

func (r *taskRow) toTask() *entity.Task {
	return &entity.Task{ID: r.ID, Title: r.Title, Status: r.Status, CreatedAt: r.CreatedAt}
}
```

The schema goes in `deploy/local/shared/mysql/schema.sql` and is documented in `docs/db/mysql.md`. All row models of a tech live in that tech's `model.go`, so rows used by several storages (e.g. `monsterRow`, embedded in `gameRow`) need no separate package.

## DynamoDB (aws-sdk-go v1)

`internal/driven/storage/ddb/task.go`:

```go
type TaskStorageConfig struct {
	DynamoClient *dynamodb.DynamoDB `validate:"nonnil"`
	TableName    string             `validate:"nonzero"`
}
```

- `model.go` has a `taskKey{ID string}` with `toDDBKey()`, plus a `taskRow` with `dynamodbav:"..."` tags, `newTaskRow(entity.Task)`, and `(r taskRow) toTask() *entity.Task`.
- Get: `GetItemWithContext`. When `len(output.Item) == 0`, return `nil, nil`. Otherwise `dynamodbattribute.UnmarshalMap` and map to the entity.
- Save: `dynamodbattribute.MarshalMap(newTaskRow(task))` then `PutItemWithContext`.
- Table creation goes in `deploy/local/shared/dynamodb/create-tables.sh` and is documented in `docs/db/dynamodb.md`.

See `internal/driven/storage/ddb/game.go` for the full version.

## External services (HTTP APIs, queues, email...)

These are driven adapters too. Define the port in Core in domain terms (e.g. `type Notifier interface { NotifyTaskDone(ctx, task entity.Task) error }`), and implement it in `internal/driven/<kind>/<name>/` with the same `XxxConfig`/`NewXxx` pattern. The vendor SDK must not leak into the port signature.

## Adapter tests (integration, real infra)

- Write tests against real MySQL/DynamoDB started by `deploy/local/tests/docker-compose.yml`. Do not mock the database.
- Use a `newXxxStorage(t)` helper (e.g. `newGameStorage(t)`). Test the round trip (save, get, `require.Equal`) and the not-found case (`require.NoError`, `require.Nil`).
- Use random IDs (`uuid.NewString()`) so tests do not collide on a shared database.
- Tests live in the adapter package itself (`package mysql`), so they can use the unexported row models. Shared helpers go in `<tech>/helper_test.go`, so they never ship in the binary. Read connection info from env (`TEST_SQL_DSN`, `LOCALSTACK_ENDPOINT`, `DDB_TABLE_<NAME>_NAME`); see `mysql/helper_test.go` and `ddb/helper_test.go`.
- Memory adapters can use plain unit tests.
