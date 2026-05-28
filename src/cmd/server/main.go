package main

import (
	"sync"

	API "watchtower/internal/api"
	Config "watchtower/internal/config"
	Logger "watchtower/pkg/logger"
	Scanner "watchtower/internal/scanner"
)

func main() {
	// Initialize Configs and loggng
	Config.InitializeConfigWrapper(
		&Config.ServerConfig,
		Config.ServerPaths.ConfigFilepath,
	)
	Logger.InitializeServerLogger()


	var wg sync.WaitGroup

	// Create background process to conduct network scans
	if Config.ServerConfig.Scanner.Enabled {
		wg.Add(2)
		go func() {
			defer wg.Done()
			Scanner.InitializeNetworkScanner()
		}()
	} else {
		wg.Add(1)
	}
	
	// Create background process to handle API endpoints
	go func() {
		defer wg.Done()
		
		API.InitializeServerAPI()
	}()
	wg.Wait()
	Logger.Close()
}
