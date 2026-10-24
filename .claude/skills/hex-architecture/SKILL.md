---
name: hex-architecture
description: Build, refactor, and review Go applications using Haraj's Hexagonal Architecture (ports & adapters, without DDD), as implemented in hex-monscape. Use when scaffolding a new Go service/project, restructuring an existing Go codebase into hex, adding a use case, core service, entity, or storage port, adding a driver adapter (REST API, Lambda, CLI subcommand, cron scheduler, queue consumer) or a driven adapter (memory, files, MySQL, DynamoDB, S3, external API), wiring dependencies in cmd/, or reviewing code for hex compliance and maintainability (e.g. "does this follow hex?", "core imports infra", "rewrite this into hex", "is this too complicated?").
---

# Hexagonal Architecture (Haraj style)

Source of truth: [docs/reference/hex-architecture.md](../../../docs/reference/hex-architecture.md) and [docs/reference/project-methodology.md](../../../docs/reference/project-methodology.md). This repo (`hex-monscape`) is the canonical example, so read the real files cited below whenever a template is not enough.

## The goal: maintainers can understand the software

Hex exists so a maintainer can open the repo and quickly answer "what does this app do, and where do I change X?". Core tells you *what* the app does. Adapters tell you *how* it talks to the outside world. Every rule below serves that goal.

So apply the conventions with common sense, not literally. If following a convention makes the code harder to read in a particular codebase, for example a package per entity in a small app, choose the readable option and keep the hard rules (dependency direction, ports in Core). Before finishing any change, reread the code as a new maintainer would: are the names honest, is anything redundant, and can each step be followed without jumping around?

## Principles

1. Clearly divide the `inside` (Core) from the `outside` (Actors and Adapters) of the application.
2. Dependencies always point **inward**. Core imports nothing from adapters, `cmd/`, or infrastructure libraries.
3. Isolate the boundary with **ports** (interfaces defined in Core) and **adapters** (implementations outside Core).
4. **Keep it simple and avoid overengineering.** Write code a junior developer can read and maintain. Do not add clever Go idioms, extra layers or abstractions, generic helpers, functional options, or "nice to have" safety tricks such as `var _ Port = (*Storage)(nil)` unless the code actually needs them.
5. **Size packages for navigation, not for purity.** A package should hold a meaningful chunk of code. Many one-file packages and aliased imports are a smell (see [Package granularity](#package-granularity)).

The four pillars:

| Pillar | What it is | In this repo |
|---|---|---|
| **Core** | The application logic: how the app behaves for each use case. Made of `Service`, `Entity`, and `Storage` interfaces. | `internal/core/` |
| **Actors** | Driver actors start an interaction (an HTTP client, a cron tick). Driven actors are called by the app (MySQL, GitHub API). | – |
| **Ports** | Interfaces inside Core. **Driver port** = `Service` interface. **Driven port** = `XxxStorage` (or `XxxClient`) interfaces. | `service/<ctx>/service.go`, `service/<ctx>/storage.go` |
| **Adapters** | Driver adapters translate a driver actor into Service calls. Driven adapters implement driven ports. | `internal/driver/rest`, `internal/driven/storage/...` |

Use the term **application logic** (not "business logic"). The team applies hex **without DDD**: there are no aggregates, value objects, or domain events. A service may update several entities directly when the app is small.

## Project layout

```
cmd/<deliverable>/              composition root: main.go, config.go, dep.go (only place that wires adapters)
internal/core/entity/           flat package of entities (plain structs + behavior methods + pure rules)
internal/core/service/<ctx>/    one package per context: service.go (driver port + impl), storage.go (driven ports)
internal/core/testutil/         fixture generators shared by tests
internal/driver/<tech>/         driver adapters, e.g. rest/{api,error,req_body,resp_body}.go, scheduler/
internal/driven/storage/...     storage adapters (see Package granularity)
internal/driven/<service>/      external service adapters, e.g. github/, payment/
build/package/<deliverable>/Dockerfile
deploy/local/{run/<variant>,tests,shared}/   docker-compose for running & integration tests
docs/{api,db,reference}/        API spec, DB schemas, architecture notes
Makefile                        run, run-<variant>, test, test-build-all
```

### Package granularity

Group driven adapters by **kind**, then split only when there is a real reason:

- **Default: one package per kind.** For example, a single `internal/driven/storage` package with one file per concern (`result.go`, `progress.go`, `catalog.go`) and a shared `model.go` for row models. `cmd/` then imports one package, without aliases.
- **Split by tech** (`storage/memory`, `storage/mysql`, `storage/dynamodb`) when there are **alternative implementations of the same ports** that `cmd/` switches between. This is why hex-monscape has per-tech packages: the server runs on memory, MySQL, or DynamoDB.
- **Split by entity** (`<tech>/<name>strg/`) only when a tech package gets big enough that one package is hard to navigate. hex-monscape does this for its three storage techs. Do not copy it into an app with one storage tech and a handful of files.

Core services stay **one package per context**, because each context owns a different set of ports. Merge contexts only when they share almost all of their ports.

## Conventions

| Thing | Convention |
|---|---|
| Constructors | Every constructor (entity, service, adapter, API) takes a `XxxConfig` struct with `validate` tags plus a `Validate()` method using `gopkg.in/validator.v2` (`nonzero`, `nonnil`, `min=1`). The constructor calls `cfg.Validate()` first. Pure functions that cannot fail (e.g. `entity.NewResult(repo, metrics, now)`) don't need a config. |
| Service | Interface `Service` (driver port) with **doc comments describing behavior and returned errors** for each method. Unexported impl `service`. `NewService(ServiceConfig) (Service, error)`. One package per context (`play`, `battle`, `crawl`). |
| Driven ports | `XxxStorage` / `XxxClient` interfaces in `service/<ctx>/storage.go`, owned by the consuming service. Duplicating a port across services (e.g. `play.GameStorage` and `battle.GameStorage`) is intended. |
| Not found | Storage and external clients return `nil, nil` (or an empty list). The service converts that into a sentinel error or a skip. Adapter-internal "not found" errors never leave the adapter. |
| Errors | Sentinel `var ErrXxx = errors.New("...")` in the service package. Infra errors are wrapped as `fmt.Errorf("unable to <do X> due: %w", err)`. Entity errors (e.g. `entity.ErrInvalidState`) are mapped to service errors. |
| Entities | Exported fields, **no JSON/DB tags**. Typed string enums (`type State string`). Behavior and pure rules (filters, scoring, metrics) live on entities or as functions in `entity`. `NewXxx(XxxConfig)` sets IDs (`uuid.NewString()`) and defaults. |
| Driven adapter | Type named after what it stores (`ResultStorage`, `GameStorage`) with `NewXxx(XxxConfig)`. In per-entity packages the legacy form `<name>strg.Storage` + `New(Config)` is fine. Persistence row models live in `model.go` with mappers both ways. |
| REST adapter | `APIConfig` holds Service interfaces only. Handlers are named `serveXxx`; request bodies `xxxReqBody` with `Validate()`. Error codes are `ERR_UPPER_SNAKE`. Every response uses the `RespBody` envelope, with `conjson` producing snake_case keys from entities. |
| Code layout | Separate logical steps with a blank line, e.g. after an `if err != nil { ... }` block and before the next step. Keep together what belongs together: the `nil` check of a call stays right under its `err` check, and `defer` stays right under the resource it closes. Use short lowercase `// step` comments for non-obvious steps (`// get partner instance`). |
| Names | Names must say what a thing *is*. When scope grows (e.g. crawl "topics" now also include company orgs), rename (`CurrentTopic` → `CurrentUnit`). Keep old wire or file keys in the row model when renaming would break existing data. |
| Tests | Core is tested with **hand-written mocks** in the `_test` package. Adapters get integration tests against real infra via docker-compose, or plain unit tests for files and fakes for SDK clients. Use `testify/require`. Time and sleep can be injected through optional `Now` / `Sleep` config fields. |

Default libraries: `go-chi/chi/v5`, `go-chi/render`, `go-chi/cors`, `Rican7/conjson`, `gosidekick/goconfig`, `jmoiron/sqlx`, `aws-sdk-go`, `google/uuid`, `gopkg.in/validator.v2`, `stretchr/testify`, `go-co-op/gocron/v2` (scheduling).

## Workflows

Pick the workflow that matches the request. Load the referenced file when you reach that step. Every workflow ends with a **readability pass**: reread what you wrote as a new maintainer would, and fix misleading names, redundant state, dense blocks, and needless packages.

### 1. Scaffold a new project
1. Follow the methodology: write a high-level overview in `README.md`, list the use cases, and write the API spec in `docs/api/`. Ask the user when requirements are unclear.
2. Group endpoints/use cases into **contexts**; each context becomes one `service/<ctx>` package.
3. Build in this order: **Service interface, Entities, Storage ports, service impl + tests, Adapters, `cmd/` wiring**. See [references/core.md](references/core.md), [references/driver-adapters.md](references/driver-adapters.md), [references/driven-adapters.md](references/driven-adapters.md).
4. Start with the simplest storage that runs without infra (memory, or files for batch jobs), then add real storage.
5. Add `cmd/<deliverable>/{main,config,dep}.go`, Dockerfiles, docker-compose for run and tests, a Makefile, and CI. See [references/wiring-and-testing.md](references/wiring-and-testing.md).
6. Verify with `go build ./... && go vet ./... && go test ./...`.

### 2. Add a feature (use case or new service)
1. Decide the context: a new method on an existing `Service`, or a new `service/<ctx>` package.
2. Add the method to the `Service` interface **with a doc comment** stating its behavior and the sentinel errors it returns.
3. Add or extend entities (behavior goes on the entity when it is about the entity's own state).
4. Add or extend the ports in `storage.go`, keeping only the methods this service needs.
5. Implement it, then add mock-based tests in `service_test.go` (extend the existing mocks).
6. Expose it in driver adapters (handler, request body, error mapping) and implement new port methods in **every** driven adapter of that port.
7. Wire new dependencies in `cmd/*/dep.go` and `main.go`, then update `docs/api` and `docs/db`.

### 3. Add an adapter
- **Driven** (new storage tech or external service): add it to the package for its kind (see [Package granularity](#package-granularity)) and implement the existing ports. Add its config to `cmd/*/config.go`, construct it in `dep.go`, add a test, and add any infra to `deploy/local/`. **Do not change Core** to fit the adapter. If a port truly must change, change it in Core first and update all adapters.
- **Driver** (REST, Lambda, CLI, cron scheduler, queue consumer): create `internal/driver/<tech>/`, depend only on core `Service` interfaces, and map service errors to the transport's semantics. Add a new `cmd/<deliverable>/`, or a subcommand when one binary serves several drivers.

See [references/driven-adapters.md](references/driven-adapters.md), [references/driver-adapters.md](references/driver-adapters.md).

### 4. Restructure an existing codebase into hex
1. Read the current code first and list the use cases, the infra it talks to, and the deliverables that actually run in production.
2. Move pure logic (filters, scoring, state changes) into `internal/core/entity`, and orchestration into `service/<ctx>`. Strip json/db tags from what becomes an entity, and move them into row models in the adapters.
3. Keep behavior and file/wire formats identical. Verify by running old and new side by side on the same input and diffing the outputs.
4. Collapse binaries that only exist to call one pipeline into subcommands of one deliverable when that is simpler to deploy and understand.
5. Do the readability pass, then commit in reviewable steps (core, adapters + wiring, deploy + docs).

### 5. Review for compliance
Run through [references/review-checklist.md](references/review-checklist.md). Report each violation with `file:line`, the rule broken, and a concrete fix. Distinguish hard violations (dependency direction, infra in Core) from convention drift and readability problems.

## Pitfalls to avoid
- Driver error mapping must use `errors.Is`, not `switch err { case ...: }`, so wrapped errors still map correctly (see `handleServiceError` in `internal/driver/rest/api.go`).
- Always check `err` first, then `nil`. Checking `game == nil` before `err` makes a storage failure surface as `ErrGameNotFound`.
- Copying hex-monscape's per-tech, per-entity storage packages into an app with one storage tech. You end up with many one-file packages and a `dep.go` full of import aliases.
- Copying style mechanically, e.g. writing no blank lines between steps because a reference file was dense. Readability beats imitation.
- Keeping state that can be derived from other state (e.g. a second "seen IDs" set that is always a subset of the first). It makes maintainers reason about a difference that does not exist.
- Letting an empty result overwrite good data. When a use case publishes or replaces data (export, sync), return a sentinel such as `ErrNoResults` instead of writing an empty file.
