package main

import (
	"time"

	API "watchtower/internal/api"
	Query "watchtower/internal/api/query"
	Submission "watchtower/internal/api/submission"
	Config "watchtower/internal/config"
	Logger "watchtower/pkg/logger"
)

func main() {
	// Initialize Configs and logging
	Config.InitializeConfigWrapper(
		&Config.AgentConfig,
		Config.AgentPaths.ConfigFilepath,
	)
	Logger.InitializeAgentLogger()

	// Initialize API
	API.InitializeAgentAPI()

	// Create timer
	ticker := time.NewTicker(
		time.Duration(Config.AgentConfig.Agent.PushIntervalSeconds) * time.Second,
	)
	defer ticker.Stop()

	// Every Config.AgentConfig.Agent.PushIntervalSeconds
	for range ticker.C {
		Query.QueryHealthCheckEndpoint()
		Submission.SubmitHostMetricsMaster()
	}

	Logger.Close()
}
