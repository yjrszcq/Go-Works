package lib

import (
	"github.com/joho/godotenv"
	"project/lib/config"
)

// 服务端配置数据结构
type Config struct {
	HttpPort   int
	AllowHost  string
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     int
	DbName     string
	DbTimeout  string
	LogFile    string
}

// 加载服务端配置
func LoadConfig() *Config {
	envMap, err := godotenv.Read("./conf/app.env")
	if err != nil {
		ErrorLog.Fatal(err)
	}
	server := config.LoadServerConfig(envMap)
	cors := config.LoadCorsConfig(envMap)
	database := config.LoadDatabaseConfig(envMap)
	logger := config.LoadLoggerConfig(envMap)
	Config := Config{
		HttpPort:   server.Port,
		AllowHost:  cors.AllowHost,
		DbUser:     database.User,
		DbPassword: database.Password,
		DbHost:     database.Host,
		DbPort:     database.Port,
		DbName:     database.DbName,
		DbTimeout:  database.Timeout,
		LogFile:    logger.LogFile,
	}
	InfoLog.Println("配置加载成功")
	return &Config
}
