CREATE TABLE check_results (
    monitor_id UUID,
    checked_at DateTime64(3, 'UTC'),
    success Bool,
    status_code UInt16,
    response_time_ms UInt32,
    error String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(checked_at)
ORDER BY (monitor_id, checked_at);
