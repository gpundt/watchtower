package main

import (
	"time"

	API "watchtower/internal/api"
	Registration "watchtower/internal/api/registration"
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
	Config.InitializeFilepaths(Config.AgentPaths)
	Logger.InitializeAgentLogger()

	// Initialize API
	API.InitializeAgentAPI()
	Registration.RegisterAgent()

	// Create timer
	ticker := time.NewTicker(
		time.Duration(Config.AgentConfig.Agent.PushIntervalSeconds) * time.Second,
	)
	defer ticker.Stop()

	// Every Config.AgentConfig.Agent.PushIntervalSeconds
	for range ticker.C {
		Submission.SubmitHostCheckIn()
		Submission.SubmitAllHostMetrics()
	}

	Logger.Close()
}
