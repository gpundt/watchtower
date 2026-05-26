package main

import (
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

	Logger.Close()
}