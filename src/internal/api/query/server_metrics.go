package query

import (
	"encoding/json"
	"fmt"
	"net/http"
	"syscall"
	"time"

	Endpoints "watchtower/pkg/endpoints"
	Handlers "watchtower/internal/api/handlers"

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
		Handlers.MakeGetHandler(GenerateHostCPUJSON),
	)
	log.Debug().Str("server_cpu", Endpoints.QueryServerCPU).
		Msg("Server CPU Endpoint Initialized")
}

// Helper function to get server CPU information and format into HTTP response
func GenerateHostCPUJSON() map[string]any {
	info, _ := cpu.Info()

	logicalCores, _ := cpu.Counts(true)
	physicalCores, _ := cpu.Counts(false)
	usedPercentage, _ := cpu.Percent(time.Second, false)

	return map[string]any{
		"cpu_model":           info[0].Model,
		"cpu_family":          info[0].Family,
		"cpu_model_name":      info[0].ModelName,
		"cpu_mhz":             info[0].Mhz / 1000,
		"cpu_cache_size":      info[0].CacheSize,
		"cpu_logical_cores":   logicalCores,
		"cpu_physical_cores":  physicalCores,
		"cpu_used_percentage": usedPercentage[0] * 100,
	}
}

// ----- Server Memory Metrics ------------------------------------------
// Function to initialize server_memory API endpoint
func initializeServerMemoryEndpoint() {
	http.HandleFunc(
		Endpoints.QueryServerMemory,
		Handlers.MakeGetHandler(GenerateHostMemoryJSON),
	)
	log.Debug().Str("server_memory", Endpoints.QueryServerMemory).
		Msg("Server Memory Endpoint Initialized")
}

// Helper function to get server Memory information and format into HTTP response
func GenerateHostMemoryJSON() map[string]any {
	v, _ := mem.VirtualMemory()

	return map[string]any{
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
		Handlers.MakeGetHandler(GenerateHostStorageJSON),
	)
	log.Debug().Str("server_storage", Endpoints.QueryServerStorage).
		Msg("Server Storage Endpoint Initialized")
}

// Helper function to get server storage information and format into HTTP response
func GenerateHostStorageJSON() map[string]any {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		log.Err(err).Msg("Failed to GenerateServerStorageJSON")
	}

	// Total storage
	totalStorageBytes := stat.Blocks * uint64(stat.Bsize)

	// Free storage
	freeStorageBytes := stat.Bfree * uint64(stat.Bsize)
	freeStoragePercentage := (float64(freeStorageBytes)/float64(totalStorageBytes))*100

	// Used storage
	usedStorageBytes := totalStorageBytes - freeStorageBytes
	usedStoragePercentage := (float64(usedStorageBytes)/float64(totalStorageBytes))*100

	return map[string]any{
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
		b, _ := json.MarshalIndent(GenerateHostTempJSON(), "", "  ")
		fmt.Fprintf(w, string(b))
	})

	log.Debug().Str("server_temp", Endpoints.QueryServerTemp).
		Msg("Server Temp Endpoint Initialized")
}

type SensorData struct {
	Sensor  string `json:"sensor"`
	Celsius *float64    `json:"celsius"`
}

type TemperatureResponse struct {
	Host string `json:"host"`
	Data []SensorData `json:"data"`
}

// Helper function to get server Temperature information and format into HTTP response
func GenerateHostTempJSON() TemperatureResponse {
	// response := map[string]any{}
	response := TemperatureResponse{
		Host: "watchtower_server",
		Data: []SensorData{},
	}

	temps, err := host.SensorsTemperatures()
	if err != nil || len(temps) == 0 {
		//log.Error().Str("server_temp", "Unavailable").Msg("Failed to get host.SensorsTemperatures()")
		newSensorData := SensorData{
			Sensor:  "Temperature Not Available",
			Celsius: nil,
		}
		response.Data = append(response.Data, newSensorData)
	}

	for _, t := range temps {
		sensor := fmt.Sprintf("sensor_%s", t.SensorKey)
		temp := t.Temperature
		newSensorData := SensorData{
			Sensor:  sensor,
			Celsius: &temp,
		}
		response.Data = append(response.Data, newSensorData)
	}

	return response
}
