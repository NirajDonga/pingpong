# Overview

PingPong is an uptime monitoring platform.

Users add website URLs. For each URL, the user selects how often it should be checked. PingPong sends requests at that interval, stores every result, and shows whether the website is up or down over time.

Example:

```txt
URL: https://example.com
Interval: 60 seconds
Timeout: 5 seconds
```

Main idea:

```txt
PostgreSQL = users, monitors, incidents, current state, scheduler timing
ClickHouse = time-series check results
scheduler = decides when to check
worker = performs the HTTP request
api = validates requests, owns schemas, stores monitors, processes results
web = user dashboard
```

First goal:

```txt
simple persistent uptime monitoring with REST polling
```

