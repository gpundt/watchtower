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

var serviceLogRegex = regexp.MustCompile(
	`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+[+-]\d{2}:\d{2})\s+` + // ISO 8601 timestamp
		`[a-zA-Z0-9_\.\-]+\s+` + // Hostname (ignored)
		`([^\[:\s][^:\[]*?)` + // Service name (non-greedy)
		`(?:\[\d+\])?:\s+` + // PID (optional)
		`(.*)$`, // Message
)

func ParseServiceLogs() []LogEntry {
	log.Info().Str("func", "ParseServiceLogs").Msg("Parsing Service Logs")
	entries := []LogEntry{}
	serviceLogFilepath := "/var/log/syslog"

	file, err := os.Open(serviceLogFilepath)
	if err != nil {
		log.Err(err).Str("func", "ParseServiceLogs").Msg("")
		return entries
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := serviceLogRegex.FindStringSubmatch(line)
		if len(matches) != 4 {
			log.Debug().Str("line", line).Msg("no service regex match")
			continue
		}

		parsedTime, err := time.Parse(time.RFC3339Nano, matches[1])
		if err != nil {
			log.Warn().Str("func", "ParseServiceLogs").Str("timestamp", matches[1]).Msg("failed to parse timestamp")
			continue
		}

		entries = append(entries, LogEntry{
			Table:       Database.ServiceLogEntryTable,
			Host:        Config.AgentConfig.Agent.Name,
			Timestamp:   parsedTime.UTC().Format("2006-01-02 15:04:05.000000-0700"),
			Severity:    "INFO",
			Message:     matches[3],
			ServiceName: strings.TrimSpace(matches[2]),
			User:        "N/A",
		})
	}

	log.Debug().Str("func", "ParseServiceLogs").Msg(
		fmt.Sprintf("Service Logs Parsed: %d", len(entries)),
	)
	return entries
}
