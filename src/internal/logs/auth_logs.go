package logs

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

func ParseAuthLogs() []LogEntry {
	entries := []LogEntry{}
	authLogFilepath := "/var/log/auth.log"

	file, err := os.Open(authLogFilepath)
	if err != nil {
		log.Err(err).Str("func", "ParseAuthLogs").Msg("")
		return entries
	}
	defer file.Close()

	// Regex to extract Timestamp, Service, and Message
	logRegex := regexp.MustCompile(`^([A-Za-z]{3}\s+\d+\s+\d{2}:\d{2}:\d{2})\s+\S+\s+([^:\[]+)(?:\[\d+\])?: (.*)$`)


	// Regex to find targeted user in message
	userRegex := regexp.MustCompile(`(?:user|for|to)\s+([a-zA-Z0-9_-]+)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := logRegex.FindStringSubmatch(line)

		if len(matches) == 4 {
			// Parse timestamp
			currentYear := time.Now().Year()
			timeStr := fmt.Sprintf("%d %s", currentYear, matches[1])
			parsedTime, _ := time.Parse("2006 Jan  2 15:04:05", timeStr)
		
			// Extract user if present in message
			userMatch := userRegex.FindStringSubmatch(matches[3])
			user := ""
			if len(userMatch) > 1 {
				user = userMatch[1]
			}

			// Append 
			entries = append(
				entries,
				LogEntry{
					Table: Database.AuthLogEntryTable,
					Host: Config.AgentConfig.Agent.Name,
					Timestamp: parsedTime.UTC().Format("2006-01-02 15:04:05.000000-0700"),
					Severity: "INFO",
					Message: matches[3],
					ServiceName: matches[2],
					User: user,
				},
			)
		}
	}

	log.Debug().Str("func", "ParseAuthLogs").Msg(
		fmt.Sprintf("Auth Logs Parsed: %d", len(entries)),
	)
	return entries
}