package config

type CorsConfig struct {
	AllowHost string
}

func LoadCorsConfig(envMap map[string]string) *CorsConfig {
	allowHost := envMap["CORS_ALLOW_HOST"]
	if allowHost == "" {
		allowHost = "example.com"
		WarningLog.Println("CORS_ALLOW_HOST is empty, default set 'example.com'")
	}
	return &CorsConfig{
		AllowHost: allowHost,
	}
}
