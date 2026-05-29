package registration

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

// Master Function to initialize registration endpoints
func InitializeRegistrationEndpoints() {
	initializeAgentRegistrationEndpoint()
}

// ----- Agent Registration -----------------------------------------------
// Struct to be populated with incoming agent registration request body
type AgentRegstrationBody struct {
	Host string `json:"host"`
}

// Function to intialize and handle incoming agent registration requests
func initializeAgentRegistrationEndpoint() {
	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.RegisterAgent),
		Handlers.MakePostHandler[AgentRegstrationBody](func(body AgentRegstrationBody) error {
			err := Database.InsertAgent(
				body.Host,
				time.Now(),
			)
			if err != nil {
				log.Err(err).Str("endpoint", Endpoints.RegisterAgent).Msg("")
				return err
			}
			return nil
		}),
	)
	log.Debug().Str("agent_registration", Endpoints.RegisterAgent).
		Msg("Agent Registration Endpoint: Initialized")
}

func RegisterAgent() {
	agentRegistrationBody := AgentRegstrationBody{
		Host: Config.AgentConfig.Agent.Name,
	}
	jsonData, err := json.MarshalIndent(agentRegistrationBody, "", "  ")
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf(
			"Failed to marshal agentRegistrationBody: %+v",
			agentRegistrationBody,
		))
		return
	}

	// Make post request
	resp, err := TLS.AgentTLSClient.Post(
		fmt.Sprintf(
			"%s%s",
			Config.AgentConfig.Agent.ServerURL,
			Endpoints.RegisterAgent,
		),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf(
			"Failed to make POST to %s",
			Endpoints.RegisterAgent,
		))
		return
	}
	defer resp.Body.Close()
	log.Debug().Str("endpoint", Endpoints.RegisterAgent).
		Msg("Agent: Registered")
}
