package query

import (
	"fmt"
	// "io"
	"net/http"

	Handlers "watchtower/internal/api/handlers"
	TLS "watchtower/internal/api/tls"
	Config "watchtower/internal/config"
	Endpoints "watchtower/pkg/endpoints"

	"github.com/rs/zerolog/log"
)

// ----- Server Health Check --------------------------------------------
func InitializeHealthCheckAPIEndpoint() {
	http.HandleFunc(
		Endpoints.QueryHealthCheck,
		Handlers.MakeGetHandler(generateHealthCheckResponse),
	)
	log.Debug().Str("health_check", Endpoints.QueryHealthCheck).
		Msg("Health Check Endpoint: Initialized")
}

// Generates JSON response for incoming health_check queries
func generateHealthCheckResponse() map[string]any {
	return map[string]any{
		"name":      "watchtower",
		"component": "server",
		"status":    "healthy",
		"uptime":    Config.GetUptimeString(),
	}
}

// Function to contact the server's health_check endpoint
func QueryHealthCheckEndpoint() error{
	resp, err := TLS.AgentTLSClient.Get(fmt.Sprintf(
		"%s%s",
		Config.AgentConfig.Agent.ServerURL,
		Endpoints.QueryHealthCheck,
	))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err := fmt.Errorf(
			"Request failed with status: %d",
			resp.StatusCode,
		)
		return err
	}
	return nil
}
