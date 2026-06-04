CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- Host CPU Usage table ---------------------------------------
CREATE TABLE IF NOT EXISTS host_cpu_usage (
    time           TIMESTAMPTZ         NOT NULL,
    hostname                TEXT                NOT NULL,
    total_cores         INTEGER,
    core_percentages DOUBLE PRECISION[] DEFAULT '{}'::DOUBLE PRECISION[]
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
    time                 TIMESTAMPTZ      NOT NULL,
    hostname             TEXT             NOT NULL,
    total_storage_bytes  NUMERIC(20,0)    CHECK (total_storage_bytes >= 0),
    free_storage_bytes   NUMERIC(20,0)    CHECK (free_storage_bytes >= 0),
    used_storage_bytes   NUMERIC(20,0)    CHECK (used_storage_bytes >= 0),
    free_storage_percent DOUBLE PRECISION,
    used_storage_percent DOUBLE PRECISION
);
SELECT create_hypertable('host_storage_usage', 'time', if_not_exists => TRUE);

CREATE INDEX idx_storage_time ON host_storage_usage (hostname, time DESC);
ALTER TABLE host_storage_usage SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'hostname'
);

-- Host Temperature Usage table ------------------------------
CREATE TABLE IF NOT EXISTS host_temperature (
    time         TIMESTAMPTZ      NOT NULL,
    hostname     TEXT             NOT NULL,
    sensor_name  TEXT             NOT NULL,
    temp_celsius DOUBLE PRECISION    
);
SELECT create_hypertable('host_temperature', 'time', if_not_exists =>  TRUE);

CREATE INDEX idx_temp_time ON host_temperature (hostname, time DESC);
ALTER TABLE host_temperature SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'hostname'
);

-- Agent Registration table ------------------------------
CREATE TABLE IF NOT EXISTS agents (
    agent_id   UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    hostname   VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
);

-- Host Port Scan Table ---------------------------------
CREATE TABLE IF NOT EXISTS port_scan (
    hostname            TEXT        PRIMARY KEY,
    open_ports          INTEGER[]   DEFAULT '{}'::INTEGER[],
    last_scan_timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

--- Host Auth Log Entry Table --------------------------------
CREATE TABLE IF NOT EXISTS auth_logs (
    time             TIMESTAMPTZ NOT NULL,
    hostname         TEXT        NOT NULL,
    severity         TEXT        NOT NULL,
    log_message      TEXT        NOT NULL,
    "service"        TEXT,
    "user"           TEXT,

    CONSTRAINT auth_logs_unique_entry UNIQUE (time, hostname, log_message)
);
SELECT create_hypertable('auth_logs', 'time', if_not_exists => TRUE);

CREATE INDEX idx_auth_log_time ON auth_logs (hostname, time DESC);

ALTER TABLE auth_logs SET(
    timescaledb.compress,
    timescaledb.compress_segmentby = 'hostname'
);

--- Host Container Log Entry Table --------------------------------
CREATE TABLE IF NOT EXISTS container_logs (
    time             TIMESTAMPTZ NOT NULL,
    hostname         TEXT        NOT NULL,
    severity         TEXT        NOT NULL,
    log_message      TEXT        NOT NULL,
    "service"        TEXT,
    "user"           TEXT,

    CONSTRAINT container_logs_unique_entry UNIQUE (time, hostname, log_message)
);
SELECT create_hypertable('container_logs', 'time', if_not_exists => TRUE);

CREATE INDEX idx_container_log_time ON container_logs (hostname, time DESC);

ALTER TABLE container_logs SET(
    timescaledb.compress,
    timescaledb.compress_segmentby = 'hostname'
);

--- Host CRON Log Entry Table --------------------------------
CREATE TABLE IF NOT EXISTS cron_logs (
    time             TIMESTAMPTZ NOT NULL,
    hostname         TEXT        NOT NULL,
    severity         TEXT        NOT NULL,
    log_message      TEXT        NOT NULL,
    "service"        TEXT,
    "user"           TEXT,

    CONSTRAINT cron_logs_unique_entry UNIQUE (time, hostname, log_message)
);
SELECT create_hypertable('cron_logs', 'time', if_not_exists => TRUE);

CREATE INDEX idx_cron_log_time ON cron_logs (hostname, time DESC);

ALTER TABLE cron_logs SET(
    timescaledb.compress,
    timescaledb.compress_segmentby = 'hostname'
);

--- Host Kernel Log Entry Table --------------------------------
CREATE TABLE IF NOT EXISTS kernel_logs (
    time             TIMESTAMPTZ NOT NULL,
    hostname         TEXT        NOT NULL,
    severity         TEXT        NOT NULL,
    log_message      TEXT        NOT NULL,
    "service"        TEXT,
    "user"           TEXT,

    CONSTRAINT kernel_logs_unique_entry UNIQUE (time, hostname, log_message)
);
SELECT create_hypertable('kernel_logs', 'time', if_not_exists => TRUE);

CREATE INDEX idx_kernel_log_time ON kernel_logs (hostname, time DESC);

ALTER TABLE kernel_logs SET(
    timescaledb.compress,
    timescaledb.compress_segmentby = 'hostname'
);

--- Host Service Log Entry Table --------------------------------
CREATE TABLE IF NOT EXISTS service_logs (
    time             TIMESTAMPTZ NOT NULL,
    hostname         TEXT        NOT NULL,
    severity         TEXT        NOT NULL,
    log_message      TEXT        NOT NULL,
    "service"        TEXT,
    "user"           TEXT,

    CONSTRAINT service_logs_unique_entry UNIQUE (time, hostname, log_message)
);
SELECT create_hypertable('service_logs', 'time', if_not_exists => TRUE);

CREATE INDEX idx_service_log_time ON service_logs (hostname, time DESC);

ALTER TABLE service_logs SET(
    timescaledb.compress,
    timescaledb.compress_segmentby = 'hostname'
);