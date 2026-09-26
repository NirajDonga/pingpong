CREATE TABLE check_results (
    monitor_id UUID,
    checked_at DateTime64(3, 'UTC'),
    success Bool,
    status_code UInt16,
    response_time_ms UInt32,
    dns_ms UInt32,
    tcp_ms UInt32,
    tls_ms UInt32,
    ttfb_ms UInt32,
    error String,
    worker_name String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(checked_at)
ORDER BY (monitor_id, checked_at);
