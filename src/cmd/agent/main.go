package main

import (
	"sync"
	"time"

	API "watchtower/internal/api"
	Query "watchtower/internal/api/query"
	Registration "watchtower/internal/api/registration"
	Submission "watchtower/internal/api/submission"
	Config "watchtower/internal/config"
	Endpoints "watchtower/pkg/endpoints"
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

	var wg sync.WaitGroup
	// Background process for submitting metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Create timer loop for metrics
		metricsTicker := time.NewTicker(
			time.Duration(Config.AgentConfig.Agent.PushIntervalSeconds) * time.Second,
		)
		defer metricsTicker.Stop()

		// Every Config.AgentConfig.Agent.PushIntervalSeconds
		for range metricsTicker.C {
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
	}()

	// Background process for submitting logs
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Create timer loop for metrics
		logsTicker := time.NewTicker(
			time.Duration(Config.AgentConfig.Agent.PushIntervalSeconds) * 20 * time.Second,
		)
		defer logsTicker.Stop()

		for range logsTicker.C {
			Submission.SubmitLogs()
			log.Info().Str("endpoint", Endpoints.SubmitLogEntry).
				Msg("Host Logs: Submitted")
		}
	}()

	wg.Wait()
	Logger.Close()
}
