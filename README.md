# Hirelly Executive Advisory — Go Backend Core Engine

A production-grade, cloud-native Go backend service designed for **Hirelly Executive Advisory** (`hirelly.in`). Built with modern Go patterns, PostgreSQL connection pooling (`pgx/v5`), Supabase REST fallback, and transactional email dispatching via Brevo.

Accelerated with **[Abeta-dev/go-libs](https://github.com/Abeta-dev/go-libs)** and **[Abeta-dev/go-app-kit](https://github.com/Abeta-dev/go-app-kit)**.

---

## 🏛️ Architecture & Clean Layout

```
hirelly-backend/
├── cmd/
│   └── server/
│       └── main.go                     # Bootstrap & graceful shutdown
├── internal/
│   ├── config/
│   │   └── config.go                   # .env & environment variable loader
│   ├── domain/
│   │   ├── application.go              # Mandate & inquiry models + validation
│   │   ├── subscriber.go               # Talent/recruiter subscription model
│   │   ├── job.go                      # Executive job posting model
│   │   └── errors.go                   # Standard domain errors
│   ├── repository/
│   │   ├── repository.go               # Persistence contracts
│   │   └── postgres/
│   │       ├── db.go                   # PostgreSQL pgxpool + Supabase REST fallback
│   │       ├── mandate_repo.go         # Mandate persistence
│   │       ├── subscriber_repo.go      # Subscriber persistence
│   │       └── job_repo.go             # Job listings
│   ├── service/
│   │   ├── advisory_service.go         # Mandate workflow & Brevo email dispatch
│   │   ├── subscriber_service.go       # Talent registration logic
│   │   └── email_service.go            # Brevo v3 integration (go-app-kit/notifications)
│   ├── handler/
│   │   ├── router.go                   # Gin router with Abeta-dev/go-libs middleware
│   │   ├── mandate_handler.go          # Mandates & advisory API endpoints
│   │   ├── subscriber_handler.go       # Waitlist & subscriber endpoints
│   │   ├── job_handler.go              # Executive opportunities endpoint
│   │   ├── health_handler.go           # Health & readiness probes
│   │   └── middleware.go               # Permissive CORS headers
├── migrations/
│   ├── 001_create_subscribers.sql      # Subscribers table schema
│   ├── 002_create_applications.sql     # Executive mandates table schema
│   └── 003_create_jobs.sql             # Executive jobs table schema
├── pkg/
│   └── uuid/                           # RFC 4122 v4 UUID generator
├── Dockerfile                          # Multi-stage production container build
├── docker-compose.yml                  # Local PostgreSQL + Backend orchestration
├── Makefile                            # Build, test, run, and docker commands
├── .env.example                        # Template environment variables
└── README.md                           # Documentation
```

---

## ⚡ Libraries & Accelerators Used

* **[`github.com/Abeta-dev/go-libs`](https://github.com/Abeta-dev/go-libs)**:
  * `ginmw.RequestID()`: Correlates every incoming HTTP request with `X-Request-ID`.
  * `ginmw.SecurityHeaders("hirelly-api")`: Defensive headers (CSP, HSTS, X-Frame-Options, X-Content-Type-Options).
  * `ginmw.Recovery()`: Panic trapping with structured stack trace logging & OpenTelemetry telemetry.
  * `ginmw.Logger()`: Structured `slog` access logging with latency, client IP, and route profiling.
  * `ginmw.LimitBodyDefault()`: Protects endpoints against oversized payload amplification attacks.
  * `health.New()`: Production health & readiness engine with pluggable database pingers.
* **[`github.com/Abeta-dev/go-app-kit`](https://github.com/Abeta-dev/go-app-kit)**:
  * `notifications.NewBrevoSender()`: Native Brevo API v3 transactional email dispatch with HTML templates and recipient routing.
* **[`github.com/jackc/pgx/v5`](https://github.com/jackc/pgx)**:
  * High-concurrency PostgreSQL connection pooling (`pgxpool.Pool`).

---

## 🔌 API Endpoints Reference

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | System health check (probes PostgreSQL, uptime, memory) |
| `GET` | `/ready` | Kubernetes / Cloud readiness probe |
| `POST` | `/api/v1/mandates` | Submit new executive search inquiry / advisory mandate |
| `GET` | `/api/v1/mandates` | Retrieve recent mandates with pagination (`?limit=20&offset=0`) |
| `GET` | `/api/v1/mandates/:id` | Retrieve single mandate by UUID |
| `POST` | `/api/v1/subscribers` | Register talent/client subscriber to platform waitlist |
| `GET` | `/api/v1/subscribers` | List recent subscribers |
| `GET` | `/api/v1/jobs` | Retrieve active executive opportunities |
| `POST` | `/api/applications` | **1:1 Drop-in compatibility** with existing frontend forms |
| `GET` | `/api/applications` | **1:1 Drop-in compatibility** for application records |
| `GET` | `/api/jobs` | **1:1 Drop-in compatibility** for job listings |

---

## 🚀 Quick Start

### 1. Prerequisites
* Go `1.24+` installed (`go version`)
* PostgreSQL or Supabase account
* Brevo API key

### 2. Configure Environment
```bash
cp .env.example .env
# Edit .env with your PostgreSQL credentials & Brevo API Key
```

### 3. Run Locally
```bash
make run
# Or directly:
go run ./cmd/server
```
The server will start at `http://localhost:8080`.

### 4. Run Unit Tests
```bash
make test
# Output:
# === RUN   TestHealthCheck
# --- PASS: TestHealthCheck (0.00s)
# === RUN   TestCreateMandate_Success
# --- PASS: TestCreateMandate_Success (0.00s)
# PASS
# ok      hirelly-backend/internal/handler
```

### 5. Build Binary
```bash
make build
# Binary compiled to bin/server
```

---

## 🐳 Docker Deployment

### Run Container Stack
```bash
docker compose up -d
```
Spins up:
1. PostgreSQL 16 Alpine on port `5432` with automated migration seeding.
2. `hirelly-backend` multi-stage container on port `8080`.

---

## 📜 Database Migrations

Apply the migration files in `migrations/` to your database:
```bash
psql $POSTGRES_URL -f migrations/001_create_subscribers.sql
psql $POSTGRES_URL -f migrations/002_create_applications.sql
psql $POSTGRES_URL -f migrations/003_create_jobs.sql
```

---

## 🛡️ License
MIT License © 2026 Hirelly International Private Limited.
