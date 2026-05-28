package database

import (
	"database/sql"
	"fmt"
	"time"

	Config "watchtower/internal/config"

	"github.com/rs/zerolog/log"
	_ "github.com/lib/pq"	// PostgreSQL Driver
)

// Constant strings for database connection
const (
	DBConnection := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		Config.ServerConfig.Database.Host,
		Config.ServerConfig.Database.Port,
		Config.ServerConfig.Database.User,
		Config.ServerConfig.Database.Password,
		Config.ServerConfig.Database.Name,
		Config.ServerConfig.Database.SSLMode,
	)
	CPUUsageTable = "host_cpu_usage"
	MemoryUsageTable = "host_memory_usage"
	StorageUsageTable = "host_storage_usage"
	TemperatureTable = "host_temperature"
)

// Helper function to insert a Host CPU Metrics Submission into the database
func InsertHostCPUUsage(
	timestamp time.Time,
	host string,
	used_percentage float64,
) error {
	db, err := sql.Open("postgres", DBConnection)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`INSERT INTO %S VALUES ($1, $2, $3)`,
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
	db, err := sql.Open("postgres", DBConnection)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`INSERT INTO %s VALUES ($1, $2, $3, $4, $5, $6, $7)`,
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
	used_percentage float64
) error {
	db, err := sql.Open("postgres", DBConnection)
	if err != nil {
		return err
	}
	defer db.close()

	sqlStatement := fmt.Sprintf(
		`INSERT INTO %s VALUES ($1, $2, $3, $4, $5, $6, $7)`,
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
	tempData map[string][]any,
) error {
	db, err := sql.Open("postgres", DBConnection)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`INSERT INTO %s VALUES ($1, $2, $3, $4)`
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