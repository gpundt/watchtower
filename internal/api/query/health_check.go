package query

import (
	"fmt"
	// "io"
	"net/http"

	Endpoints "watchtower/internal/api/endpoints"
	TLS "watchtower/internal/api/tls"
	Config "watchtower/internal/config"

	"github.com/rs/zerolog/log"
)

func initializeHealthCheckAPI() {
	http.HandleFunc(Endpoints.QueryHealthCheck, makeGetHandler(generateHealthCheckResponse))
	log.Debug().Str("health_check", Endpoints.QueryHealthCheck).Msg("Health Check Endpoint: Initialized")
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

func QueryHealthCheckEndpoint() {
	resp, err := TLS.AgentTLSClient.Get(fmt.Sprintf(
		"%s%s",
		Config.AgentConfig.Agent.ServerURL,
		Endpoints.QueryHealthCheck,
	))
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s unavailable", Endpoints.QueryHealthCheck))
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err := fmt.Errorf("Request failed with status: %d", resp.StatusCode)
		log.Err(err)
	}

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Err(err)
	// }
	
	// log.Debug().Str("health_check", Endpoints.QueryHealthCheck).Msg(string(body))
}
