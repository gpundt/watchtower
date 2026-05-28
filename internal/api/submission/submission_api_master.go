package submission

import (
	"encoding/json"
	"fmt"
	"net/http"

	Query "watchtower/internal/api/query"
	Endpoints "watchtower/pkg/endpoints"

	"github.com/rs/zerolog/log"
)

// Submits gathered host metrics to server-side endpoints
func SubmitHostMetricsMaster() {
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
func InitializeHostMetricsAPIEndpoints() {
	initializeHostCPUSubmissionEndpoint()
	initializeHostMemorySubmissionEndpoint()
	initializeHostStorageSubmissionEndpoint()
	initializeHostTempSubmissionEndpoint()
}

// Paramater type constraint for postHandler outputStruct
type OutputStructConstraint interface {
	HostCPUBody | HostMemoryBody | HostStorageBody | HostTempBody
}

// Master function to create a POST endpoint that populates an outputStruct
func makePostHandler[T OutputStructConstraint](
	outputStruct T,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&outputStruct); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Debug().Str("url", r.URL.Path).
			Str("request_addr", r.RemoteAddr).
			Str("method", r.Method).
			Msg(fmt.Sprintf("%+v", outputStruct))
	}
}
