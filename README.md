# Go gRPC Guestbook

A Go service that exposes a **Guestbook** API over gRPC and HTTP (grpc-gateway), with **JWT authentication** for protected RPCs, **MySQL** persistence via GORM, and optional **Redis** caching for list queries. Protobuf definitions drive code generation and **OpenAPI/Swagger** docs.

## Project structure

| Path | Purpose |
|------|---------|
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
- Adjust `server.port` (gRPC), `server.httpPort` (HTTP gateway + Swagger), `database`, `redis`, and `jwt` to match your environment.

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

**GuestbookService (gRPC + HTTP via gateway)**

- **`AddEntry`** — create an entry (JWT required; see interceptor in `cmd/server/main.go`)
- **`ListEntries`** — paginated list of entries (optional Redis cache)

**Authentication**

- Proto also defines **`AuthService`** (`Login`, `Register`) for documentation alignment with OpenAPI.
- The running server serves **`POST /auth/login`** and **`POST /auth/register`** as JSON via `internal/controller` and `internal/dto`, backed by `internal/service` and MySQL users.

Obtain a JWT from login or register, then call protected HTTP/gateway methods with header: `Authorization: Bearer <token>`.
