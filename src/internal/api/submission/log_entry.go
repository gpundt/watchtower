package submission

import (
	"fmt"
	"net/http"
	"time"

	Handlers "watchtower/internal/api/handlers"
	Database "watchtower/internal/database"
	Endpoints "watchtower/pkg/endpoints"

	"github.com/rs/zerolog/log"
)

type LogEntry struct {
	Host        string `json:"host"`
	Timestamp   string `json:"timestamp"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	User        string `json:"user"`
}

func InitializeLogSubmissionEndpoint() {
	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.SubmitLogEntry),
		Handlers.MakePostHandler[LogEntry](func(body LogEntry) error {
			err := Database.InsertLogEntry(
				time.Now(),
				body.Host,
				body.Severity,
				body.Message,
				body.ServiceName,
				body.User,
			)
			if err != nil {
				log.Err(err).Str("endpoint", Endpoints.SubmitLogEntry).Msg("")
				return err
			}
			return nil
		}),
	)
	log.Debug().Str("log_entry", Endpoints.SubmitLogEntry).
		Msg("Log Entry Submission Endpoint: Initialized")
}
