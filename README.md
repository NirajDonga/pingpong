# PingPong 

A lightweight uptime monitoring platform. Add a URL, pick a check interval, and PingPong continuously pings it — tracking response times, detecting downtime, and logging incidents automatically.

## How It Works

```
web → api → PostgreSQL
              ↑
scheduler ────┘ (reads due monitors, publishes jobs)
              ↓
            NATS
              ↓
           worker (performs HTTP check, publishes result)
              ↓
            NATS
              ↓
            api → ClickHouse (stores result history)
                → PostgreSQL (updates status, opens/closes incidents)
```

1. User creates a monitor with a URL and interval
2. Scheduler finds due monitors every second and publishes check jobs to NATS
3. Worker receives the job, performs the HTTP request, and publishes the result
4. API consumes the result, stores it in ClickHouse, updates monitor status, and manages incidents

## Services

| Service | Language | Role |
|---|---|---|
| `api` | Go (Gin) | REST API, auth, result processing, incident management |
| `scheduler` | Go | Finds due monitors, publishes check jobs to NATS |
| `worker` | Go | Performs HTTP checks, measures response time |
| `frontend` | Next.js + TypeScript | User dashboard |

## Tech Stack

- **Databases** — PostgreSQL (app data & state) · ClickHouse (time-series check history)
- **Messaging** — NATS
- **Auth** — Session cookies / JWT
- **Frontend** — Next.js · Tailwind CSS · TanStack Query · Recharts
- **Infra** — Docker Compose

## Getting Started

**Start infrastructure:**

```bash
cd infra
docker compose up -d
```

This starts PostgreSQL (`:5432`), ClickHouse (`:8123`), and NATS (`:4222`).

**Run each service** (copy `.env.example` → `.env` in each folder first):

```bash
# API
cd backend/api && go run ./cmd/api

# Scheduler
cd backend/scheduler && go run ./cmd/scheduler

# Worker
cd backend/worker && go run ./cmd/worker

# Frontend
cd frontend && pnpm install && pnpm dev
```

## Project Structure

```
pingpong/
├── backend/
│   ├── api/          # REST API + result consumer
│   ├── scheduler/    # Check job publisher
│   └── worker/       # HTTP checker
├── data/
│   └── tinybird/     # Analytics (future)
├── docs/         # Documentation
├── frontend/     # Next.js dashboard
└── infra/        # Docker Compose
```

## Data Model (PostgreSQL)

- `users` — accounts
- `monitors` — URL, interval, timeout, expected status, current state
- `incidents` — auto-opened on failure, auto-closed on recovery

Check results (response time, status code, DNS/TCP/TLS/TTFB timings) are stored in ClickHouse.
