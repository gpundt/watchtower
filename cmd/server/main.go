package main

import (
	"sync"

	API "watchtower/internal/api"
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		API.InitializeServerAPI()

	}()
	

	wg.Wait()
	Logger.Close()
}
