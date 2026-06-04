package database

import (
	"database/sql"
	"fmt"
	"time"

	Query "watchtower/internal/api/query"
	Config "watchtower/internal/config"

	"github.com/lib/pq" // PostgreSQL Driver
	"github.com/rs/zerolog/log"
)

// Constant strings for database connection
const (
	CPUUsageTable     = "host_cpu_usage"
	MemoryUsageTable  = "host_memory_usage"
	StorageUsageTable = "host_storage_usage"
	TemperatureTable  = "host_temperature"
	AgentsTable       = "agents"
	PortScanTable     = "port_scan"
	AuthLogEntryTable     = "auth_logs"
	ContainerLogEntryTable = "conainer_logs"
	CronLogEntryTable = "cron_logs"
	KernelLogEntryTable = "kernel_logs"
	ServiceLogEntryTable = "service_logs"
)

// Helper function to build the connection string fresh each call
func dbConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		Config.ServerConfig.Database.Host,
		Config.ServerConfig.Database.Port,
		Config.ServerConfig.Database.User,
		Config.ServerConfig.Database.Password,
		Config.ServerConfig.Database.Name,
		Config.ServerConfig.Database.SSLMode,
	)
}

// Helper function to insert a Host CPU Metrics Submission into the database
func InsertHostCPUUsage(
	timestamp time.Time,
	host string,
	totalCores int,
	corePercentages []float64,
) error {
	db, err := sql.Open("postgres", dbConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`INSERT INTO %s (time, hostname, total_cores, core_percentages) VALUES ($1, $2, $3, $4)`,
		CPUUsageTable,
	)

	_, execErr := db.Exec(
		sqlStatement,
		timestamp,
		host,
		totalCores,
		corePercentages,
	)
	if execErr != nil {
		log.Err(execErr).Str("func", "Database.InsertHostCPUUsage").Msg("")
		return execErr
	}
	return nil
}

// Helper function to insert a Host Memory Metrics Submission into the database
func InsertHostMemoryUsage(
	timestamp time.Time,
	host string,
	total_bytes,
	free_bytes,
	used_bytes,
	free_percentage,
	used_percentage float64,
) error {
	db, err := sql.Open("postgres", dbConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`INSERT INTO %s (time, hostname, total_memory_bytes, free_memory_bytes, used_memory_bytes, free_memory_percent, used_memory_percent) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		MemoryUsageTable,
	)

	_, execErr := db.Exec(
		sqlStatement,
		timestamp,
		host,
		total_bytes,
		free_bytes,
		used_bytes,
		free_percentage,
		used_percentage,
	)
	if execErr != nil {
		log.Err(execErr).Str("func", "Database.InsertHostMemoryUsage").Msg("")
		return execErr
	}
	return nil
}

// Helper function to insert a Host Storage Metrics submission into the database
func InsertHostStorageUsage(
	timestamp time.Time,
	host string,
	total_bytes,
	free_bytes,
	used_bytes uint64,
	free_percentage,
	used_percentage float64,
) error {
	db, err := sql.Open("postgres", dbConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`INSERT INTO %s (time, hostname, total_storage_bytes, free_storage_bytes, used_storage_bytes, free_storage_percent, used_storage_percent) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		StorageUsageTable,
	)

	_, execErr := db.Exec(
		sqlStatement,
		timestamp,
		host,
		total_bytes,
		free_bytes,
		used_bytes,
		free_percentage,
		used_percentage,
	)
	if execErr != nil {
		log.Err(execErr).Str("func", "Database.InsertHostStorageUsage").Msg("")
		return execErr
	}
	return nil
}

// Helper function to insert a Host Temperature submission into the database
func InsertHostTemperature(
	timestamp time.Time,
	host string,
	tempData []Query.SensorData,
) error {
	db, err := sql.Open("postgres", dbConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`INSERT INTO %s (time, hostname, sensor_name, temp_celsius) VALUES ($1, $2, $3, $4)`,
		TemperatureTable,
	)

	for _, sensorStruct := range tempData {
		_, execErr := db.Exec(
			sqlStatement,
			timestamp,
			host,
			sensorStruct.Sensor,
			sensorStruct.Celsius,
		)
		if execErr != nil {
			log.Err(execErr).Str("func", "Database.InsertHostTemperature").Msg("")
			return execErr
		}
	}
	return nil
}

// Function to insert or update an agent into the agents db
func InsertAgent(
	host string,
	timestamp time.Time,
) error {
	db, err := sql.Open("postgres", dbConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`INSERT INTO %s (hostname, created_at, updated_at) VALUES ($1, $2, $3) ON CONFLICT (hostname) DO UPDATE SET updated_at = $3`,
		AgentsTable,
	)
	_, execErr := db.Exec(
		sqlStatement,
		host,
		timestamp,
		timestamp,
	)
	if execErr != nil {
		log.Err(execErr).Str("func", "Database.InsertAgent").Msg("")
		return execErr
	}
	return nil
}

// Function to update a row in the agent table, triggered by HostCheckIn
func UpdateAgent(
	host string,
	timestamp time.Time,
) error {
	db, err := sql.Open("postgres", dbConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`UPDATE %s SET updated_at = $1 WHERE hostname = $2`,
		AgentsTable,
	)
	_, execErr := db.Exec(
		sqlStatement,
		timestamp,
		host,
	)
	if execErr != nil {
		log.Err(execErr).Str("func", "Database.UpdateAgent").Msg("")
		return execErr
	}
	return nil
}

// Function to insert individual port scan results into psql db
func InsertPortScan(
	host string,
	openPorts []int,
	timestamp time.Time,
) error {
	db, err := sql.Open("postgres", dbConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`INSERT INTO %s (hostname, open_ports, last_scan_timestamp) VALUES ($1, $2, $3) ON CONFLICT (hostname) DO UPDATE SET open_ports = $2, last_scan_timestamp = $3`,
		PortScanTable,
	)

	_, execErr := db.Exec(
		sqlStatement,
		host,
		pq.Array(openPorts),
		timestamp,
	)
	if execErr != nil {
		log.Err(execErr).Str("func", "Database.InsertPortScan").Msg("")
		return execErr
	}

	return nil
}

// Function to insert individual log entry into logs database
func InsertLogEntry(
	tableName string,
	timestamp time.Time,
	hostname,
	severity,
	logMessage,
	serviceName,
	user string,
) error {
	db, err := sql.Open("postgres", dbConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`INSERT INTO %s (time, hostname, severity, log_message, "service", "user") VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (time, hostname, log_message) DO NOTHING`,
		tableName,
	)

	_, execErr := db.Exec(
		sqlStatement,
		timestamp,
		hostname,
		severity,
		logMessage,
		serviceName,
		user,
	)
	if execErr != nil {
		log.Err(execErr).Str("func", "Database.InsertLogEntry").Msg("")
		return execErr
	}

	return nil
}
