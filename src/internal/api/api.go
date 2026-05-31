package api

import (
	Query "watchtower/internal/api/query"
	Registration "watchtower/internal/api/registration"
	Submission "watchtower/internal/api/submission"
	TLS "watchtower/internal/api/tls"

	"github.com/rs/zerolog/log"
)

func InitializeServerAPI() {
	// Initialize individual query endpoints
	Query.InitializeHealthCheckAPIEndpoint()
	Query.InitializeServerCPUEndpoint()
	Query.InitializeServerStorageEndpoint()
	Query.InitializeServerMemoryEndpoint()
	Query.InitializeServerTempEndpoint()

	// Initialize individual submission endpoints
	Submission.InitializeHostCheckInSubmissionEndpoint()
	Submission.InitializeHostCPUSubmissionEndpoint()
	Submission.InitializeHostMemorySubmissionEndpoint()
	Submission.InitializeHostStorageSubmissionEndpoint()
	Submission.InitializeHostTempSubmissionEndpoint()

	// Initialize individual registraion endpoints
	Registration.InitializeAgentRegistrationEndpoint()

	TLS.InitializeServermTLS()

	log.Info().Msg("Server API: Initialized")
}

func InitializeAgentAPI() {
	TLS.InitializeAgentmTLS()
	if err := Query.QueryHealthCheckEndpoint(); err != nil {
		log.Fatal().Err(err).Str("func", "Query.QueryHealthCheckEndpoint").Msg("")
	}
}