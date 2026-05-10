# Go gRPC Guestbook

A Go service that exposes a **Guestbook** API over gRPC and HTTP (grpc-gateway), with **JWT authentication** for protected RPCs, **MySQL** persistence via GORM, and optional **Redis** caching for list queries. Protobuf definitions drive code generation and **OpenAPI/Swagger** docs.

## Project structure

| Path | Purpose |
| ---- | ------- |
| `cmd/server` | Application entrypoint: config, DB, DI wiring, gRPC + HTTP servers |
| `proto` | `.proto` sources and OpenAPI/grpc-gateway options |
| `gen` | Generated Go, gRPC, gateway, and Swagger JSON (created by `make generate`) |
| `internal/controller` | HTTP and gRPC adapters (handlers only; no business rules) |
| `internal/service` | Use cases: guestbook operations, auth, validation contracts |
| `internal/repository` | Persistence **ports** (interfaces) and **GORM adapters** |
| `internal/model` | Domain entities mapped to the database |
| `internal/di` | Composition root: wires repository → service → gRPC controller |
| `internal/dto` | JSON request/response shapes for HTTP auth endpoints |
| `internal/mapper` | Conversion between domain models and protobuf messages |
| `pkg/framework` | gRPC server wrapper (interceptors, lifecycle) |
| `pkg/config` | JSON configuration loading (`config.local.json` / `config.prod.json`) |
| `pkg/auth` | JWT issue/validate and auth interceptor |
| `pkg/cache` | Cache abstraction and Redis implementation |
| `pkg/router` | HTTP mux: auth routes, gateway under `/v1/`, Swagger in non-production |
| `test` | Integration tests against a real gRPC stack |
| `swagger` | Copy of OpenAPI JSON used by `make swagger` static server |

Layers follow **SOLID**-friendly boundaries: dependencies point inward (handlers → services → repository interfaces), with construction centralized in `internal/di`.

## Prerequisites

- **Go 1.24+** (see `go.mod`)
- **Make** (recommended)
- **`protoc`** — e.g. macOS: `brew install protobuf`
- **MySQL** — database for the app (DSN in `pkg/config/config.local.json`)
- **Redis** (optional) — set `redis.enabled` in config; list caching is skipped if disabled or unreachable
- **Python 3** — used by `make swagger` to serve static Swagger UI

## Configuration

- Default config file: `pkg/config/config.local.json`
- Production: set `APP_ENV=production` to load `pkg/config/config.prod.json`
- Environment variables with `APP_` prefix override file values (Viper), for example:

```bash
APP_SERVER_PORT=60000 APP_SERVER_HTTPPORT=9090 go run cmd/server/main.go
APP_DATABASE_HOST=127.0.0.1 APP_DATABASE_PASSWORD=secret go run cmd/server/main.go
APP_JWT_SECRETKEY=my-prod-secret APP_JWT_TOKENDURATION=12 go run cmd/server/main.go
```

- Adjust `server.port` (gRPC), `server.httpPort` (HTTP gateway + Swagger), `database`, `redis`, and `jwt` to match your environment.

### Runtime behavior at startup

- Loads config from `config.local.json` or `config.prod.json` depending on `APP_ENV`
- Connects to MySQL using GORM and runs `AutoMigrate` for `GuestbookEntry` and `User`
- Tries to connect Redis only when `redis.enabled=true`; if Redis is unreachable, app still runs with cache disabled
- Applies JWT interceptor only to `GuestbookService/AddEntry` (list is public)

## Code generation

The `Makefile` installs Go plugins (`protoc-gen-go`, `protoc-gen-go-grpc`, `protoc-gen-grpc-gateway`, `protoc-gen-openapiv2`) and fetches proto dependencies.

```bash
make generate
```

This installs tools as needed, generates Go and gateway code under `gen/`, and emits OpenAPI JSON (e.g. `gen/guestbook/v1/guestbook.swagger.json`).

## Swagger / OpenAPI

After `make generate`:

1. **With the full server** (`make run`): open **`http://localhost:8080`** or **`http://localhost:8080/swagger-ui.html`** (port from `httpPort` in config). Raw spec: `/swagger/guestbook.swagger.json` when the server exposes it (non-production).
2. **Static UI only**: `make swagger` serves on **`http://localhost:8081/swagger-ui.html`**.

## Running the server

Ensure MySQL is running and the database exists (see config). Then:

```bash
go run cmd/server/main.go
```

Or:

```bash
make build
./bin/server
```

- **gRPC**: listens on the port in config (default `50051`)
- **HTTP**: gateway and auth routes on `httpPort` (default `8080`)

## Make targets

- `make generate`: install protoc plugins, fetch proto deps, generate gRPC/gateway/OpenAPI code into `gen/`
- `make build`: build binary to `bin/server`
- `make run`: run server via `go run cmd/server/main.go`
- `make test`: run `go test ./...`
- `make swagger`: copy generated spec to `swagger/` and serve static UI on port `8081`
- `make clean`: remove `gen/`, `bin/`, and `swagger/`

## Testing

Run the full module test suite:

```bash
make test
# or
go test ./...
```

Focused packages:

```bash
go test ./internal/...
go test ./test/...
```

## API overview

Defined in `proto/guestbook/v1/guestbook.proto`:

### GuestbookService (gRPC + HTTP via gateway)

- **`AddEntry`** — create an entry (JWT required; protected method)
- **`ListEntries`** — paginated list (query params: `limit`, `offset`; defaults `10/0`, max `limit=100`), optional Redis cache

### Authentication

- Proto also defines **`AuthService`** (`Login`, `Register`) for documentation alignment with OpenAPI.
- The running server serves **`POST /auth/login`** and **`POST /auth/register`** as JSON via `internal/controller` and `internal/dto`, backed by `internal/service` and MySQL users.

### HTTP endpoint summary

- `POST /auth/register` (public, JSON body: `username`, `email`, `password`)
- `POST /auth/login` (public, JSON body: `username`, `password`)
- `GET /v1/guestbook` (public, supports `limit` and `offset`)
- `POST /v1/guestbook` (requires `Authorization: Bearer <token>`)

### Example flow (HTTP)

Register:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","email":"alice@example.com","password":"secret123"}'
```

List entries:

```bash
curl 'http://localhost:8080/v1/guestbook?limit=10&offset=0'
```

Add entry (replace `<token>`):

```bash
curl -X POST http://localhost:8080/v1/guestbook \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <token>' \
  -d '{"name":"Alice","email":"alice@example.com","message":"Hello from HTTP"}'
```

Obtain a JWT from login or register, then call protected HTTP/gateway methods with header: `Authorization: Bearer <token>`.
