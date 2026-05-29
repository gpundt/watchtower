CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- Host CPU Usage table ---------------------------------------
CREATE TABLE IF NOT EXISTS host_cpu_usage (
    time           TIMESTAMPTZ         NOT NULL,
    hostname                TEXT                NOT NULL,
    cpu_used_percentage DOUBLE PRECISION
);
SELECT create_hypertable('host_cpu_usage', 'time', if_not_exists => TRUE);

CREATE INDEX idx_cpu_time ON host_cpu_usage (hostname, time DESC);

-- Enable compression
ALTER TABLE host_cpu_usage SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'hostname'
);

-- Host Memory Usage table ------------------------------------
CREATE TABLE IF NOT EXISTS host_memory_usage (
    time               TIMESTAMPTZ     NOT NULL,
    hostname                    TEXT            NOT NULL,
    total_memory_bytes      DOUBLE PRECISION,
    free_memory_bytes       DOUBLE PRECISION,
    used_memory_bytes       DOUBLE PRECISION,
    free_memory_percent     DOUBLE PRECISION,
    used_memory_percent     DOUBLE PRECISION
);
SELECT create_hypertable('host_memory_usage', 'time', if_not_exists => TRUE);

CREATE INDEX idx_memory_time ON host_memory_usage (hostname, time DESC);
ALTER TABLE host_memory_usage SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'hostname'
);

-- Host Storage Usage table ----------------------------------
CREATE TABLE IF NOT EXISTS host_storage_usage (
    time               TIMESTAMPTZ     NOT NULL,
    hostname                    TEXT            NOT NULL,
    total_storage_bytes      NUMERIC(20,0) CHECK (total_storage_bytes >= 0),
    free_storage_bytes       NUMERIC(20,0) CHECK (free_storage_bytes >= 0),
    used_storage_bytes       NUMERIC(20,0) CHECK (used_storage_bytes >= 0),
    free_storage_percent     DOUBLE PRECISION,
    used_storage_percent     DOUBLE PRECISION
);
SELECT create_hypertable('host_storage_usage', 'time', if_not_exists => TRUE);

CREATE INDEX idx_storage_time ON host_storage_usage (hostname, time DESC);
ALTER TABLE host_storage_usage SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'hostname'
);

-- Host Temperature Usage table ------------------------------
CREATE TABLE IF NOT EXISTS host_temperature (
    time       TIMESTAMPTZ     NOT NULL,
    hostname            TEXT            NOT NULL,
    sensor_name     TEXT            NOT NULL,
    temp_celsius    DOUBLE PRECISION    
);
SELECT create_hypertable('host_temperature', 'time', if_not_exists =>  TRUE);

CREATE INDEX idx_temp_time ON host_temperature (hostname, time DESC);
ALTER TABLE host_temperature SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'hostname'
);

-- Agent Registration table ------------------------------
CREATE TABLE IF NOT EXISTS agents (
    agent_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hostname VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Host Port Scan Table ---------------------------------
CREATE TABLE IF NOT EXISTS port_scan (
    hostname    TEXT        PRIMARY KEY,
    open_ports INTEGER[] DEFAULT '{}'::INTEGER[],
    last_scan_timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)