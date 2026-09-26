# Build Phases

Each phase should have one purpose. Do not mix unrelated work into a phase.

## Phase 1 - Service Skeletons

Goal: create the target service folders and bootable empty services.

Sub-phases:

- 1.1 Create `api` service skeleton.
- 1.2 Create `scheduler` service skeleton.
- 1.3 Create `worker` service skeleton.
- 1.4 Create `web` service skeleton.
- 1.5 Create `infra/docker-compose.yml`.

Done when:

```txt
Each service starts and prints a health/log message.
No business logic yet.
```

## Phase 2 - Infrastructure

Goal: run required local dependencies.

Sub-phases:

- 2.1 Add PostgreSQL to `infra`.
- 2.2 Add ClickHouse to `infra`.
- 2.3 Add NATS to `infra`.
- 2.4 Add config values for each service.

Done when:

```txt
PostgreSQL, ClickHouse, and NATS run locally.
api can connect to PostgreSQL and ClickHouse.
scheduler can connect to PostgreSQL and NATS.
worker can connect to NATS.
```

## Phase 3 - Database Schema

Goal: define persistent data.

Sub-phases:

- 3.1 Add PostgreSQL migrations under `api`.
- 3.2 Add `users` table.
- 3.3 Add `monitors` table.
- 3.4 Add `incidents` table.
- 3.5 Add ClickHouse migration for `check_results`.
- 3.6 Add migration runner or documented migration command.

Done when:

```txt
All tables can be created from migrations.
Schema lives under api ownership.
```

## Phase 4 - Auth

Goal: identify the current user.

Sub-phases:

- 4.1 Add register endpoint.
- 4.2 Add login endpoint.
- 4.3 Add password hashing.
- 4.4 Add JWT or secure session.
- 4.5 Add current-user middleware.

Done when:

```txt
User can register, login, and call an authenticated endpoint.
```

## Phase 5 - Monitor CRUD

Goal: users can manage monitors.

Sub-phases:

- 5.1 Add create monitor endpoint.
- 5.2 Add list monitors endpoint.
- 5.3 Add get monitor endpoint.
- 5.4 Add update monitor endpoint.
- 5.5 Add pause/resume monitor endpoint or field update.
- 5.6 Add delete monitor behavior.
- 5.7 Enforce `user_id` ownership in every monitor query.

Done when:

```txt
User can create and manage only their own monitors.
Monitor stores interval_seconds and next_check_at.
```

## Phase 6 - Scheduler

Goal: scheduler dispatches due checks.

Sub-phases:

- 6.1 Add scheduler store for reading due monitors from PostgreSQL.
- 6.2 Add scheduler loop.
- 6.3 Add NATS publisher for check jobs.
- 6.4 Update `next_check_at` after dispatch.
- 6.5 Add limit/batch size so one tick cannot dispatch unlimited jobs.

Done when:

```txt
scheduler finds due monitors and publishes check jobs.
scheduler does not perform HTTP requests.
```

## Phase 7 - Worker

Goal: worker performs one HTTP check per job.

Sub-phases:

- 7.1 Add NATS consumer for check jobs.
- 7.2 Add HTTP checker.
- 7.3 Add timeout handling.
- 7.4 Add expected status validation.
- 7.5 Add response timing fields.
- 7.6 Publish check result to NATS.

Done when:

```txt
worker receives job, checks URL, publishes result.
worker does not access PostgreSQL or ClickHouse.
```

## Phase 8 - Result Processing

Goal: api stores check results.

Sub-phases:

- 8.1 Add NATS consumer in `api` for check results.
- 8.2 Validate result payload.
- 8.3 Insert raw result into ClickHouse.
- 8.4 Expose monitor check history endpoint.

Done when:

```txt
Check results are stored in ClickHouse and readable through api.
```

## Phase 9 - Monitor Status

Goal: api updates current monitor state.

Sub-phases:

- 9.1 Track consecutive failures.
- 9.2 Track consecutive successes.
- 9.3 Mark monitor up/down using thresholds.
- 9.4 Store current status in PostgreSQL.

Done when:

```txt
Monitor status changes based on check results.
One failed check does not instantly mark monitor down.
```

## Phase 10 - Incidents

Goal: track downtime periods.

Sub-phases:

- 10.1 Open incident when monitor becomes down.
- 10.2 Close incident when monitor recovers.
- 10.3 Add incident list endpoint.
- 10.4 Add monitor incidents endpoint.

Done when:

```txt
Downtime has start time, end time, status, and reason.
```

## Phase 11 - Web MVP

Goal: usable dashboard without realtime.

Sub-phases:

- 11.1 Add auth screens.
- 11.2 Add monitor list.
- 11.3 Add create monitor form.
- 11.4 Add monitor detail page.
- 11.5 Add recent checks table.
- 11.6 Add response time chart.
- 11.7 Add incidents view.
- 11.8 Add REST polling.

Done when:

```txt
User can manage monitors and see uptime data from the browser.
No WebSocket required.
```

## Phase 12 - Realtime Later

Goal: replace polling with live updates later.

Sub-phases:

- 12.1 Add WebSocket endpoint.
- 12.2 Push monitor status changes.
- 12.3 Push new check results.
- 12.4 Update dashboard without polling.

Done when:

```txt
Dashboard updates without manual refresh or polling.
```

