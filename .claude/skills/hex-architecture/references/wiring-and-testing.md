# Wiring (`cmd/`), build, and testing

`cmd/<deliverable>/` is the **composition root** and the only place where concrete adapters are constructed and handed to Core. Canonical sources: `cmd/server/{main.go,config.go,dep.go}`, `cmd/lambda/main.go`.

## `config.go`: env config via goconfig

`github.com/gosidekick/goconfig` derives env var names from nested `cfg` tags joined with `_` and upper-cased. `Storage.MySQL.SQLDSN` with tags `storage` → `mysql` → `sql_dsn` becomes `STORAGE_MYSQL_SQL_DSN`.

```go
package main

type config struct {
	Port    string        `cfg:"port" cfgDefault:"9186"`
	Storage storageConfig `cfg:"storage"`
}

type storageConfig struct {
	Type  storageType        `cfg:"type" cfgDefault:"memory"`
	MySQL storageMySQLConfig `cfg:"mysql"`
}

type storageMySQLConfig struct {
	SQLDSN string `cfg:"sql_dsn"`
}

type storageType string

const (
	storageTypeMemory storageType = "memory"
	storageTypeMySQL  storageType = "mysql"
)
```

Use `cfgRequired:"true"` for deliverables with a single fixed storage (see `cmd/lambda`).

## `dep.go`: build adapters, typed by ports

```go
type storageDeps struct {
	TodoTaskStorage todo.TaskStorage // field typed by the PORT, named <Service><Port>
}

func initStorageDeps(cfg config) (*storageDeps, error) {
	var deps storageDeps
	switch cfg.Storage.Type {
	case storageTypeMemory:
		deps.TodoTaskStorage = memtaskstrg.New()
	case storageTypeMySQL:
		sqlClient, err := sqlx.Open("mysql", cfg.Storage.MySQL.SQLDSN)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize sql client due: %v", err)
		}
		taskStorage, err := sqltaskstrg.New(sqltaskstrg.Config{SQLClient: sqlClient})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize task storage due: %v", err)
		}
		deps.TodoTaskStorage = taskStorage
	default:
		return nil, fmt.Errorf("unknown storage type: %v", cfg.Storage.Type)
	}
	return &deps, nil
}
```

Import aliases are prefixed by tech: `memtaskstrg`, `sqltaskstrg`, `ddbtaskstrg`.

## `main.go`: wiring sequence

```go
func main() {
	// initialize configs
	var cfg config
	err := goconfig.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse configs due: %v", err)
	}
	// initialize storages
	deps, err := initStorageDeps(cfg)
	if err != nil {
		log.Fatalf("unable to initialize storages due: %v", err)
	}
	// initialize todo service
	todoService, err := todo.NewService(todo.ServiceConfig{TaskStorage: deps.TodoTaskStorage})
	if err != nil {
		log.Fatalf("unable to initialize todo service due: %v", err)
	}
	// initialize rest api
	api, err := rest.NewAPI(rest.APIConfig{TodoService: todoService})
	if err != nil {
		log.Fatalf("unable to initialize rest api due: %v", err)
	}
	// run server
	server := &http.Server{Addr: ":" + cfg.Port, Handler: api.GetHandler(), ReadTimeout: 3 * time.Second}
	log.Printf("[INFO] server is listening on :%v...", cfg.Port)
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("unable to start server due: %v", err)
	}
}
```

Blank-import DB drivers (`_ "github.com/go-sql-driver/mysql"`) in `main.go`, not in adapters.

## Testing strategy

| Layer | How | Infra needed |
|---|---|---|
| Entities | Table-driven unit tests, `package entity_test` | none |
| Services | Hand-written mocks of ports, `package <ctx>_test` | none |
| Driven adapters | Integration tests against real MySQL/DynamoDB(localstack) | docker-compose |
| Driver adapters | `httptest` + fake `Service` (optional) | none |

Test the essential scenarios: success, each sentinel error, and each dependency failure. 100% coverage is not the goal.

## Docker and Makefile

- `build/package/<deliverable>/Dockerfile`: `golang:<ver>-alpine`. Copy `go.mod`/`go.sum`, run `go mod download`, then copy `cmd/<deliverable>` and `internal`. `server` uses `air` for hot reload during local runs.
- `build/package/tests/Dockerfile`: copies `internal` and runs `go test -v -coverprofile=./.cov-report/profile.cov ./...`.
- `deploy/local/tests/docker-compose.yml`: infra services with healthchecks (mysql with `schema.sql` mounted to `/docker-entrypoint-initdb.d/`, localstack with `create-tables.sh` mounted to `/etc/localstack/init/ready.d/`), plus a `tests` service that `depends_on: condition: service_healthy` and gets the `TEST_*` / `LOCALSTACK_ENDPOINT` / `DDB_TABLE_*` env vars.
- `deploy/local/run/<variant>/docker-compose.yml`: one per runnable variant (`rest-memory`, `rest-mysql`, ...), passing `STORAGE_TYPE` etc. as env.
- Makefile targets: `run` (default variant), `run-<variant>` (`docker compose down --remove-orphans` then `up --build --attach=server`), `test` (`up --build --attach=tests --exit-code-from=tests`), and `test-build-all` (build every image).
- CI (`.github/workflows/ci.yml`): `make test`, upload coverage, `make test-build-all`.

Copy these files from this repo and adjust names. Do not re-invent them.
