package registration

import (
	"fmt"
	"net/http"

	Handlers "watchtower/internal/api/handlers"
	Endpoints "watchtower/pkg/endpoints"

	"github.com/rs/zerolog/log"
)

// Master Function to initialize registration endpoints
func InitializeRegistrationEndpoints() {
	initializeAgentRegistrationEndpoint()
	initializeUserRegistrationEndpoint()
}


// ----- Agent Registration -----------------------------------------------
// Struct to be populated with incoming agent registration request body
type AgentRegstrationBody struct {
	AgentIPAddr		string		`json:"agent_ip_addr"`
	AgentHostname	string 		`json:"agent_hostname"`
}

// Function to intialize and handle incoming agent registration requests
func initializeAgentRegistrationEndpoint() {
	var outputStruct AgentRegstrationBody
	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.RegisterAgent),
		Handlers.MakePostHandler[AgentRegstrationBody](outputStruct),
	)
	log.Debug().Str("agent_registration", Endpoints.RegisterAgent).
		Msg("Agent Registration Endpoint: Initialized")
}

// Function to be called on a successful agent registration request
func registerAgent() {
	return
}

// ----- User Registration ------------------------------------------------
// Struct to be populated with incoming user registration request body
type UserRegistrationBody struct {
	Username	string		`json:"username"`
	SourceIP	string		`json:"source_ip"`
}

// Function to initialize and handle incoming user registration requests
func initializeUserRegistrationEndpoint() {
	var outputStruct UserRegistrationBody
	http.HandleFunc(
		fmt.Sprintf("POST %s", Endpoints.RegisterUser),
		Handlers.MakePostHandler[UserRegistrationBody](outputStruct),
	)
	log.Debug().Str("user_registration", Endpoints.RegisterUser).
		Msg("User Registration Endpoint: Initialized")
}

// Function to be called on a successful user registration request
func registerUser() {
	return
}
