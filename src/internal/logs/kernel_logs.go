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

var kernelUserRegex = regexp.MustCompile(`(?:user|uid|auid)=(\d+)`)

func ParseKernelLogs() []LogEntry {
	log.Info().Str("func", "ParseKernelLogs").Msg("Parsing Kernel Logs")
	entries := []LogEntry{}
	kernLogFilepath := "/var/log/kern.log"

	file, err := os.Open(kernLogFilepath)
	if err != nil {
		log.Err(err).Str("func", "ParseKernelLogs").Msg("")
		return entries
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		//Skip if the line doesnt match the standard pattern
		if !LogRegex.MatchString(line) {
			log.Debug().Str("line", line).Msg("no kernel regex match")
			continue
		}

		// Extract named submatches into a map
		match := LogRegex.FindStringSubmatch(line)
		result := make(map[string]string)
		for i, name := range LogRegex.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		// Parse timestamp
		currentYear := time.Now().Year()
		timeStr := fmt.Sprintf("%d %s", currentYear, result["timestamp"])
		parsedTime, err := time.Parse("2006 Jan  2 15:04:05", timeStr)
		if err != nil {
			parsedTime, _ = time.Parse("2006 Jan 2 15:04:05", timeStr)
		}

		// Extract user from log message
		name := "N/A"
		userMatches := kernelUserRegex.FindStringSubmatch(result["message"])
		if len(userMatches) > 1 {
			name = userMatches[1]
		}

		entry := LogEntry{
			Table:       Database.KernelLogEntryTable,
			Host:        Config.AgentConfig.Agent.Name,
			Timestamp:   parsedTime.UTC().Format("2006-01-02 15:04:05.000000-0700"),
			Severity:    "KERN",
			Message:     result["message"],
			ServiceName: strings.TrimSpace(result["service"]),
			User:        name,
		}

		entries = append(entries, entry)
	}

	log.Debug().Str("func", "ParseKernelLogs").Msg(
		fmt.Sprintf("Kernel Logs Parsed: %d", len(entries)),
	)
	return entries
}
