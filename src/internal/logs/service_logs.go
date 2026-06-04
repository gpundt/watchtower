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
	`^([A-Z][a-z]{2}\s+\d+\s+\d{2}:\d{2}:\d{2})\s+` + // Timestamp (e.g., Jun  4 12:00:00)
		`([a-zA-Z0-9_\.\-]+)` + // Host (ignored in capture)
		`\s+([^\[:]+)` + // Service Name
		`(?:\[(\d+)\])?:\s+` + // PID (optional capture)
		`(.*)$`, // Message
)

func ParseServiceLogs() []LogEntry {
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
		if len(matches) < 6 {
			continue
		}

		currentYear := time.Now().Year()
		timeStr := fmt.Sprintf("%s %d", matches[1], currentYear)
		parsedTime, _ := time.Parse("2006 Jan  2 15:04:05", timeStr)

		entry := LogEntry{
			Table: Database.ServiceLogEntryTable,
			Host: Config.AgentConfig.Agent.Name,
			Timestamp: parsedTime.UTC().Format("2006-01-02 15:04:05.000000-0700"),
			Severity: "INFO",
			Message: matches[5],
			ServiceName: strings.TrimSpace(matches[3]),
			User: "N/A",
		}

		entries = append(entries, entry)
	}
	 
	log.Debug().Str("func", "ParseServiceLogs").Msg(
		fmt.Sprintf("Service Logs Parsed: %d", len(entries)),
	)
	return entries
}