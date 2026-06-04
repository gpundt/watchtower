package logs

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	Config "watchtower/internal/config"
	Database "watchtower/internal/database"

	"github.com/rs/zerolog/log"
)

type LogEntry struct {
	Table       string `json:"table"`
	Host        string `json:"host"`
	Timestamp   string `json:"timestamp"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	User        string `json:"user"`
}

// Regex to extract Timestamp, Service, and Message
var LogRegex = regexp.MustCompile(
	`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+[+-]\d{2}:\d{2})\s+` + // ISO 8601 timestamp
		`[a-zA-Z0-9_\.\-]+\s+` + // Hostname (ignored)
		`([^\[:\s][^:\[]*?)` + // Service name (non-greedy)
		`(?:\[\d+\])?:\s+` + // PID (optional)
		`(.*)$`, // Message
)

// Regex to find targeted user in message
var authUserRegex = regexp.MustCompile(`(?:user|for|to)\s+([a-zA-Z0-9_-]+)`)

func ParseAuthLogs() []LogEntry {
	log.Info().Str("func", "ParseAuthLogs").Msg("Parsing Auth Logs")
	entries := []LogEntry{}
	authLogFilepath := "/var/log/auth.log"

	file, err := os.Open(authLogFilepath)
	if err != nil {
		log.Err(err).Str("func", "ParseAuthLogs").Msg("")
		return entries
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := LogRegex.FindStringSubmatch(line)

		if len(matches) != 4 {
			log.Debug().Str("line", line).Msg("no auth regex match")
			continue
		}

		// Parse timestamp
		parsedTime, err := time.Parse(time.RFC3339Nano, matches[1])
		if err != nil {
			log.Warn().Str("func", "ParseAuthLogs").Str("timestamp", matches[1]).Msg("failed to parse timestamp")
			continue
		}

		// Extract user if present in message
		userMatch := authUserRegex.FindStringSubmatch(matches[3])
		user := ""
		if len(userMatch) > 1 {
			user = userMatch[1]
		}

		// Append
		entries = append(entries, LogEntry{
			Table:       Database.AuthLogEntryTable,
			Host:        Config.AgentConfig.Agent.Name,
			Timestamp:   parsedTime.UTC().Format("2006-01-02 15:04:05.000000-0700"),
			Severity:    "INFO",
			Message:     matches[3],
			ServiceName: strings.TrimSpace(matches[2]),
			User:        user,
		})
	}

	log.Debug().Str("func", "ParseAuthLogs").Msg(
		fmt.Sprintf("Auth Logs Parsed: %d", len(entries)),
	)
	return entries
}
