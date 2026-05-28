CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- Host CPU Usage table ---------------------------------------
CREATE TABLE IF NOT EXISTS host_cpu_usage (
    timestamp           TIMESTAMPTZ         NOT NULL,
    host                TEXT                NOT NULL,
    cpu_used_percentage DOUBLE PRECISION
);
SELECT create_hypertable('host_cpu_usage', 'timestamp', if_not_exists => TRUE);

CREATE INDEX idx_cpu_time ON host_cpu_usage (host, timestamp DESC);

-- Enable compression
ALTER TABLE host_cpu_usage SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'host'
);

-- Host Memory Usage table ------------------------------------
CREATE TABLE IF NOT EXISTS host_memory_usage (
    timestamp               TIMESTAMPTZ     NOT NULL,
    host                    TEXT            NOT NULL,
    total_memory_bytes      DOUBLE PRECISION,
    free_memory_bytes       DOUBLE PRECISION,
    used_memory_bytes       DOUBLE PRECISION,
    free_memory_percent     DOUBLE PRECISION,
    used_memory_percent     DOUBLE PRECISION
);
SELECT create_hypertable('host_memory_usage', 'timestamp', if_not_exists => TRUE);

CREATE INDEX idx_memory_time ON host_memory_usage (host, timestamp DESC);
ALTER TABLE host_memory_usage SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'host'
);

-- Host Storage Usage table ----------------------------------
CREATE TABLE IF NOT EXISTS host_storage_usage (
    timestamp               TIMESTAMPTZ     NOT NULL,
    host                    TEXT            NOT NULL,
    total_storage_bytes      NUMERIC(20,0) CHECK (total_storage_bytes >= 0),
    free_storage_bytes       NUMERIC(20,0) CHECK (free_storage_bytes >= 0),
    used_storage_bytes       NUMERIC(20,0) CHECK (used_storage_bytes >= 0),
    free_storage_percent     DOUBLE PRECISION,
    used_storage_percent     DOUBLE PRECISION
);
SELECT create_hypertable('host_storage_usage', 'timestamp', if_not_exists => TRUE);

CREATE INDEX idx_storage_time ON host_storage_usage (host, timestamp DESC);
ALTER TABLE host_storage_usage SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'host'
);

-- Host Temperature Usage table ------------------------------
CREATE TABLE IF NOT EXISTS host_temperature (
    timestamp       TIMESTAMPTZ     NOT NULL,
    host            TEXT            NOT NULL,
    sensor_name     TEXT            NOT NULL,
    temp_celsius    DOUBLE PRECISION    
);
SELECT create_hypertable('host_temperature', 'timestamp', if_not_exists =>  TRUE);

CREATE INDEX idx_temp_time ON host_temperature (host, timestamp DESC);
ALTER TABLE host_temperature SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'host'
);

-- Agent Registration table ------------------------------
CREATE TABLE IF NOT EXISTS agents (
    agent_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    host VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);