package submission

import (
	Query "watchtower/internal/api/query"
	Config "watchtower/internal/config"
	Endpoints "watchtower/pkg/endpoints"

	"github.com/rs/zerolog/log"
)

// Submits gathered host metrics to server-side endpoints
func SubmitAllHostMetrics() {
	submitHostMetrics(
		Endpoints.SubmitHostCPU,
		Query.GenerateHostCPUJSON(
			Config.AgentConfig.Agent.Name,
		),
	)
	submitHostMetrics(
		Endpoints.SubmitHostMemory,
		Query.GenerateHostMemoryJSON(
			Config.AgentConfig.Agent.Name,
		),
	)
	submitHostMetrics(
		Endpoints.SubmitHostStorage,
		Query.GenerateHostStorageJSON(
			Config.AgentConfig.Agent.Name,
		),
	)
	submitHostMetrics(
		Endpoints.SubmitHostTemp,
		Query.GenerateHostTempJSON(
			Config.AgentConfig.Agent.Name,
		),
	)
	log.Info().Str("endpoint", Endpoints.SubmissionEndpoint).
		Msg("Host Metrics: Submitted")
}

// Initializes Server-side endpoints to receive host metrics
func InitializeSubmissionEndpoints() {
	initializeHostCheckInSubmissionEndpoint()
	initializeHostCPUSubmissionEndpoint()
	initializeHostMemorySubmissionEndpoint()
	initializeHostStorageSubmissionEndpoint()
	initializeHostTempSubmissionEndpoint()
}
