package submission

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	Handlers "watchtower/internal/api/handlers"
	TLS "watchtower/internal/api/tls"
	Config "watchtower/internal/config"
	Database "watchtower/internal/database"
	Endpoints "watchtower/pkg/endpoints"

	"github.com/rs/zerolog/log"
)

type HostCheckInBody struct {
	Host string `json:"host"`
}

func SubmitHostCheckIn() {
	hostCheckInBody := HostCheckInBody{
		Host: Config.AgentConfig.Agent.Name,
	}
	jsonData, err := json.MarshalIndent(hostCheckInBody, "", "  ")
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf(
			"Failed to marshal agentRegistrationBody: %+v",
			hostCheckInBody,
		))
		return
	}

	// Make post request
	resp, err := TLS.AgentTLSClient.Post(
		fmt.Sprintf(
			"%s%s",
			Config.AgentConfig.Agent.ServerURL,
			Endpoints.SubmitHostCheckIn,
		),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf(
			"Failed to make POST to %s",
			Endpoints.SubmitHostCheckIn,
		))
		return
	}
	defer resp.Body.Close()
	log.Debug().Str("endpoint", Endpoints.SubmitHostCheckIn).
		Msg("Agent: Registered")
}

func initializeHostCheckInSubmissionEndpoint() {
	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.SubmitHostCheckIn),
		Handlers.MakePostHandler[HostCheckInBody](func(body HostCheckInBody) error {
			err := Database.UpdateAgent(
				body.Host,
				time.Now(),
			)
			if err != nil {
				log.Err(err).Str("endpoint", Endpoints.SubmitHostCheckIn).Msg("")
				return err
			}
			return nil
		}),
	)
	log.Debug().Str("host_check-in", Endpoints.SubmitHostCheckIn).
		Msg("Host Check-in Endpoint: Initialized")
}