package query

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func InitializeQueryAPI() {
	initializeHealthCheckAPI()
	initializeServerMetricsAPIEndpoints()
}

// helper functin to initialize an HTTP endpoint to accept GET requests
func makeGetHandler(generate func() map[string]any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Debug().Str("url", r.URL.Path).Str("request_addr", r.RemoteAddr).Str("method", r.Method).Msg("")

		w.Header().Set("Content-Type", "application/json")
		b, _ := json.MarshalIndent(generate(), "", "  ")
		fmt.Fprintf(w, string(b))
	}
}
