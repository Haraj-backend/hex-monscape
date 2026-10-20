# Driven adapters

A driven adapter implements one or more **driven ports** (`XxxStorage` interfaces in `internal/core/service/<ctx>/storage.go`) using a concrete technology. It imports `internal/core/entity` for types and may import the service package only to assert that it satisfies the port. It never contains application logic.

Canonical sources:
- `internal/driven/storage/memory/{gamestrg,battlestrg,monstrg}/`
- `internal/driven/storage/mysql/{gamestrg,...,shared}/`
- `internal/driven/storage/dynamodb/{gamestrg,...,shared}/`

The templates continue the `todo` example from [core.md](core.md).

## Package rules
- Path: `internal/driven/storage/<tech>/<name>strg/`, where `<name>` is the entity (`taskstrg`, `gamestrg`, `monstrg`). The files are `storage.go`, `model.go`, and `storage_test.go`.
- The type is always `Storage`, built by `New(Config) (*Storage, error)`. The memory adapter with no deps may use `New() *Storage`.
- One adapter may satisfy several ports from different services. For example, `monstrg.Storage` satisfies both `play.PartnerStorage` and `battle.MonsterStorage`.
- **Not found returns `nil, nil`.** All other failures are wrapped with context.
- Persistence models (`XxxRow`, `xxxKey`) live in `model.go` with mappers both ways. Never put `db:`/`dynamodbav:` tags on entities.
- Recommended (new, not yet used in this repo): add a compile-time assertion so drift is caught by the compiler:
  ```go
  var _ todo.TaskStorage = (*Storage)(nil)
  ```

## Memory: `internal/driven/storage/memory/taskstrg/storage.go`

```go
package taskstrg

import (
	"context"
	"sync"

	"example.com/todoapp/internal/core/entity"
)

type Storage struct {
	mu   sync.RWMutex
	data map[string]entity.Task
}

func (s *Storage) GetTask(ctx context.Context, taskID string) (*entity.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, exist := s.data[taskID]
	if !exist {
		// if item is not found, returns nil as expected by the port
		return nil, nil
	}
	return &t, nil
}

func (s *Storage) SaveTask(ctx context.Context, task entity.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[task.ID] = task
	return nil
}

func New() *Storage {
	return &Storage{data: make(map[string]entity.Task)}
}
```

(The repo's memory adapters omit the mutex. Add it for anything served concurrently over HTTP.) Seed data is passed as `[]byte` through `Config` (see `memory/monstrg`), not read from disk inside the adapter.

## MySQL (sqlx)

`storage.go`:

```go
package taskstrg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/todoapp/internal/core/entity"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
)

type Storage struct {
	sqlClient *sqlx.DB
}

type Config struct {
	SQLClient *sqlx.DB `validate:"nonnil"`
}

func (c Config) Validate() error {
	return validator.Validate(c)
}

func New(cfg Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &Storage{sqlClient: cfg.SQLClient}, nil
}

func (s *Storage) GetTask(ctx context.Context, taskID string) (*entity.Task, error) {
	var row TaskRow
	query := `SELECT id, title, status, created_at FROM task WHERE id = ?`
	err := s.sqlClient.GetContext(ctx, &row, query, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to get task %s due: %w", taskID, err)
	}
	return row.ToTask(), nil
}

func (s *Storage) SaveTask(ctx context.Context, task entity.Task) error {
	query := `
		REPLACE INTO task (id, title, status, created_at)
		VALUES (:id, :title, :status, :created_at)
	`
	_, err := s.sqlClient.NamedExecContext(ctx, query, NewTaskRow(&task))
	if err != nil {
		return fmt.Errorf("unable to execute query due: %w", err)
	}
	return nil
}
```

`model.go`:

```go
type TaskRow struct {
	ID        string        `db:"id"`
	Title     string        `db:"title"`
	Status    entity.Status `db:"status"`
	CreatedAt int64         `db:"created_at"`
}

func (r *TaskRow) ToTask() *entity.Task {
	return &entity.Task{ID: r.ID, Title: r.Title, Status: r.Status, CreatedAt: r.CreatedAt}
}

func NewTaskRow(t *entity.Task) *TaskRow {
	return &TaskRow{ID: t.ID, Title: t.Title, Status: t.Status, CreatedAt: t.CreatedAt}
}
```

The schema goes in `deploy/local/shared/mysql/schema.sql` and is documented in `docs/db/mysql.md`. Rows shared by several adapters (e.g. `MonsterRow`) go in `mysql/shared/model.go`.

## DynamoDB (aws-sdk-go v1)

```go
type Config struct {
	DynamoClient *dynamodb.DynamoDB `validate:"nonnil"`
	TableName    string             `validate:"nonzero"`
}
```

- `model.go` has a `taskKey{ID string}` with `toDDBKey()`, plus a `taskRow` with `dynamodbav:"..."` tags, `toTaskRow(entity.Task)`, and `(r taskRow) toTask() *entity.Task`.
- Get: `GetItemWithContext`. When `len(output.Item) == 0`, return `nil, nil`. Otherwise `dynamodbattribute.UnmarshalMap` and map to the entity.
- Save: `dynamodbattribute.MarshalMap(toTaskRow(task))` then `PutItemWithContext`.
- Table creation goes in `deploy/local/shared/dynamodb/create-tables.sh` and is documented in `docs/db/dynamodb.md`.

See `internal/driven/storage/dynamodb/gamestrg/` for the full version.

## External services (HTTP APIs, queues, email...)

These are driven adapters too. Define the port in Core in domain terms (e.g. `type Notifier interface { NotifyTaskDone(ctx, task entity.Task) error }`), and implement it in `internal/driven/<kind>/<name>/` with the same `Config`/`New` pattern. The vendor SDK must not leak into the port signature.

## Adapter tests (integration, real infra)

- Write tests against real MySQL/DynamoDB started by `deploy/local/tests/docker-compose.yml`. Do not mock the database.
- Use a `newStorage(t)` helper. Test the round trip (save, get, `require.Equal`) and the not-found case (`require.NoError`, `require.Nil`).
- Use random IDs (`uuid.NewString()`) so tests do not collide on a shared database.
- Shared helpers go in `<tech>/shared/tests.go`. Read connection info from env (`TEST_SQL_DSN`, `LOCALSTACK_ENDPOINT`, `DDB_TABLE_<NAME>_NAME`); see `mysql/shared/tests.go` and `dynamodb/shared/tests.go`.
- Memory adapters can use plain unit tests.
