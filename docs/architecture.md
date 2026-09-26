# Architecture

Target services:

```txt
api
scheduler
worker
web
infra
```

No `uptime-` prefix in folder names.

High-level architecture:

```txt
+---------------------+
| web                 |
| Next.js dashboard   |
+----------+----------+
           |
           | REST + polling
           v
+---------------------+          +----------------------+
| api                 |          | PostgreSQL           |
| REST API            +--------->| users                |
| auth                |          | monitors             |
| monitor service     |          | incidents            |
| result processor    |          | scheduler timing     |
| incident service    |          +----------+-----------+
+----------+----------+                     ^
           |                                |
           | consumes check results         | reads due monitors
           v                                |
+---------------------+          +----------+-----------+
| ClickHouse          |          | scheduler            |
| time-series checks  |          | finds due monitors   |
| response history    |          | publishes jobs       |
+---------------------+          +----------+-----------+
                                             |
                                             | check jobs
                                             v
                                  +----------+-----------+
                                  | NATS                 |
                                  | jobs and results     |
                                  +----------+-----------+
                                             |
                                             | job
                                             v
                                  +----------+-----------+
                                  | worker               |
                                  | performs HTTP check  |
                                  | publishes result     |
                                  +----------------------+
```

Flow:

```txt
1. User creates a monitor from web.
2. api receives the request.
3. api validates the request.
4. api stores the monitor in PostgreSQL.
5. Monitor row contains interval_seconds and next_check_at.
6. scheduler reads due monitors from PostgreSQL.
7. scheduler finds monitors where next_check_at <= now.
8. scheduler publishes one check job to NATS.
9. scheduler updates next_check_at = now + interval_seconds.
10. worker receives the check job from NATS.
11. worker sends HTTP request to the monitor URL.
12. worker publishes the check result to NATS.
13. api consumes the result from NATS.
14. api stores the raw check result in ClickHouse.
15. api updates monitor current status in PostgreSQL.
16. api opens or closes incidents if status changed.
17. web reads latest status/history through REST polling.
```

Boundaries:

```txt
api = validates user requests and owns business data
scheduler = timing only
worker = HTTP request only
ClickHouse = time-series result history
PostgreSQL = app data and current state
```

The API does not keep a request open for 60 seconds. It validates and saves the monitor. The repeated "every 60 seconds" work belongs to scheduler.

