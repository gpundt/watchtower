package config

type Filepaths struct {
	EtcDirectory	string
	ConfigFilepath  string
	OptDirectory	string
	BinaryDirectory string
	BinaryFilepath	string
	LogDirectory    string
	LogFilepath     string
	TLSDirectory	string
	CACertFilepath  string
	CertFilepath	string
	KeyFilepath		string
}

var AgentPaths = Filepaths{
	EtcDirectory: 	 "/etc/watchtower/",
	ConfigFilepath:  "/etc/watchtower/agent.yaml",
	OptDirectory:	 "/opt/watchtower/",
	BinaryDirectory: "/opt/watchtower/bin/",
	BinaryFilepath:  "/opt/watchtower/bin/watchtower_agent",
	LogDirectory:    "/var/log/watchtower/",
	LogFilepath:     "",
	TLSDirectory:    "/opt/watchtower/tls/",
	CACertFilepath:  "/opt/watchtower/tls/ca.crt",
	CertFilepath:    "/opt/watchtower/tls/test_agent.crt",
	KeyFilepath:     "/opt/watchtower/tls/test_agent.key",
}

var ServerPaths = Filepaths{
	EtcDirectory: 	 "/etc/watchtower/",
	ConfigFilepath:  "/etc/watchtower/server.yaml",
	OptDirectory:	 "/opt/watchtower/",
	BinaryDirectory: "/opt/watchtower/bin/",
	BinaryFilepath:  "/opt/watchtower/bin/watchtower_server",
	LogDirectory:    "/var/log/watchtower/",
	LogFilepath:     "",
	TLSDirectory:    "/opt/watchtower/tls/",
	CACertFilepath:  "/opt/watchtower/tls/ca.crt",
	CertFilepath:    "/opt/watchtower/tls/server.crt",
	KeyFilepath:     "/opt/watchtower/tls/server.key",
}
