package config

import (
	"log"
	"os"
)

var (
	InfoLog    = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	WarningLog = log.New(os.Stdout, "[WARNING] ", log.LstdFlags)
	ErrorLog   = log.New(os.Stdout, "[ERROR] ", log.LstdFlags)
)

type LoggerConfig struct {
	LogFile string
}

func LoadLoggerConfig(envMap map[string]string) *LoggerConfig {
	logFile := envMap["LOG_FILE"]
	if logFile == "" {
		logFile = "./log/app.log"
		WarningLog.Println("LOG_FILE is empty, default set './log/app.log'")
	}
	return &LoggerConfig{
		LogFile: logFile,
	}
}
