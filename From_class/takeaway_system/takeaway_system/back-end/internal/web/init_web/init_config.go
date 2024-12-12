package init_web

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
)

var configFilePath = "config.ini"

type Config struct {
	// database
	DbUsername string
	DbPassword string
	DbHost     string
	DbPort     int
	DbName     string
	DbTimeout  string
	// server
	ServerPort string
	// cors
	AllowHost string
	// cookie
	CookieName string
	// admin
	AdminUsername string
	AdminPassword string
	// other
	DefaultPassword string
}

func NewConfig() *Config { // 默认配置
	return &Config{
		// database
		DbUsername: "root",
		DbPassword: "1234",
		DbHost:     "127.0.0.1",
		DbPort:     3306,
		DbName:     "test",
		DbTimeout:  "10s",
		// server
		ServerPort: "1000",
		//cors
		AllowHost: "example.com",
		// cookie
		CookieName: "sid",
		// admin
		AdminUsername: "admin",
		AdminPassword: "123456",
		// other
		DefaultPassword: "123456",
	}
}

func InitConfig() *Config { // 加载配置文件
	defaultConfig := NewConfig()

	// 检查配置文件是否存在
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fmt.Println("配置文件不存在，创建默认配置文件...")
		if err := writeDefaultConfig(defaultConfig); err != nil {
			log.Fatalf("创建默认配置文件失败: %v", err)
		}
		fmt.Println("默认配置文件已创建。")
		return defaultConfig
	}

	// 加载现有配置文件
	cfg, err := ini.Load(configFilePath)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 检查配置是否完整
	if !isConfigComplete(cfg) {
		fmt.Println("配置文件字段不完整，重置为默认配置...")
		if err := writeDefaultConfig(defaultConfig); err != nil {
			log.Fatalf("重置配置文件失败: %v", err)
		}
		fmt.Println("配置文件已重置为默认配置。")
		return defaultConfig
	}

	// 配置文件存在且完整，读取配置
	config := readConfig(cfg, defaultConfig)
	return config
}

// writeDefaultConfig 创建或重置配置文件为默认配置
func writeDefaultConfig(config *Config) error {
	cfg := ini.Empty()

	// 设置数据库配置
	secDatabase := cfg.Section("database")
	secDatabase.Key("username").SetValue(config.DbUsername)
	secDatabase.Key("password").SetValue(config.DbPassword)
	secDatabase.Key("host").SetValue(config.DbHost)
	secDatabase.Key("port").SetValue(fmt.Sprintf("%d", config.DbPort))
	secDatabase.Key("name").SetValue(config.DbName)
	secDatabase.Key("timeout").SetValue(config.DbTimeout)

	// 设置服务器配置
	secServer := cfg.Section("server")
	secServer.Key("port").SetValue(config.ServerPort)

	// 设置cors配置
	secCors := cfg.Section("cors")
	secCors.Key("allow_host").SetValue(config.AllowHost)

	// 设置cookie配置
	secCookie := cfg.Section("cookie")
	secCookie.Key("name").SetValue(config.CookieName)

	// 设置管理员配置
	secAdmin := cfg.Section("admin")
	secAdmin.Key("username").SetValue(config.AdminUsername)
	secAdmin.Key("password").SetValue(config.AdminPassword)

	// 设置其他配置
	secOther := cfg.Section("other")
	secOther.Key("default_password").SetValue(config.DefaultPassword)
	// 保存到文件
	return cfg.SaveTo(configFilePath)
}

// isConfigComplete 检查配置文件是否包含所有必需的字段
func isConfigComplete(cfg *ini.File) bool {
	// 定义必需的配置结构
	requiredConfig := map[string][]string{
		"database": {"username", "password", "host", "port", "name", "timeout"},
		"server":   {"port"},
		"cors":     {"allow_host"},
		"cookie":   {"name"},
		"admin":    {"username", "password"},
		"other":    {"default_password"},
	}

	for section, keys := range requiredConfig {
		sec, err := cfg.GetSection(section)
		if err != nil {
			fmt.Printf("缺少段落 [%s]\n", section)
			return false
		}
		for _, key := range keys {
			if !sec.HasKey(key) {
				fmt.Printf("段落 [%s] 缺少键: %s\n", section, key)
				return false
			}
		}
	}
	return true
}

// readConfig 从配置文件中读取配置项并返回 Config 实例
func readConfig(cfg *ini.File, defaultConfig *Config) *Config {
	config := &Config{}
	// 读取数据库配置
	dbSection := cfg.Section("database")
	config.DbUsername = dbSection.Key("username").MustString(defaultConfig.DbUsername)
	config.DbPassword = dbSection.Key("password").MustString(defaultConfig.DbPassword)
	config.DbHost = dbSection.Key("host").MustString(defaultConfig.DbHost)
	config.DbPort = dbSection.Key("port").MustInt(defaultConfig.DbPort)
	config.DbName = dbSection.Key("name").MustString(defaultConfig.DbName)
	config.DbTimeout = dbSection.Key("timeout").MustString(defaultConfig.DbTimeout)

	// 读取服务器配置
	serverSection := cfg.Section("server")
	config.ServerPort = serverSection.Key("port").MustString(defaultConfig.ServerPort)

	// 读取cors配置
	corsSection := cfg.Section("cors")
	config.AllowHost = corsSection.Key("allow_host").MustString(defaultConfig.AllowHost)

	// 读取cookie配置
	cookieSection := cfg.Section("cookie")
	config.CookieName = cookieSection.Key("name").MustString(defaultConfig.CookieName)

	// 读取管理员配置
	adminSection := cfg.Section("admin")
	config.AdminUsername = adminSection.Key("username").MustString(defaultConfig.AdminUsername)
	config.AdminPassword = adminSection.Key("password").MustString(defaultConfig.AdminPassword)

	// 读取其他配置
	otherSection := cfg.Section("other")
	config.DefaultPassword = otherSection.Key("default_password").MustString(defaultConfig.DefaultPassword)

	return config
}
