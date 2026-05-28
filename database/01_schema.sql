-- 1. Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- 2. Create host_cpu_usage table
CREATE TABLE IF NOT EXISTS host_cpu_usage (
    timestamp           TIMESTAMPTZ         NOT NULL,
    host                TEXT                NOT NULL,
    cpu_used_percentage DOUBLE PRECISION    NOT NULL
)

-- 3. Covert the table into a hypertable
SELECT create_hypertable("host_cpu_usage", "timestamp", if_not_exists => TRUE);

-- 4. Create an index for common query patterns
CREATE INDEX ON host_cpu_usage (host, timestamp DESC)