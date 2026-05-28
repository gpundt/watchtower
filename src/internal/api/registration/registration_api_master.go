package registration

import (
	"fmt"
	"net/http"

	Database "watchtower/internal/database"
	Handlers "watchtower/internal/api/handlers"
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
	AgentHostname	string 		`json:"agent_hostname"`
}

// Function to intialize and handle incoming agent registration requests
func initializeAgentRegistrationEndpoint() {
	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.RegisterAgent),
		Handlers.MakePostHandler[AgentRegstrationBody](func(body AgentRegstrationBody) error {
			return Database.InsertAgentRegistration(body.AgentHostname)
		}),
	)
	log.Debug().Str("agent_registration", Endpoints.RegisterAgent).
		Msg("Agent Registration Endpoint: Initialized")
}
