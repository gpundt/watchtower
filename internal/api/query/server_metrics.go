package query

import (
	"encoding/json"
	"fmt"
	"net/http"
	"syscall"
	"time"

	Endpoints "watchtower/internal/api/endpoints"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v4/mem"
)

// Root function to initialize all api endpoints related to server metrics
func initializeServerMetricsAPIEndpoints() {
	initializeServerCPUEndpoint()
	initializeServerStorageEndpoint()
	initializeServerMemoryEndpoint()
	initializeServerTempEndpoint()
}

// Function to initialize server_cpu API endpoint
func initializeServerCPUEndpoint() {
	http.HandleFunc(Endpoints.QueryServerCPU, makeGetHandler(GenerateHostCPUJSON))
	log.Debug().Str("server_cpu", Endpoints.QueryServerCPU).Msg("Server CPU Endpoint Initialized")
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
		"cpu_used_percentage": usedPercentage[0],
	}
}

// Function to initialize server_memory API endpoint
func initializeServerMemoryEndpoint() {
	http.HandleFunc(Endpoints.QueryServerMemory, makeGetHandler(GenerateHostMemoryJSON))
	log.Debug().Str("server_memory", Endpoints.QueryServerMemory).Msg("Server Memory Endpoint Initialized")
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

// Function to initialize server_storage API endpoint
func initializeServerStorageEndpoint() {
	http.HandleFunc(Endpoints.QueryServerStorage, makeGetHandler(GenerateHostStorageJSON))
	log.Debug().Str("server_storage", Endpoints.QueryServerStorage).Msg("Server Storage Endpoint Initialized")
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
	freeStoragePercentage := fmt.Sprintf(
		"%.2f",
		(float64(freeStorageBytes)/float64(totalStorageBytes))*100,
	)

	// Used storage
	usedStorageBytes := totalStorageBytes - freeStorageBytes
	usedStoragePercentage := fmt.Sprintf(
		"%.2f",
		(float64(usedStorageBytes)/float64(totalStorageBytes))*100,
	)

	return map[string]any{
		"total_storage_bytes":     totalStorageBytes,
		"free_storage_bytes":      freeStorageBytes,
		"free_storage_percentage": freeStoragePercentage,
		"used_storage_bytes":      usedStorageBytes,
		"used_storage_percentage": usedStoragePercentage,
	}
}

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

	log.Debug().Str("server_temp", Endpoints.QueryServerTemp).Msg("Server Temp Endpoint Initialized")
}

type TempData struct {
	Sensor  string `json:"sensor"`
	Celsius any    `json:"celsius"`
}

// Helper function to get server Temperature information and format into HTTP response
func GenerateHostTempJSON() map[string][]TempData {
	// response := map[string]any{}
	response := map[string][]TempData{
		"data": []TempData{},
	}

	temps, err := host.SensorsTemperatures()
	if err != nil || len(temps) == 0 {
		log.Error().Str("server_temp", "Unavailable").Msg("Failed to get host.SensorsTemperatures()")
		newTempData := TempData{
			Sensor:  "Temperature Not Available",
			Celsius: nil,
		}
		response["data"] = append(response["data"], newTempData)
	}

	for _, t := range temps {
		sensor := fmt.Sprintf("sensor_%s", t.SensorKey)
		newTempData := TempData{
			Sensor:  sensor,
			Celsius: t.Temperature,
		}
		response["data"] = append(response["data"], newTempData)
	}

	return response
}
