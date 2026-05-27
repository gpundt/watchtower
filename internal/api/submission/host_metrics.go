package submission

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	Endpoints "watchtower/internal/api/endpoints"
	Query "watchtower/internal/api/query"
	TLS "watchtower/internal/api/tls"
	Config "watchtower/internal/config"

	"github.com/rs/zerolog/log"
)

type MetricsStructConstraint interface {
	map[string]any | map[string][]Query.TempData
}

// Agent-side function to submit gathered host metrics
func submitHostMetrics[T MetricsStructConstraint](endpoint string, metricsStruct T) {
	jsonData, err := json.MarshalIndent(metricsStruct, "", "  ")
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf(
			"Failed to marshal metricsStrict: %+v",
			metricsStruct,
		))
		return
	}

	resp, err := TLS.AgentTLSClient.Post(
		fmt.Sprintf("%s%s", Config.AgentConfig.Agent.ServerURL, endpoint),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf(
			"Failed to make POST to %s",
			endpoint,
		))
		return
	}
	defer resp.Body.Close()
}

// Struct to be populated with incoming host CPU metrics submission
type HostCPUBody struct {
	Model          string  `json:"cpu_model"`
	Family         string  `json:"cpu_family"`
	ModelName      string  `json:"cpu_model_name"`
	Mhz            float64 `json:"cpu_mhz"`
	CacheSize      int32   `json:"cpu_cache_size"`
	LogicalCores   int     `json:"cpu_logical_cores"`
	PhysicalCores  int     `json:"cpu_physical_cores"`
	UsedPercentage float64 `json:"cpu_used_percentage"`
}

// Initializes server-side endpoint to receive CPU metrics
func initializeHostCPUSubmissionEndpoint() {
	var outputStruct HostCPUBody
	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.SubmitHostCPU),
		makePostHandler(outputStruct),
	)
	log.Debug().Str("host_cpu", Endpoints.SubmitHostCPU).Msg("Host CPU Submission Endpoint: Initialized")
}

// Struct to be populated with Incoming host memory metrics submission
type HostMemoryBody struct {
	TotalMemoryBytes     float64 `json:"total_memory_bytes"`
	FreeMemoryBytes      float64 `json:"free_memory_bytes"`
	FreeMemoryPercentage float64 `json:"free_memory_percentage"`
	UsedMemoryBytes      float64 `json:"used_memory_bytes"`
	UsedMemoryPercentage float64 `json:"used_memory_percentage"`
}

// Initializes server-side endpoint to receive memory metrics
func initializeHostMemorySubmissionEndpoint() {
	var outputStruct HostMemoryBody
	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.SubmitHostMemory),
		makePostHandler(outputStruct),
	)
	log.Debug().Str("host_cpu", Endpoints.SubmitHostMemory).Msg("Host Memory Submission Endpoint: Initialized")
}

// Struct to be populated with Incoming host storage metrics submission
type HostStorageBody struct {
	TotalStoageBytes      uint64 `json:"total_storage_bytes"`
	FreeStorageBytes      uint64 `json:"free_storage_bytes"`
	FreeStoragePercentage string `json:"free_storage_percentage"`
	UsedStorageBytes      uint64 `json:"used_storage_bytes"`
	UsedStoragePercentage string `json:"used_storage_percentage"`
}

// Initializes server-side endpoint to receive storage metrics
func initializeHostStorageSubmissionEndpoint() {
	var outputStruct HostStorageBody
	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.SubmitHostStorage),
		makePostHandler(outputStruct),
	)
	log.Debug().Str("host_cpu", Endpoints.SubmitHostStorage).Msg("Host Storage Submission Endpoint: Initialized")

}

// Struct to be populated with Incoming host Temperature metrics submission
type HostTempBody struct {
	Data []Query.TempData `json:"data"`
}

// Initializes server-side endpoint to receive temperature metrics
func initializeHostTempSubmissionEndpoint() {
	var outputStruct HostTempBody
	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.SubmitHostTemp),
		makePostHandler(outputStruct),
	)
	log.Debug().Str("host_cpu", Endpoints.SubmitHostTemp).Msg("Host Temp Submission Endpoint: Initialized")
}

