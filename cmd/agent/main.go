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
	Config.InitializeConfigWrapper(
		&Config.AgentConfig,
		Config.AgentPaths.ConfigFilepath,
	)
	Config.PrepareFilepaths(Config.AgentPaths)
	Logger.InitializeAgentLogger()

	API.InitializeAgentAPI()

	ticker := time.NewTicker(
		time.Duration(Config.AgentConfig.Agent.PushIntervalSeconds) * time.Second,
	)
	defer ticker.Stop()

	for range ticker.C {
		Query.QueryHealthCheckEndpoint()
		Submission.SubmitHostMetricsMaster()
	}

	Logger.Close()
}
