package logger

import (
	"fmt"
	"os"
	"strconv"
	"time"

	Config "watchtower/internal/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var LogFile *os.File

// Root function to initialize server logger for file and STDOUT
func InitializeServerLogger() {
	if Config.ServerConfig.Server.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	initializeLogger("server", Config.ServerPaths)

	log.Debug().Str("Server Host", Config.ServerConfig.Server.Host).Msg("")
	log.Debug().Str("Server Port", strconv.Itoa(Config.ServerConfig.Server.Port)).Msg("")
	log.Debug().Str("Database Host", Config.ServerConfig.Database.Host).Msg("")
	log.Debug().Str("Database Port", strconv.Itoa(Config.ServerConfig.Database.Port)).Msg("")
	log.Debug().Str("Database Name", Config.ServerConfig.Database.Name).Msg("")
	log.Debug().Str("Database User", Config.ServerConfig.Database.User).Msg("")
}

// Root function to iniialize agent logger for file and STDOUT
func InitializeAgentLogger() {
	if Config.AgentConfig.Agent.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	initializeLogger("agent", Config.AgentPaths)

	log.Debug().Str("Server URL", Config.AgentConfig.Agent.ServerURL).Msg("")
	log.Debug().Str("Push Interval (sec)", strconv.Itoa(Config.AgentConfig.Agent.PushIntervalSeconds)).Msg("")
}

// Helper function to intialize log string Format
// Sets up file logger and stdout logger
func initializeLogger(component string, filepathsStruct Config.Filepaths) {
	zerolog.TimeFieldFormat = time.RFC3339

	date := time.Now().Format("2006-01-02")
	logFilename := fmt.Sprintf(
		"%s/%s_%s.log",
		filepathsStruct.LogDirectory,
		date,
		component,
	)
	filepathsStruct.LogFilepath = logFilename

	var err error
	LogFile, err = os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Str("path", logFilename).Msg("failed to open log file")
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
		NoColor:    false,
	}

	// JSON to file, pretty to stdout
	multi := zerolog.MultiLevelWriter(consoleWriter, LogFile)

	log.Logger = zerolog.New(multi).
		With().
		Timestamp().
		Logger()

	log.Info().Msg("Logger initialized")
}

// Call this on shutdown to cleanly close the log file
func Close() {
	if LogFile != nil {
		LogFile.Close()
	}
}
