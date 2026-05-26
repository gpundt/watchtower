package config

import (
	"log"
	"os"
)

type Filepaths struct {
	ConfigFilepath string
	LogDirectory   string
	LogFilepath    string
}

var AgentPaths = Filepaths{
	ConfigFilepath: "../config/agent.yaml",
	LogDirectory:   "../logs/",
	LogFilepath:    "",
}

var ServerPaths = Filepaths{
	ConfigFilepath: "../config/server.yaml",
	LogDirectory:   "../logs",
	LogFilepath:    "",
}

func PrepareFilepaths(filepathsStruct Filepaths) {
	dirs := []string{
		filepathsStruct.LogDirectory,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf(
				"Filepath prep failed '%s': %v",
				dir,
				err,
			)
		}
	}
}
