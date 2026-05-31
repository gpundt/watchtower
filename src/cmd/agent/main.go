package main

import (
	"time"

	API "watchtower/internal/api"
	Registration "watchtower/internal/api/registration"
	Endpoints "watchtower/pkg/endpoints"
	Query "watchtower/internal/api/query"
	Submission "watchtower/internal/api/submission"
	Config "watchtower/internal/config"
	Logger "watchtower/pkg/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize Configs, filepaths, and logging
	Config.InitializeConfigWrapper(
		&Config.AgentConfig,
		Config.AgentPaths.ConfigFilepath,
	)
	Config.InitializeFilepaths(Config.AgentPaths)
	Logger.InitializeAgentLogger()

	// Initialize mTLS and contact server registration endpoint
	API.InitializeAgentAPI()
	Registration.RegisterAgent()

	// Create timer loop
	ticker := time.NewTicker(
		time.Duration(Config.AgentConfig.Agent.PushIntervalSeconds) * time.Second,
	)
	defer ticker.Stop()

	// Every Config.AgentConfig.Agent.PushIntervalSeconds
	for range ticker.C {
		// Submit data to each individual submisson endpoint
		Submission.SubmitHostCheckIn()
		Submission.SubmitHostMetrics(
			Endpoints.SubmitHostCPU,
			Query.GenerateHostCPUJSON(
				Config.AgentConfig.Agent.Name,
			),
		)
		Submission.SubmitHostMetrics(
			Endpoints.SubmitHostMemory,
			Query.GenerateHostMemoryJSON(
				Config.AgentConfig.Agent.Name,
			),
		)
		Submission.SubmitHostMetrics(
			Endpoints.SubmitHostStorage,
			Query.GenerateHostStorageJSON(
				Config.AgentConfig.Agent.Name,
			),
		)
		Submission.SubmitHostMetrics(
			Endpoints.SubmitHostTemp,
			Query.GenerateHostTempJSON(
				Config.AgentConfig.Agent.Name,
			),
		)
		log.Info().Str("endpoint", Endpoints.SubmissionEndpoint).
			Msg("Host Metrics: Submitted")
	}

	Logger.Close()
}
