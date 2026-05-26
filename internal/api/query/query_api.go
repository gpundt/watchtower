package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	Endpoints "watchtower/internal/api/endpoints"
	Config "watchtower/internal/config"

	"github.com/rs/zerolog/log"
)

func InitializeQueryAPI() {
	initializeTestEndpoint()
	initializeHealthCheckAPI()
}

func initializeTestEndpoint() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, encrypted world!"))
	})
	log.Debug().Str("test", "/").Msg("Test Endpoint: Initialized")
}

func initializeHealthCheckAPI() {
	http.HandleFunc(Endpoints.HealthCheckEndpoint, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		b, _ := json.MarshalIndent(generateHealthCheckResponse(), "", "  ")
		fmt.Fprintf(w, string(b))
	})
	log.Debug().Str("health_check", Endpoints.HealthCheckEndpoint).Msg("Health Check Endpoint: Initialized")
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