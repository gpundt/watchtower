package query

import (
	"encoding/json"
	"fmt"
	"net/http"
	"syscall"
	"time"

	Handlers "watchtower/internal/api/handlers"
	Config "watchtower/internal/config"
	Endpoints "watchtower/pkg/endpoints"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v4/mem"
)

// ----- Server CPU Metrics ---------------------------------------------
// Function to initialize server_cpu API endpoint
func initializeServerCPUEndpoint() {
	http.HandleFunc(
		Endpoints.QueryServerCPU,
		Handlers.MakeGetHandler(func() map[string]any {
			return GenerateHostCPUJSON(Config.ServerConfig.Server.Host)
		}),
	)
	log.Debug().Str("server_cpu", Endpoints.QueryServerCPU).
		Msg("Server CPU Endpoint Initialized")
}

// Helper function to get server CPU information and format into HTTP response
func GenerateHostCPUJSON(sourceHost string) map[string]any {
	percentages, err := cpu.Percent(1*time.Second, true)
	if err != nil {
		log.Err(err).Str("func", "GenerateHostCPUJSON").Msg("")
		return nil
	}

	jsonData := map[string]any{
		"host":             sourceHost,
		"total_cores":      len(percentages),
		"core_percentages": percentages,
	}

	return jsonData
}

// ----- Server Memory Metrics ------------------------------------------
// Function to initialize server_memory API endpoint
func initializeServerMemoryEndpoint() {
	http.HandleFunc(
		Endpoints.QueryServerMemory,
		Handlers.MakeGetHandler(func() map[string]any {
			return GenerateHostMemoryJSON(Config.ServerConfig.Server.Host)
		}),
	)
	log.Debug().Str("server_memory", Endpoints.QueryServerMemory).
		Msg("Server Memory Endpoint Initialized")
}

// Helper function to get server Memory information and format into HTTP response
func GenerateHostMemoryJSON(sourceHost string) map[string]any {
	v, _ := mem.VirtualMemory()

	return map[string]any{
		"host":                   sourceHost,
		"total_memory_bytes":     float64(v.Total),
		"free_memory_bytes":      float64(v.Free),
		"free_memory_percentage": float64(100 - v.UsedPercent),
		"used_memory_bytes":      float64(v.Used),
		"used_memory_percentage": float64(v.UsedPercent),
	}
}

// ----- Server Storage Metrics ----------------------------------------
// Function to initialize server_storage API endpoint
func initializeServerStorageEndpoint() {
	http.HandleFunc(
		Endpoints.QueryServerStorage,
		Handlers.MakeGetHandler(func() map[string]any {
			return GenerateHostStorageJSON(Config.ServerConfig.Server.Host)
		}),
	)
	log.Debug().Str("server_storage", Endpoints.QueryServerStorage).
		Msg("Server Storage Endpoint Initialized")
}

// Helper function to get server storage information and format into HTTP response
func GenerateHostStorageJSON(sourceHost string) map[string]any {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		log.Err(err).Msg("Failed to GenerateServerStorageJSON")
	}

	// Total storage
	totalStorageBytes := stat.Blocks * uint64(stat.Bsize)

	// Free storage
	freeStorageBytes := stat.Bfree * uint64(stat.Bsize)
	freeStoragePercentage := (float64(freeStorageBytes) / float64(totalStorageBytes)) * 100

	// Used storage
	usedStorageBytes := totalStorageBytes - freeStorageBytes
	usedStoragePercentage := (float64(usedStorageBytes) / float64(totalStorageBytes)) * 100

	return map[string]any{
		"host":                    sourceHost,
		"total_storage_bytes":     totalStorageBytes,
		"free_storage_bytes":      freeStorageBytes,
		"free_storage_percentage": freeStoragePercentage,
		"used_storage_bytes":      usedStorageBytes,
		"used_storage_percentage": usedStoragePercentage,
	}
}

// ----- Server Temperature Metrics -----------------------------------
// Function to initialize server_temp API endpoint
func initializeServerTempEndpoint() {
	http.HandleFunc(Endpoints.QueryServerTemp, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		b, _ := json.MarshalIndent(
			GenerateHostTempJSON(Config.ServerConfig.Server.Host), "", "  ",
		)
		fmt.Fprintf(w, string(b))
	})

	log.Debug().Str("server_temp", Endpoints.QueryServerTemp).
		Msg("Server Temp Endpoint Initialized")
}

type SensorData struct {
	Sensor  string  `json:"sensor"`
	Celsius float64 `json:"celsius"`
}

type TemperatureResponse struct {
	Host string       `json:"host"`
	Data []SensorData `json:"data"`
}

// Helper function to get server Temperature information and format into HTTP response
func GenerateHostTempJSON(sourceHost string) TemperatureResponse {
	// response := map[string]any{}
	response := TemperatureResponse{
		Host: sourceHost,
		Data: []SensorData{},
	}

	temps, err := host.SensorsTemperatures()
	if err != nil || len(temps) == 0 {
		//log.Error().Str("server_temp", "Unavailable").Msg("Failed to get host.SensorsTemperatures()")
		newSensorData := SensorData{
			Sensor:  "Temperature Not Available",
			Celsius: 0.0,
		}
		response.Data = append(response.Data, newSensorData)
	}

	for _, t := range temps {
		sensor := fmt.Sprintf("sensor_%s", t.SensorKey)
		temp := t.Temperature
		newSensorData := SensorData{
			Sensor:  sensor,
			Celsius: temp,
		}
		response.Data = append(response.Data, newSensorData)
	}

	return response
}
