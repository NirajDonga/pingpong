# Services

## web

User dashboard.

Responsibilities:

- register/login UI
- monitor list
- create monitor form
- interval selector
- pause/resume monitor controls
- monitor detail page
- response time chart
- uptime history
- recent check table
- incident list/detail

MVP uses REST + polling. WebSocket can come later.

## api

Main backend service.

Responsibilities:

- authentication
- user management
- validate monitor requests
- create/list/update/delete monitors
- enforce user ownership
- consume check results from NATS
- write raw check results to ClickHouse
- update monitor current status in PostgreSQL
- open/close incidents
- expose REST APIs for web
- own database schemas and migrations

## scheduler

Separate backend service.

Responsibilities:

- read monitor schedule data from PostgreSQL
- find monitors due for checking
- publish check jobs to NATS
- update `next_check_at`

How scheduler knows what to check:

```txt
api saves monitors in PostgreSQL.
Each monitor has enabled, interval_seconds, and next_check_at.
scheduler runs every second.
scheduler asks PostgreSQL for monitors where:
  enabled = true
  next_check_at <= current time
Those monitors are due.
scheduler publishes one check job for each due monitor.
scheduler updates next_check_at for each monitor.
```

Example query:

```sql
SELECT *
FROM monitors
WHERE enabled = true
AND next_check_at <= NOW();
```

Example NATS job:

```json
{
  "monitorId": "mon_123",
  "url": "https://example.com",
  "timeoutSeconds": 5,
  "expectedStatus": 200
}
```

After publishing the job:

```sql
UPDATE monitors
SET next_check_at = NOW() + interval_seconds
WHERE id = 'mon_123';
```

## worker

Background service that performs HTTP checks.

Responsibilities:

- receive check job from NATS
- perform one HTTP request
- measure response time
- collect DNS/TCP/TLS/TTFB timings when possible
- publish success/failure result to NATS

Worker input:

```txt
monitor_id
url
timeout_seconds
expected_status
```

Worker output:

```txt
monitor_id
checked_at
success
status_code
response_time_ms
dns_ms
tcp_ms
tls_ms
ttfb_ms
error
worker_name
```

Worker does not access PostgreSQL or ClickHouse.

## infra

Local dependency setup.

Responsibilities:

- PostgreSQL
- ClickHouse
- NATS
- Docker Compose config
- DB init files if needed

`infra` runs dependencies. It does not own app schema.

