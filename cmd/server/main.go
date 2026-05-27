package main

import (
	"sync"

	API "watchtower/internal/api"
	Config "watchtower/internal/config"
	Logger "watchtower/pkg/logger"
)

func main() {
	// Initialize Configs and loggng
	Config.InitializeConfigWrapper(
		&Config.ServerConfig,
		Config.ServerPaths.ConfigFilepath,
	)
	Logger.InitializeServerLogger()

	// Create background process to handle API endpoints
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		API.InitializeServerAPI()
	}()

	wg.Wait()
	Logger.Close()
}
