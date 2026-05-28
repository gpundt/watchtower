package submission

import (
	Query "watchtower/internal/api/query"
	Endpoints "watchtower/pkg/endpoints"

	"github.com/rs/zerolog/log"
)

// Submits gathered host metrics to server-side endpoints
func SubmitAllHostMetrics() {
	submitHostMetrics(
		Endpoints.SubmitHostCPU,
		Query.GenerateHostCPUJSON(),
	)
	submitHostMetrics(
		Endpoints.SubmitHostMemory,
		Query.GenerateHostMemoryJSON(),
	)
	submitHostMetrics(
		Endpoints.SubmitHostStorage,
		Query.GenerateHostStorageJSON(),
	)
	submitHostMetrics(
		Endpoints.SubmitHostTemp,
		Query.GenerateHostTempJSON(),
	)
	log.Info().Str("endpoint", Endpoints.SubmissionEndpoint).
		Msg("Host Metrics: Submitted")
}

// Initializes Server-side endpoints to receive host metrics
func InitializeSubmissionEndpoints() {
	initializeHostCPUSubmissionEndpoint()
	initializeHostMemorySubmissionEndpoint()
	initializeHostStorageSubmissionEndpoint()
	initializeHostTempSubmissionEndpoint()
}
