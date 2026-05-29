package submission

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	Handlers "watchtower/internal/api/handlers"
	Query "watchtower/internal/api/query"
	TLS "watchtower/internal/api/tls"
	Config "watchtower/internal/config"
	Database "watchtower/internal/database"
	Endpoints "watchtower/pkg/endpoints"

	"github.com/rs/zerolog/log"
)

// ----- Metrics Submission ------------------------------------------------
type MetricsStructConstraint interface {
	map[string]any | Query.TemperatureResponse
}

// Agent-side function to submit gathered host metrics
func submitHostMetrics[T MetricsStructConstraint](
	endpoint string,
	metricsStruct T,
) {
	// Construct JSON body
	jsonData, err := json.MarshalIndent(metricsStruct, "", "  ")
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf(
			"Failed to marshal metricsStruct: %+v",
			metricsStruct,
		))
		return
	}

	// Make post request
	resp, err := TLS.AgentTLSClient.Post(
		fmt.Sprintf(
			"%s%s",
			Config.AgentConfig.Agent.ServerURL,
			endpoint,
		),
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
	log.Debug().Str("endpoint", endpoint).
		Msg("Metrics: Submitted")
}

// ----- Host CPU Metrics ------------------------------------------------
// Struct to be populated with incoming host CPU metrics submission
type HostCPUBody struct {
	Host           string  `json:"host"`
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
	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.SubmitHostCPU),
		Handlers.MakePostHandler[HostCPUBody](func(body HostCPUBody) error {
			err := Database.InsertHostCPUUsage(
				time.Now(),
				body.Host,
				body.UsedPercentage,
			)
			if err != nil {
				log.Err(err).Str("endpoint", Endpoints.SubmitHostCPU).Msg("")
				return err
			}
			return nil
		}),
	)
	log.Debug().Str("host_cpu", Endpoints.SubmitHostCPU).
		Msg("Host CPU Submission Endpoint: Initialized")
}

// ----- Host Memory Metrics ----------------------------------------------
// Struct to be populated with Incoming host memory metrics submission
type HostMemoryBody struct {
	Host                 string  `json:"host"`
	TotalMemoryBytes     float64 `json:"total_memory_bytes"`
	FreeMemoryBytes      float64 `json:"free_memory_bytes"`
	FreeMemoryPercentage float64 `json:"free_memory_percentage"`
	UsedMemoryBytes      float64 `json:"used_memory_bytes"`
	UsedMemoryPercentage float64 `json:"used_memory_percentage"`
}

// Initializes server-side endpoint to receive memory metrics
func initializeHostMemorySubmissionEndpoint() {

	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.SubmitHostMemory),
		Handlers.MakePostHandler[HostMemoryBody](func(body HostMemoryBody) error {
			err := Database.InsertHostMemoryUsage(
				time.Now(),
				body.Host,
				body.TotalMemoryBytes,
				body.FreeMemoryBytes,
				body.UsedMemoryBytes,
				body.FreeMemoryPercentage,
				body.UsedMemoryPercentage,
			)
			if err != nil {
				log.Err(err).Str("endpoint", Endpoints.SubmitHostMemory).Msg("")
				return err
			}
			return nil
		}),
	)
	log.Debug().Str("host_cpu", Endpoints.SubmitHostMemory).
		Msg("Host Memory Submission Endpoint: Initialized")
}

// ----- Host Storage Metrics ---------------------------------------------
// Struct to be populated with Incoming host storage metrics submission
type HostStorageBody struct {
	Host                  string  `json:"host"`
	TotalStorageBytes     uint64  `json:"total_storage_bytes"`
	FreeStorageBytes      uint64  `json:"free_storage_bytes"`
	FreeStoragePercentage float64 `json:"free_storage_percentage"`
	UsedStorageBytes      uint64  `json:"used_storage_bytes"`
	UsedStoragePercentage float64 `json:"used_storage_percentage"`
}

// Initializes server-side endpoint to receive storage metrics
func initializeHostStorageSubmissionEndpoint() {
	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.SubmitHostStorage),
		Handlers.MakePostHandler[HostStorageBody](func(body HostStorageBody) error {
			err := Database.InsertHostStorageUsage(
				time.Now(),
				body.Host,
				body.TotalStorageBytes,
				body.FreeStorageBytes,
				body.UsedStorageBytes,
				body.FreeStoragePercentage,
				body.UsedStoragePercentage,
			)
			if err != nil {
				log.Err(err).Str("endpoint", Endpoints.SubmitHostStorage).Msg("")
				return err
			}
			return nil
		}),
	)
	log.Debug().Str("host_cpu", Endpoints.SubmitHostStorage).
		Msg("Host Storage Submission Endpoint: Initialized")

}

// ----- Host Temperature Metrics------------------------------------------
// Struct to be populated with Incoming host Temperature metrics submission
type HostTempBody struct {
	Host string             `json:"host"`
	Data []Query.SensorData `json:"data"`
}

// Initializes server-side endpoint to receive temperature metrics
func initializeHostTempSubmissionEndpoint() {
	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.SubmitHostTemp),
		Handlers.MakePostHandler[HostTempBody](func(body HostTempBody) error {
			err := Database.InsertHostTemperature(
				time.Now(),
				body.Host,
				body.Data,
			)
			if err != nil {
				log.Err(err).Str("endpoint", Endpoints.SubmitHostTemp).Msg("")
				return err
			}
			return nil
		}),
	)
	log.Debug().Str("host_cpu", Endpoints.SubmitHostTemp).
		Msg("Host Temp Submission Endpoint: Initialized")
}
