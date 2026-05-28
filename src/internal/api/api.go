package api

import (
	Query "watchtower/internal/api/query"
	Registration "watchtower/internal/api/registration"
	Submission "watchtower/internal/api/submission"
	TLS "watchtower/internal/api/tls"

	"github.com/rs/zerolog/log"
)

func InitializeServerAPI() {
	Query.InitializeQueryEndpoints()
	Submission.InitializeSubmissionEndpoints()
	Registration.InitializeRegistrationEndpoints()

	TLS.InitializeServermTLS()

	log.Info().Msg("Server API: Initialized")
}

func InitializeAgentAPI() {
	TLS.InitializeAgentmTLS()
	Query.QueryHealthCheckEndpoint()
}
