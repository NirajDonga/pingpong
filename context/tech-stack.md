# Tech Stack

## api

```txt
Language: Go
HTTP framework: Gin
Queue/broker: NATS
Relational DB: PostgreSQL
Time-series DB: ClickHouse
Auth: JWT or secure session cookies
```

## scheduler

```txt
Language: Go
Queue/broker: NATS
Database: PostgreSQL
```

## worker

```txt
Language: Go
HTTP client: net/http or Resty
Timing: httptrace if using net/http
Queue/broker: NATS
```

## web

```txt
Framework: Next.js
Language: TypeScript
Styling: Tailwind CSS
Charts: Recharts
Icons: lucide-react
```

## infra

```txt
PostgreSQL
ClickHouse
NATS
Docker Compose
```

