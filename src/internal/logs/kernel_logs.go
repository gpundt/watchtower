package logs

import (
	"fmt"
	"os"
	"bufio"
	"regexp"
	"time"

	Config "watchtower/internal/config"
	Database "watchtower/internal/database"

	"github.com/rs/zerolog/log"
)

func ParseKernelLogs() []LogEntry {
	entries := []LogEntry{}
	kernLogFilepath := "/var/log/kern.log"

	file, err := os.Open(kernLogFilepath)
	if err != nil {
		log.Err(err).Str("func", "ParseKernelLogs").Msg("")
		return entries
	}
	defer file.Close()

	logRegex := regexp.MustCompile(
		`^(?P<timestamp>[A-Z][a-z]{2}\s+\d+\s+\d{2}:\d{2}:\d{2})\s+` + // Timestamp
			`[a-zA-Z0-9_\-]+\s+` +                                  // Hostname (ignored)
			`(?P<service>[a-zA-Z0-9_\-]+)(?:\[\d+\])?:?\s+` +      // Service (e.g., kernel)
			`(?P<message>.*)$`, 
	)
	
	userRegex := regexp.MustCompile(`(?:user|uid|auid)=(\d+)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		//Skip if the line doesnt match the standard pattern
		if !logRegex.MatchString(line) {
			continue
		}

		// Extract named submatches into a map
		match := logRegex.FindStringSubmatch(line)
		result := make(map[string]string)
		for i, name := range logRegex.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		// Parse timestamp
		currentYear := time.Now().Year()
		timeStr := fmt.Sprintf("%d %s", currentYear, result["timestamp"])
		parsedTime, _ := time.Parse("2006 Jan  2 15:04:05", timeStr)

		// Extract user from log message
		name := "N/A"
		userMatches := userRegex.FindStringSubmatch(result["message"])
		if len(userMatches) > 1 {
			name = userMatches[1]
		}

		entry := LogEntry{
			Table: Database.KernelLogEntryTable,
			Host: Config.AgentConfig.Agent.Name,
			Timestamp: parsedTime.UTC().Format("2006-01-02 15:04:05.000000-0700"),
			Severity: "KERN",
			Message: result["message"],
			ServiceName: result["service"],
			User: name,
		}

		entries = append(entries, entry)
	}

	log.Debug().Str("func", "ParseKernelLogs").Msg(
		fmt.Sprintf("Kernel Logs Parsed: %d", len(entries)),
	)
	return entries
}