package main

import (
	Config "watchtower/internal/config"
	Logger "watchtower/pkg/logger"
)

func main() {
	Config.InitializeConfigWrapper(
		&Config.ServerConfig,
		Config.ServerPaths.ConfigFilepath,
	)
	Config.PrepareFilepaths(Config.ServerPaths)
	Logger.InitializeServerLogger()

	Logger.Close()
}
