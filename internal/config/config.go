package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Structure of the server YAML config
type Server struct {
	Host	string	`yaml:"host"`
	Port	int		`yaml:"port"`
	Secret	string	`yaml:"secret"`
	Verbose bool	`yaml:"verbose"`
} `yaml:"server"`

// Structure of the database YAML config
type Database struct {
	Host	string	`yaml:"host"`
	Port	int		`yaml:"port"`
	Name	string	`yaml:"name"`
	User	string	`yaml:"user"`
	Password string	`yaml:"password"`
	SSLMode	string	`yaml:"ssl_mode"`
} `yaml:"database"`

// Structure of the scanner YAML config
type Scanner struct {
	Enabled		bool	`yaml:"enabled"`
	IntervalSeconds int `yaml:"interval_seconds"`
	PortScanIntervalSeconds int	`yaml:"port_scan_interval_seconds"`
	Targets []string	`yaml:"targets"`
	Ports	[]int		`yaml:"ports"`
} `yaml:"scanner"`

// Structure of the agent YAML config
type Agent struct {
	PushIntervalSeconds int	`yaml:"push_interval_seconds"`
	ServerURL string	`yaml:"server_url"`
	Token	string	`yaml:"token"`
	Verbose	bool	`yaml:"verbose"`
} `yaml:"agent"`

// Structure of the alerts YAML config
type Alerts struct {
	Discord struct {
		Enabled	bool	`yaml:"enabled"`
		WebhookURL	string	`yaml:"webhook_url"`
	} `yaml:"discord"`
	Email struct {
		Enabled bool 	`yaml:"enabled"`
		SMTPHost string `yaml:"smtp_host"`
		SMTPPort int	`yaml:"smtp_port"`
		From	string	`yaml:"from"`
		To		string  `yaml:"to"`
	} `yaml:"email"`
}

