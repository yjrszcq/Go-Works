package config

import (
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	Timeout  string
}

func LoadDatabaseConfig(envMap map[string]string) *DatabaseConfig {
	host := envMap["DB_HOST"]
	if host == "" {
		host = "localhost"
		WarningLog.Println("DB_HOST is empty, default set 'localhost'")
	}
	portString := envMap["DB_PORT"]
	if portString == "" {
		portString = "3306"
		WarningLog.Println("DB_PORT is empty, default set 3306")
	}
	port, err := strconv.Atoi(portString)
	if err != nil || port == 0 {
		if err != nil {
			WarningLog.Println("DB_PORT must be an integer, default set '3306'")
		} else {
			WarningLog.Println("DB_PORT can not be set '0', default set '3306'")
		}
		port = 3306
	}
	user := envMap["DB_USER"]
	if user == "" {
		user = "root"
		WarningLog.Println("DB_USER is empty, default set 'root'")
	}
	password := envMap["DB_PASSWORD"]
	if password == "" {
		password = "1234"
		WarningLog.Println("DB_PASSWORD is empty, default set '1234'")
	}
	dbName := envMap["DB_NAME"]
	if dbName == "" {
		dbName = "test"
		WarningLog.Println("DB_NAME is empty, default set 'test'")
	}
	dbTimeout := envMap["DB_TIMEOUT"]
	if dbTimeout == "" {
		dbTimeout = "30s"
		WarningLog.Println("DB_TIMEOUT is empty, default set '30s'")
	}

	return &DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DbName:   dbName,
		Timeout:  dbTimeout,
	}
}
