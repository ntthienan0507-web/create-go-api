# create-go-api

Interactive CLI tool that scaffolds production-ready Go API projects from [gostack-kit](https://github.com/ntthienan0507-web/gostack-kit).

Pick only the stacks you need — the tool clones the latest `gostack-kit` source, filters files based on your selection, and generates a clean project with the correct module path.

## Install

```bash
go install github.com/ntthienan0507-web/create-go-api/cmd/create-go-api@latest
```

## Quick Start

```bash
create-go-api
```

The interactive TUI walks you through:

```
┌──────────────────────────────────────────────┐
│  Go module path: github.com/myorg/my-api     │
│  Server port: 8080                           │
├──────────────────────────────────────────────┤
│  Database driver (pick one):                 │
│  > GORM — ORM, recommended                  │
│    SQLC — raw SQL, type-safe                 │
│    Both — switchable via config              │
├──────────────────────────────────────────────┤
│  Auth provider:                              │
│  > JWT (local HMAC tokens)                   │
│    Keycloak (OIDC)                           │
│    Both                                      │
├──────────────────────────────────────────────┤
│  Include sample User module? Yes             │
│  Include Swagger/OpenAPI?    No              │
├──────────────────────────────────────────────┤
│  Infrastructure (space to select):           │
│  [x] Redis — cache, sessions, pub/sub       │
│  [x] Kafka — event streaming, outbox         │
│  [ ] Encryption — AES-256 for PII            │
├──────────────────────────────────────────────┤
│  Features:                                   │
│  [ ] Cron — scheduled background jobs        │
│  [x] WebSocket — real-time communication     │
│  [ ] OpenTelemetry — distributed tracing     │
├──────────────────────────────────────────────┤
│  External services:                          │
│  [x] SendGrid — email                       │
│  [ ] Stripe — payments                      │
│  [ ] IceWarp — mail server (XML)            │
│  [x] Firebase — push notifications          │
│  [ ] Elasticsearch — search                 │
├──────────────────────────────────────────────┤
│  Kubernetes manifests? No                    │
│  SonarQube config?     No                    │
│  CI/CD: GitLab CI                            │
└──────────────────────────────────────────────┘

  Cloning gostack-kit...
  Selecting files based on config...
  Created 87 files
  go.mod → github.com/myorg/my-api
  go mod tidy
  git init

  Done! Next steps:
    cd my-api
    docker compose up -d
    go run . serve
```

## How It Works

```
create-go-api
     │
     ├─ Interactive prompts (charmbracelet/huh)
     │
     ├─ git clone --depth=1 gostack-kit
     │     └─ https://github.com/ntthienan0507-web/gostack-kit
     │
     ├─ Filter files based on your selections
     │     ├─ DB: keep gorm.go OR postgres.go+store.go (not both unless "both" selected)
     │     ├─ Auth: keep jwt.go OR keycloak.go (not both unless "both")
     │     ├─ Optional: broker/, ws/, cron/, crypto/, tracing/, external/*
     │     └─ DevOps: deployments/k8s/, .gitlab-ci.yml, .github/
     │
     ├─ Replace module path
     │     gostack-kit → github.com/myorg/my-api
     │
     ├─ go mod tidy
     └─ git init
```

**Always uses the latest source** — when `gostack-kit` gets a bug fix or new feature, the next `create-go-api` run automatically includes it. No template sync needed.

## Available Stacks

### Core (always included)

| Package | What |
|---------|------|
| `pkg/app` | DI container, graceful shutdown, readiness probe |
| `pkg/config` | Viper config from `.env` + env vars |
| `pkg/auth` | JWT / Keycloak (based on selection) |
| `pkg/database` | PostgreSQL / GORM / MongoDB (based on selection) |
| `pkg/response` | Generic typed JSON responses `Response[T]` |
| `pkg/apperror` | Structured errors with i18n-ready keys |
| `pkg/middleware` | Recovery, CORS, logging, auth, validation, response audit |
| `pkg/async` | Worker pool, parallel execution, context safety |
| `pkg/logger` | Zap structured logging |

### Selectable

| Stack | Flag in prompt | Packages included |
|-------|---------------|-------------------|
| **Redis** | Infrastructure → Redis | `pkg/database/redis.go`, `pkg/cache` |
| **Kafka** | Infrastructure → Kafka | `pkg/broker` (producer, consumer, dispatcher, batcher, outbox, relay, idempotency, topics) |
| **Encryption** | Infrastructure → Encryption | `pkg/crypto` (AES-256-GCM, bcrypt, random) |
| **Cron** | Features → Cron | `pkg/cron`, `cmd/cron.go` |
| **WebSocket** | Features → WebSocket | `pkg/ws` (hub, client, rooms, message routing) |
| **OpenTelemetry** | Features → OpenTelemetry | `pkg/tracing` |
| **SendGrid** | External → SendGrid | `pkg/external/sendgrid` + `pkg/httpclient` + `pkg/retry` |
| **Stripe** | External → Stripe | `pkg/external/stripe` + `pkg/httpclient` + `pkg/retry` |
| **IceWarp** | External → IceWarp | `pkg/external/icewarp` (XML codec) + `pkg/httpclient` |
| **Firebase** | External → Firebase | `pkg/external/firebase` (FCM push, auth verify) |
| **Elasticsearch** | External → Elasticsearch | `pkg/external/elasticsearch` + `pkg/httpclient` |
| **Kubernetes** | DevOps → K8s | `deployments/k8s/` (Kustomize base + dev/staging/prod overlays) |
| **SonarQube** | DevOps → SonarQube | `sonar-project.properties`, `scripts/setup-sonar.sh` |
| **GitHub Actions** | DevOps → CI/CD | `.github/workflows/ci.yml` |
| **GitLab CI** | DevOps → CI/CD | `.gitlab-ci.yml` |

### Auto-included dependencies

These are included automatically when needed — you don't select them:

| Package | Included when |
|---------|--------------|
| `pkg/httpclient` | Any external service (SendGrid, Stripe, IceWarp, Elasticsearch) |
| `pkg/retry` | Any external service or Kafka |
| `pkg/circuitbreaker` | Any external service or Kafka |
| `pkg/app/services.go` | Any external service, Redis, Kafka, or Encryption |

## Generated Project Structure

```
my-api/
├── main.go
├── cmd/                         # CLI commands (serve, migrate, db, cron)
├── modules/                     # Business modules
│   └── user/                    # Sample module (if selected)
├── pkg/                         # Shared infrastructure
│   ├── app/                     # DI + lifecycle
│   ├── apperror/                # Structured errors
│   ├── async/                   # Worker pool
│   ├── auth/                    # JWT / Keycloak
│   ├── broker/                  # Kafka (if selected)
│   ├── cache/                   # Redis cache (if selected)
│   ├── circuitbreaker/          # Circuit breaker (auto)
│   ├── config/                  # Configuration
│   ├── cron/                    # Scheduler (if selected)
│   ├── crypto/                  # Encryption (if selected)
│   ├── database/                # DB connections
│   ├── external/                # 3rd party clients (if selected)
│   ├── httpclient/              # HTTP client (auto)
│   ├── logger/                  # Zap logging
│   ├── middleware/              # HTTP middleware
│   ├── response/                # Response helpers
│   ├── retry/                   # Retry logic (auto)
│   ├── tracing/                 # OpenTelemetry (if selected)
│   └── ws/                      # WebSocket (if selected)
├── db/migrations/               # SQL migrations
├── deployments/k8s/             # Kubernetes (if selected)
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .env.example
```

## Source Template

This tool generates projects from [**gostack-kit**](https://github.com/ntthienan0507-web/gostack-kit) — a production-ready Go API template with:

- Clean architecture (Controller → Service → Repository)
- Pluggable database (SQLC / GORM / MongoDB)
- Pluggable auth (JWT / Keycloak)
- Kafka with Transactional Outbox Pattern
- Key-sharded parallel consumers
- Structured errors with i18n-ready keys
- Response audit middleware (catches hardcoded `ctx.JSON`)
- Request validation middleware (no more EOF)
- WebSocket with rooms and message routing
- AES-256 field-level encryption
- Graceful shutdown with connection draining
- Kubernetes deployment manifests (Kustomize)
- CI/CD pipelines (GitHub Actions / GitLab CI)

See the [gostack-kit README](https://github.com/ntthienan0507-web/gostack-kit/blob/main/README.md) and [ARCHITECTURE.md](https://github.com/ntthienan0507-web/gostack-kit/blob/main/ARCHITECTURE.md) for full documentation.

## Development

```bash
# Build
go build -o create-go-api ./cmd/create-go-api

# Test (clones gostack-kit — requires network)
go test ./...

# Install locally
go install ./cmd/create-go-api
```

## License

MIT
