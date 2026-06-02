package logs

type LogEntry struct {
	Host        string `json:"host"`
	Timestamp   string `json:"timestamp"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	User        string `json:"user"`
}

func ParseAuthLogs() []LogEntry {
	entries := []LogEntry{}

	return entries
}