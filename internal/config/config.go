package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Structure of the server YAML config
var ServerConfig ServerConfigWrapper

type ServerConfigWrapper struct {
	Server struct {
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		Secret  string `yaml:"secret"`
		Verbose bool   `yaml:"verbose"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"ssl_mode"`
	} `yaml:"database"`
	Scanner struct {
		Enabled                 bool     `yaml:"enabled"`
		IntervalMinutes         int      `yaml:"interval_minutes"`
		MaxConcurrentScans		int		 `yaml:"max_concurrent_scans"`
		Ports                   []int    `yaml:"ports"`
	} `yaml:"scanner"`
	Alerts struct {
		Discord struct {
			Enabled    bool   `yaml:"enabled"`
			WebhookURL string `yaml:"webhook_url"`
		} `yaml:"discord"`
		Email struct {
			Enabled  bool   `yaml:"enabled"`
			SMTPHost string `yaml:"smtp_host"`
			SMTPPort int    `yaml:"smtp_port"`
			From     string `yaml:"from"`
			To       string `yaml:"to"`
		} `yaml:"email"`
	} `yaml:"alerts"`
}

// Structure of the agent YAML config
var AgentConfig AgentConfigWrapper

type AgentConfigWrapper struct {
	Agent struct {
		PushIntervalSeconds int    `yaml:"push_interval_seconds"`
		ServerURL           string `yaml:"server_url"`
		Name				string `yaml:"name"`
		Verbose             bool   `yaml:"verbose"`
	} `yaml:"agent"`
	TLS struct {
		CACert	string	`yaml:"ca_cert"`
		AgentCert	string	`yaml:"agent_cert"`
		AgentKey	string	`yaml:"agent_key"`
	} `yaml:"tls"`
}

type ConfigConstraint interface {
	AgentConfigWrapper | ServerConfigWrapper
}

func InitializeConfigWrapper[T ConfigConstraint](configWrapper *T, filepath string) {
	yamlData, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = yaml.Unmarshal(yamlData, configWrapper)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
