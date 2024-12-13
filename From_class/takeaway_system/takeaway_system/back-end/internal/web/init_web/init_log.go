package init_web

import (
	"back-end/internal/web/web_log"
	"log"
	"os"
)

func initLog(cfg *Config) *web_log.WebLog {
	logger := log.New(os.Stdout, "WARNING: ", log.LstdFlags)
	if cfg.LogFile != "" {
		file, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.Println("打开日志文件成功")
			return &web_log.WebLog{
				InfoLogger:  log.New(file, "INFO: ", log.LstdFlags),
				ErrorLogger: log.New(file, "ERROR: ", log.LstdFlags),
			}
		}
		logger.Println("打开日志文件失败: %v，使用标准输出", err)
	} else {
		logger.Println("未配置日志文件，使用标准输出")
	}
	return &web_log.WebLog{
		InfoLogger:  log.New(os.Stdout, "INFO: ", log.LstdFlags),
		ErrorLogger: log.New(os.Stdout, "ERROR: ", log.LstdFlags),
	}
}
