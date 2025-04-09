package config

import (
	"strconv"
)

type ServerConfig struct {
	Port int
}

func LoadServerConfig(envMap map[string]string) *ServerConfig {
	portString := envMap["HTTP_PORT"]
	if portString == "" {
		portString = "8080"
		WarningLog.Println("HTTP_PORT is empty, default set 8080")
	}
	port, err := strconv.Atoi(portString)
	if err != nil || port == 0 {
		if err != nil {
			WarningLog.Println("HTTP_PORT must be an integer, default set '8080'")
		} else {
			WarningLog.Println("HTTP_PORT can not be set '0', default set '8080'")
		}
		port = 8080
	}
	return &ServerConfig{
		Port: port,
	}
}
