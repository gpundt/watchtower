package config

import (
	"log"
	"os"
)

type Filepaths struct {
	ConfigFilepath  string
	LogDirectory    string
	LogFilepath     string
	CACertFilepath  string
	CertFilepath	string
	KeyFilepath		string
}

var AgentPaths = Filepaths{
	ConfigFilepath: "../config/agent.yaml",
	LogDirectory:   "../logs/",
	LogFilepath:    "",
	CACertFilepath: "../certs/ca/ca.crt",
	CertFilepath:   "../certs/agents/test_agent.crt",
	KeyFilepath:    "../certs/agents/test_agent.key",
}

var ServerPaths = Filepaths{
	ConfigFilepath: "../config/server.yaml",
	LogDirectory:   "../logs",
	LogFilepath:    "",
	CACertFilepath: "../certs/ca/ca.crt",
	CertFilepath: 	"../certs/server/server.crt",
	KeyFilepath:	"../certs/server/server.key",
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
