package api

import (
	Query "watchtower/internal/api/query"
	Submission "watchtower/internal/api/submission"
	TLS "watchtower/internal/api/tls"

	"github.com/rs/zerolog/log"
)

func InitializeServerAPI() {
	Query.InitializeQueryAPI()
	Submission.InitializeHostMetricsAPIEndpoints()

	TLS.InitializeServermTLS()

	log.Info().Msg("Server API: Initialized")
}

func InitializeAgentAPI() {
	TLS.InitializeAgentmTLS()
	Query.QueryHealthCheckEndpoint()
}
