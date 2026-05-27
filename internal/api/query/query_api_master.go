package query

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Global function to start listening on API endpoint
func InitializeQueryAPI() {
	initializeHealthCheckAPIEndpoint()
	initializeServerMetricsAPIEndpoints()
}

// helper function to initialize an HTTP endpoint to accept GET requests
func makeGetHandler(
	responseGenerator func() map[string]any,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Debug().Str("url", r.URL.Path).
			Str("request_addr", r.RemoteAddr).
			Str("method", r.Method).
			Msg("")

		w.Header().Set("Content-Type", "application/json")
		b, _ := json.MarshalIndent(responseGenerator(), "", "  ")
		fmt.Fprintf(w, string(b))
	}
}
