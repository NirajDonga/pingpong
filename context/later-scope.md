# Later Scope

Do not build these in the first version:

```txt
multi-region checks
alerts
notification-manager
teams/workspaces
RBAC
public status pages
billing
TCP/DNS/TLS checks
complex response assertions
WebSocket realtime
Redis scheduler queue
```

First goal:

```txt
simple persistent uptime monitoring with REST polling
```

Redis can be added later if PostgreSQL scheduler writes become too much.

