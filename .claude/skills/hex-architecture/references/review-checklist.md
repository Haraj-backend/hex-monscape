# Hex architecture review checklist

Report each finding as: `file:line`, the rule broken, why it matters, and a concrete fix. Group the findings into **Violations**, which break the architecture and must be fixed, and **Convention drift**, which is inconsistent with team style and should be fixed. Do not report something as a violation unless you verified it in the code.

Keep fixes simple. Propose only fixes for real problems, written so a junior developer can follow them. Do not suggest extra abstractions or clever idioms as improvements. Overengineered code is itself drift, so report it.

Set `MOD` to the module path from `go.mod` before running the commands.

```bash
MOD=$(go list -m)
```

## A. Dependency direction (violations)

1. **Core imports no adapters, cmd, or infra.**
   ```bash
   go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./internal/core/... \
     | grep -E "$MOD/(internal/driver|internal/driven|cmd)|net/http|database/sql|jmoiron/sqlx|aws-sdk-go|go-chi|go-sql-driver"
   ```
   Expect no output. Test-only imports (`.TestImports`) of `testutil` are fine.
2. **Driver adapters do not import driven adapters.**
   ```bash
   go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./internal/driver/... | grep "$MOD/internal/driven"
   ```
3. **Driven adapters do not import services.** They may import only `internal/core/entity` and their own `shared` package.
   ```bash
   go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./internal/driven/... | grep "$MOD/internal/core/service"
   ```
   Expect no output.
4. **Only `cmd/` constructs adapters and services.** Search for `strg.New(` or `NewAPI(` outside `cmd/` and tests.
5. **Ports are defined in Core.** Interfaces that services depend on live in `internal/core/service/<ctx>/storage.go` (or `service.go`), not in adapter packages. Services hold interface fields, never concrete adapter types.

## B. Core design (violations unless noted)

6. Service methods take and return primitives or entities, never transport or persistence types (`*http.Request`, req bodies, `XxxRow`, SDK types).
7. Entities have **no** `json:`, `db:`, or `dynamodbav:` tags.
   ```bash
   grep -rnE '`[^`]*(json|db|dynamodbav):' internal/core/entity --include='*.go' | grep -v _test.go
   ```
8. No application logic in adapters: no state transitions, business rules, or decisions in handlers or storage. Storage only maps and persists; handlers only decode, validate shape, call, and map errors.
9. Every `Service` method has a doc comment describing its behavior and the sentinel errors it returns. *(drift)*
10. Ports document their not-found contract, and adapters honor it by returning `nil, nil`. Check each `GetXxx` in `internal/driven/**`.
11. Services check `err` **before** checking for a `nil` result. Look for `if x == nil` that precedes `if err != nil` on the same call.

## C. Error handling

12. Expected outcomes are sentinel `ErrXxx` values returned unwrapped. Infra errors are wrapped with `%w` and context (`"unable to X due: %w"`). *(drift)*
13. Driver error mapping uses `errors.Is`/`errors.As`, not `==` or `switch err {`.
    ```bash
    grep -rnE 'switch err \{|err == [a-z]+\.Err' internal/driver
    ```
14. All service errors are mapped in one place per driver (`handleServiceError`). Inline mapping in a single handler is drift.
15. Every sentinel error has a mapped REST error and an `ERR_*` code documented in `docs/api/`.

## D. Construction and conventions (drift)

16. Every constructor takes `XxxConfig` with `validate` tags and calls `Validate()` first. Required deps are tagged `nonnil`/`nonzero`.
17. Naming: `Service`/`service`/`NewService`, `XxxStorage` ports, `<name>strg` packages with type `Storage` and `New(Config)`, `serveXxx` handlers, `xxxReqBody` request bodies, tech-prefixed aliases in `cmd/`.
18. Persistence row models live in the adapter's `model.go` with mappers both ways.

## E. Tests and docs

19. Each service has `TestNewService` (nil deps) plus a test per use case using hand-written mocks, not real storage. *(violation if a service test hits real infra)*
20. Each driven adapter has an integration test for the round trip and for not found, wired into `deploy/local/tests/docker-compose.yml`.
21. New ports or methods are implemented by **every** adapter of that port (memory, mysql, dynamodb). `go build ./...` catches this where `cmd/` assigns the adapter to the port. Do not suggest `var _ Port = (*Storage)(nil)` assertions.
22. `docs/api`, `docs/db`, and the README reflect the change.

## Finish

Run `go build ./... && go vet ./... && go test ./internal/core/...`. Adapter tests need `make test`.
