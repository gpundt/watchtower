package main

import (
	"sync"
	"time"

	API "watchtower/internal/api"
	Config "watchtower/internal/config"

	Scanner "watchtower/internal/scanner"
	Logger "watchtower/pkg/logger"
)

func main() {
	// Initialize Configs and loggng
	Config.InitializeConfigWrapper(
		&Config.ServerConfig,
		Config.ServerPaths.ConfigFilepath,
	)
	Config.InitializeFilepaths(Config.ServerPaths)
	Logger.InitializeServerLogger()

	var wg sync.WaitGroup

	// Create background process to conduct network scans
	if Config.ServerConfig.Scanner.Enabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Scanner.StartNetworkScanner()
			ticker := time.NewTicker(
				time.Duration(Config.ServerConfig.Scanner.IntervalMinutes) * time.Minute,
			)
			defer ticker.Stop()
			for range ticker.C {
				Scanner.StartNetworkScanner()
			}
		}()
	}
	// Create background process to handle API endpoints
	go func() {
		defer wg.Done()
		API.InitializeServerAPI()
	}()

	wg.Wait()
	Logger.Close()
}
