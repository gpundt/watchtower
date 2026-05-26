package main

import (
	API "watchtower/internal/api"
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

	Logger.Close()
}