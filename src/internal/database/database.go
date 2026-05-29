package database

import (
	"database/sql"
	"fmt"
	"time"

	Query "watchtower/internal/api/query"
	Config "watchtower/internal/config"

	// "github.com/rs/zerolog/log"
	_ "github.com/lib/pq" // PostgreSQL Driver
)

// Constant strings for database connection
const (
	CPUUsageTable     = "host_cpu_usage"
	MemoryUsageTable  = "host_memory_usage"
	StorageUsageTable = "host_storage_usage"
	TemperatureTable  = "host_temperature"
	AgentsTable       = "agents"
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
	used_percentage float64,
) error {
	db, err := sql.Open("postgres", dbConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`INSERT INTO %s (time, hostname, cpu_used_percentage) VALUES ($1, $2, $3)`,
		CPUUsageTable,
	)

	_, err = db.Exec(
		sqlStatement,
		timestamp,
		host,
		used_percentage,
	)
	if err != nil {
		return err
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

	_, err = db.Exec(
		sqlStatement,
		timestamp,
		host,
		total_bytes,
		free_bytes,
		used_bytes,
		free_percentage,
		used_percentage,
	)
	if err != nil {
		return err
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

	_, err = db.Exec(
		sqlStatement,
		timestamp,
		host,
		total_bytes,
		free_bytes,
		used_bytes,
		free_percentage,
		used_percentage,
	)
	if err != nil {
		return err
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
		_, err = db.Exec(
			sqlStatement,
			timestamp,
			host,
			sensorStruct.Sensor,
			sensorStruct.Celsius,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertAgentRegistration(
	host string,
) error {
	db, err := sql.Open("postgres", dbConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`INSERT INTO %s (hostname) VALUES ($1)`,
		AgentsTable,
	)

	_, err = db.Exec(
		sqlStatement,
		host,
	)
	if err != nil {
		return err
	}

	return nil
}
