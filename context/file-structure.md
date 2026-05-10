# File Structure

Target file structure:

```txt
pingpong/
├── api/
│   ├── go.mod
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   └── internal/
│       ├── auth/
│       ├── config/
│       ├── db/
│       │   ├── postgres.go
│       │   └── clickhouse.go
│       ├── migrations/
│       │   ├── postgres/
│       │   │   └── 001_initial.sql
│       │   └── clickhouse/
│       │       └── 001_check_results.sql
│       ├── monitor/
│       ├── incident/
│       ├── nats/
│       └── shared/
│
├── scheduler/
│   ├── go.mod
│   ├── cmd/
│   │   └── scheduler/
│   │       └── main.go
│   └── internal/
│       ├── config/
│       ├── db/
│       ├── nats/
│       └── scheduler/
│
├── worker/
│   ├── go.mod
│   ├── cmd/
│   │   └── worker/
│   │       └── main.go
│   └── internal/
│       ├── checker/
│       ├── config/
│       ├── nats/
│       └── worker/
│
├── web/
│   ├── app/
│   ├── components/
│   │   ├── monitor/
│   │   ├── incident/
│   │   └── ui/
│   ├── lib/
│   │   └── api.ts
│   └── types/
│       ├── monitor.ts
│       ├── result.ts
│       └── incident.ts
│
├── infra/
│   ├── docker-compose.yml
│   ├── postgres/
│   │   └── init.sql
│   ├── clickhouse/
│   │   └── init.sql
│   └── nats/
│       └── nats.conf
│
├── README.md
└── context/
    ├── README.md
    ├── overview.md
    ├── architecture.md
    ├── services.md
    ├── tech-stack.md
    ├── requirements.md
    ├── build-phases.md
    ├── file-structure.md
    ├── coding-style.md
    └── later-scope.md
```

Schema ownership:

```txt
api owns all database schemas and migrations.
api uses PostgreSQL and ClickHouse.
scheduler uses PostgreSQL scheduling tables.
worker does not access PostgreSQL or ClickHouse.
infra runs dependencies, but does not own app schema.
```

