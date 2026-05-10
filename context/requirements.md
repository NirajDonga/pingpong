# Requirements

## User Requirements

- User can register/login.
- User can create a monitor.
- User can set URL.
- User can set check interval.
- User can set timeout.
- User can pause/resume monitor.
- User can see current monitor status.
- User can see response time history.
- User can see recent check results.
- User can see downtime incidents.

## Monitor Requirements

Each monitor must have:

```txt
id
user_id
name
url
interval_seconds
timeout_seconds
expected_status
enabled
current_status
consecutive_failures
consecutive_successes
next_check_at
created_at
updated_at
```

## Check Result Requirements

Each check result must store:

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

## Multi-User Requirements

Every monitor belongs to one user.

User-facing queries must only return data owned by the current user.

Correct monitor query:

```sql
SELECT *
FROM monitors
WHERE user_id = ?;
```

Check history must verify ownership through the monitor.

## Status Requirements

Do not mark a monitor down after one failed request.

Default rule:

```txt
3 failed checks in a row = down
1 successful check after down = up again
```

Initial statuses:

```txt
unknown
up
down
```

Add `degraded` later.

## Incident Requirements

Incident means one downtime period.

When monitor becomes down:

```txt
open incident
set started_at
store reason
```

When monitor comes back up:

```txt
close incident
set ended_at
```

