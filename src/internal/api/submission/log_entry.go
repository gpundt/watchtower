package submission

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	TLS "watchtower/internal/api/tls"
	Handlers "watchtower/internal/api/handlers"
	Config "watchtower/internal/config"
	Database "watchtower/internal/database"
	Endpoints "watchtower/pkg/endpoints"
	Logs "watchtower/internal/logs"

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
			timestampObj, parseErr := time.Parse("2006-01-02 15:04:05.000000-0700", body.Timestamp)
			if parseErr != nil {
				log.Err(parseErr).Str("func", "InitializeLogSubmissionEndpoint").Msg("")
			}
			err := Database.InsertLogEntry(
				timestampObj,
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

// Function to submit an individual log entry
func submitLog(entry Logs.LogEntry) {
	jsonData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf(
			"Failed to marshal LogEntry: %+v",
			entry,
		))
		return
	}
	// Make post request
	resp, err := TLS.AgentTLSClient.Post(
		fmt.Sprintf(
			"%s%s",
			Config.AgentConfig.Agent.ServerURL,
			Endpoints.SubmitLogEntry,
		),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf(
			"Failed to make POST to %s",
			Endpoints.SubmitLogEntry,
		))
		return
	}
	defer resp.Body.Close()
}

// Function to call every parse log function and submit parsed logs.
func SubmitLogs() {
	// --- Auth logs ---------------------------------------------
	for _, authLog := range Logs.ParseAuthLogs() {
		submitLog(authLog)
	}
	log.Info().Str("endpoint", Endpoints.SubmitLogEntry).
		Msg("Auth Logs: Submitted")

	// --- Container logs -----------------------------------------
	for _, containerLog := range Logs.ParseContainerLogs() {
		submitLog(containerLog)
	}
	log.Info().Str("endpoint", Endpoints.SubmitLogEntry).
		Msg("Container Logs: Submitted")

	// --- Kernel logs -------------------------------------------
	for _, kernelLog := range Logs.ParseKernelLogs() {
		submitLog(kernelLog)
	}
	log.Info().Str("endpoint", Endpoints.SubmitLogEntry).
		Msg("Kernel Logs: Submitted")

	// --- Scheduled Task logs -----------------------------------
	for _, scheduledTaskLog := range Logs.ParseScheduledTaskLogs() {
		submitLog(scheduledTaskLog)
	}
	log.Info().Str("endpoint", Endpoints.SubmitLogEntry).
		Msg("Scheduled Task Logs: Submitted")

	// --- Service logs ------------------------------------------
	for _, serviceLog := range Logs.ParseServiceLogs() {
		submitLog(serviceLog)
	}
	log.Info().Str("endpoint", Endpoints.SubmitLogEntry).
		Msg("Service Logs: Submitted")
}
