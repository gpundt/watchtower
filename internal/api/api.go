package api

import (
	Query "watchtower/internal/api/query"
	TLS "watchtower/internal/api/tls"

	"github.com/rs/zerolog/log"
)

func InitializeServerAPI() {
	Query.InitializeQueryAPI()

	TLS.InitializeServermTLS()

	log.Info().Msg("Server API: Initialized")
}

func InitializeAgentAPI() {
	TLS.InitializeAgentmTLS()
	Query.QueryHealthCheckEndpoint()
}
