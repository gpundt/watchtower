package config

type Filepaths struct {
	EtcDirectory    string
	ConfigFilepath  string
	OptDirectory    string
	BinaryDirectory string
	BinaryFilepath  string
	LogDirectory    string
	LogFilepath     string
	CACertFilepath  string
	CertFilepath    string
	KeyFilepath     string
}

var AgentPaths = Filepaths{
	EtcDirectory:    "/etc/watchtower/",
	ConfigFilepath:  "/etc/watchtower/agent.yaml",
	OptDirectory:    "/opt/watchtower/",
	BinaryDirectory: "/opt/watchtower/bin/",
	BinaryFilepath:  "/opt/watchtower/bin/watchtower_agent",
	LogDirectory:    "/var/log/watchtower/",
	LogFilepath:     "",
}

var ServerPaths = Filepaths{
	EtcDirectory:    "/etc/watchtower/",
	ConfigFilepath:  "/etc/watchtower/server.yaml",
	OptDirectory:    "/opt/watchtower/",
	BinaryDirectory: "/opt/watchtower/bin/",
	BinaryFilepath:  "/opt/watchtower/bin/watchtower_server",
	LogDirectory:    "/var/log/watchtower/",
	LogFilepath:     "",
}
