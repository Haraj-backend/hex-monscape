---
name: hex-architecture
description: Build and review Go applications using Haraj's Hexagonal Architecture (ports & adapters, without DDD), as implemented in hex-monscape. Use when scaffolding a new Go service/project, adding a use case, core service, entity, or storage port, adding a driver adapter (REST API, Lambda, CLI, queue consumer) or a driven adapter (memory, MySQL, DynamoDB, external API), wiring dependencies in cmd/, or reviewing code for hex architecture compliance (e.g. "does this follow hex?", "core imports infra").
---

# Hexagonal Architecture (Haraj style)

Source of truth: [docs/reference/hex-architecture.md](../../../docs/reference/hex-architecture.md) and [docs/reference/project-methodology.md](../../../docs/reference/project-methodology.md). This repo (`hex-monscape`) is the canonical example, so read the real files cited below whenever a template is not enough.

## Principles

1. Clearly divide the `inside` (Core) from the `outside` (Actors and Adapters) of the application.
2. Dependencies always point **inward**. Core imports nothing from adapters, `cmd/`, or infrastructure libraries.
3. Isolate the boundary with **ports** (interfaces defined in Core) and **adapters** (implementations outside Core).
4. **Keep it simple and avoid overengineering.** Write code a junior developer can read and maintain. Do not add clever Go idioms, extra layers or abstractions, generic helpers, functional options, or "nice to have" safety tricks such as `var _ Port = (*Storage)(nil)` unless the code actually needs them. When in doubt, follow what this repo already does.

The four pillars:

| Pillar | What it is | In this repo |
|---|---|---|
| **Core** | The application logic: how the app behaves for each use case. Made of `Service`, `Entity`, and `Storage` interfaces. | `internal/core/` |
| **Actors** | Driver actors start an interaction (an HTTP client). Driven actors are called by the app (MySQL). | – |
| **Ports** | Interfaces inside Core. **Driver port** = `Service` interface. **Driven port** = `XxxStorage` interfaces. | `service/<ctx>/service.go`, `service/<ctx>/storage.go` |
| **Adapters** | Driver adapters translate a driver actor into Service calls. Driven adapters implement Storage ports. | `internal/driver/rest`, `internal/driven/storage/<tech>/<name>strg` |

Use the term **application logic** (not "business logic"). The team applies hex **without DDD**: there are no aggregates, value objects, or domain events. A service may update several entities directly when the app is small.

## Project layout

```
cmd/<deliverable>/              composition root: main.go, config.go, dep.go (only place that wires adapters)
internal/core/entity/           flat package of entities (plain structs + behavior methods)
internal/core/service/<ctx>/    one package per context: service.go (driver port + impl), storage.go (driven ports)
internal/core/testutil/         fixture generators shared by tests
internal/driver/<tech>/         driver adapters, e.g. rest/{api,error,req_body,resp_body}.go
internal/driven/storage/<tech>/<name>strg/   driven adapters: storage.go, model.go, storage_test.go
internal/driven/storage/<tech>/shared/       row models + test helpers shared by that tech
build/package/<deliverable>/Dockerfile
deploy/local/{run/<variant>,tests,shared}/   docker-compose for running & integration tests
docs/{api,db,reference}/        API spec, DB schemas, architecture notes
Makefile                        run, run-<variant>, test, test-build-all
```

## Conventions

| Thing | Convention |
|---|---|
| Constructors | Every constructor (entity, service, adapter, API) takes a `XxxConfig` struct with `validate` tags plus a `Validate()` method using `gopkg.in/validator.v2` (`nonzero`, `nonnil`). The constructor calls `cfg.Validate()` first. |
| Service | Interface `Service` (driver port) with **doc comments describing behavior and returned errors** for each method. Unexported impl `service`. `NewService(ServiceConfig) (Service, error)`. One package per context (`play`, `battle`). |
| Driven ports | `XxxStorage` interfaces in `service/<ctx>/storage.go`, owned by the consuming service. Duplicating a port across services (e.g. `play.GameStorage` and `battle.GameStorage`) is intended. |
| Not found | Storage returns `nil, nil`. The service converts that into a sentinel error. |
| Errors | Sentinel `var ErrXxx = errors.New("...")` in the service package. Infra errors are wrapped as `fmt.Errorf("unable to <do X> due: %w", err)`. Entity errors (e.g. `entity.ErrInvalidState`) are mapped to service errors. |
| Entities | Exported fields, **no JSON/DB tags**. Typed string enums (`type State string`). Behavior lives in pointer-receiver methods that guard state. `NewXxx(XxxConfig)` sets IDs (`uuid.NewString()`) and defaults. |
| Driven adapter | Package `<name>strg`, type `Storage`, `New(Config) (*Storage, error)`. Persistence row models live in `model.go` with mappers both ways. In `cmd/`, import aliases get a tech prefix: `memgamestrg`, `sqlgamestrg`, `ddbgamestrg`. |
| REST adapter | `APIConfig` holds Service interfaces only. Handlers are named `serveXxx`; request bodies `xxxReqBody` with `Validate()`. Error codes are `ERR_UPPER_SNAKE`. Every response uses the `RespBody` envelope, with `conjson` producing snake_case keys from entities. |
| Comments | Short lowercase `// step` comments inside method bodies (`// get partner instance`). |
| Tests | Core is tested with **hand-written mocks** in the `_test` package. Adapters get integration tests against real infra via docker-compose. Use `testify/require`. |

Default libraries: `go-chi/chi/v5`, `go-chi/render`, `go-chi/cors`, `Rican7/conjson`, `gosidekick/goconfig`, `jmoiron/sqlx`, `aws-sdk-go`, `google/uuid`, `gopkg.in/validator.v2`, `stretchr/testify`.

## Workflows

Pick the workflow that matches the request. Load the referenced file when you reach that step.

### 1. Scaffold a new project
1. Follow the methodology: write a high-level overview in `README.md`, list the use cases, and write the API spec in `docs/api/`. Ask the user when requirements are unclear.
2. Group endpoints/use cases into **contexts**; each context becomes one `service/<ctx>` package.
3. Build in this order: **Service interface, Entities, Storage ports, service impl + tests, Adapters, `cmd/` wiring**. See [references/core.md](references/core.md), [references/driver-adapters.md](references/driver-adapters.md), [references/driven-adapters.md](references/driven-adapters.md).
4. Start with the memory storage adapter so the app runs with no infra, then add real storage.
5. Add `cmd/<deliverable>/{main,config,dep}.go`, Dockerfiles, docker-compose for run and tests, a Makefile, and CI. See [references/wiring-and-testing.md](references/wiring-and-testing.md).
6. Verify with `go build ./... && go vet ./... && go test ./internal/core/...`.

### 2. Add a feature (use case or new service)
1. Decide the context: a new method on an existing `Service`, or a new `service/<ctx>` package.
2. Add the method to the `Service` interface **with a doc comment** stating its behavior and the sentinel errors it returns.
3. Add or extend entities (behavior goes on the entity when it is about the entity's own state).
4. Add or extend the ports in `storage.go`, keeping only the methods this service needs.
5. Implement it, then add mock-based tests in `service_test.go` (extend the existing mocks).
6. Expose it in driver adapters (handler, request body, error mapping) and implement new port methods in **every** driven adapter (memory, mysql, dynamodb...).
7. Wire new dependencies in `cmd/*/dep.go` and `main.go`, then update `docs/api` and `docs/db`.

### 3. Add an adapter
- **Driven** (new storage tech or external service): create `internal/driven/storage/<tech>/<name>strg/` implementing the existing ports. Add the config section in `cmd/*/config.go`, a `case` in `initStorageDeps` in `dep.go`, an integration test, and the infra service in `deploy/local/tests/docker-compose.yml` and `deploy/local/run/<variant>/`. **Do not change Core** to fit the adapter; if a port truly must change, change it in Core first and update all adapters.
- **Driver** (REST, Lambda, CLI, queue consumer): create `internal/driver/<tech>/`, depend only on core `Service` interfaces, map service errors to the transport's error format, and add a new `cmd/<deliverable>/`.

See [references/driven-adapters.md](references/driven-adapters.md), [references/driver-adapters.md](references/driver-adapters.md).

### 4. Review for compliance
Run through [references/review-checklist.md](references/review-checklist.md). Report each violation with `file:line`, the rule broken, and a concrete fix. Distinguish hard violations (dependency direction, infra in Core) from convention drift.

## Pitfalls to avoid
- Driver error mapping must use `errors.Is`, not `switch err { case ...: }`, so wrapped errors still map correctly (see `handleServiceError` in `internal/driver/rest/api.go`).
- Always check `err` first, then `nil`. Checking `game == nil` before `err` makes a storage failure surface as `ErrGameNotFound`.
