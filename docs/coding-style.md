# Coding Style

The backend style should stay close to the current Go code: simple, direct, readable, and not over-abstracted.

## Main Rule

Write boring code that is easy to follow.

Do not turn simple logic into a framework.
Do not add abstractions before they are useful.
Do not create clever helper layers just because they look clean.

Good code here should feel like:

```txt
load config
connect dependency
register handler
validate input
call service/function
return response
```

## Backend Style

Backend code is Go.

Follow the current backend style:

- small packages
- clear function names
- explicit error handling
- early returns
- simple structs for JSON messages
- config loaded from environment
- services communicate using typed structs
- logs use clear messages
- no unnecessary interfaces
- no generic framework-style architecture

Prefer this:

```go
target := c.Query("target")
if target == "" {
    c.String(http.StatusBadRequest, "Bad Request: 'target' parameter is required")
    return
}
```

Over this:

```go
validator := NewRequestValidator(...)
result := validator.Validate(...)
responseMapper.Write(...)
```

Unless that abstraction is actually needed.

## Package Style

Each package should own one clear thing.

Examples:

```txt
config      loads environment config
monitor     monitor model, routes, service, repository
scheduler   finds due monitors and publishes jobs
worker      consumes jobs and performs checks
incident    incident state logic
nats        NATS client/publisher/consumer helpers
db          database connections
shared      message structs shared inside one service
```

Avoid dumping unrelated code into one package.

## Function Style

Functions should be short enough to understand without scrolling too much.

Use early returns:

```go
if err != nil {
    log.Printf("failed to publish check job: %v", err)
    return
}
```

Keep the happy path easy to see.

Good split:

```txt
handler validates request
service applies business logic
repository talks to database
publisher sends NATS message
```

## Error Handling

Always handle errors explicitly.

Use `log.Fatalf` only during startup when the service cannot run.

Use `log.Printf` and return/continue for runtime errors.

Do not ignore errors from important operations.

## Config Style

Config should be loaded once at service startup.

Use explicit config structs.

Required environment variables should fail fast.

Config package should only load config. It should not connect to databases or run business logic.

## Message Style

NATS messages should use typed structs.

Do not pass around unstructured maps when the message shape is known.

Example:

```go
type CheckJob struct {
    MonitorID      string `json:"monitorId"`
    TargetURL      string `json:"targetUrl"`
    TimeoutSeconds int    `json:"timeoutSeconds"`
    ExpectedStatus int    `json:"expectedStatus"`
}
```

JSON field names should be stable because they are the contract between services.

## Service Boundaries

Keep service responsibilities strict.

```txt
api = validates user requests, owns schemas, stores monitors, processes results
scheduler = decides when checks should happen
worker = performs HTTP checks
web = dashboard
infra = local dependency setup
```

Rules:

- `worker` does not access PostgreSQL.
- `worker` does not access ClickHouse.
- `scheduler` does not perform HTTP checks.
- `web` does not talk directly to databases.
- database schemas and migrations live under `api`.

## Database Style

Keep SQL/database code out of HTTP handlers.

Use repository/store files for database operations.

Example shape:

```txt
monitor/routes.go      HTTP handlers
monitor/service.go     business logic
monitor/repository.go  PostgreSQL queries
```

ClickHouse is for time-series check results.
PostgreSQL is for user data, monitors, incidents, current status, and scheduler timing.

## Frontend Style

Frontend style is not established yet.

For now, keep frontend code simple:

- clear component names
- small components
- API calls isolated in `lib/api.ts` or an API folder
- shared types in `types/`
- no huge components with everything inside
- no unnecessary state management library until needed
- use REST polling for MVP

Frontend style can be defined properly after the first real screens are built.

## What To Avoid

Avoid:

- huge files with mixed responsibilities
- premature interfaces
- generic service factories
- deeply nested abstractions
- hidden global state
- magic string contracts spread everywhere
- database queries inside UI or HTTP handler code
- worker code that knows user/auth/database details
- scheduler code that performs HTTP requests
- code that looks generated but nobody wants to maintain

The code should look like it was written by someone who understands the system.

