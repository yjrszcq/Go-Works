package lib

import (
	"log"
	"os"
	"path/filepath"
	"project/lib/config"
)

var (
	Logger     *WebLog
	InfoLog    = config.InfoLog
	WarningLog = config.WarningLog
	ErrorLog   = config.ErrorLog
)

type WebLog struct {
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

func InitLog(cfg *Config) {
	if cfg.LogFile != "" {
		logDir := filepath.Dir(cfg.LogFile) // 获取目录路径
		if err := os.MkdirAll(logDir, 0755); err != nil {
			ErrorLog.Fatal("创建日志目录失败: ", err)
		}
		file, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			WarningLog.Println("打开日志文件成功")
			Logger = &WebLog{
				Info:    log.New(file, "[INFO] ", log.LstdFlags),
				Warning: log.New(file, "[WARNING] ", log.LstdFlags),
				Error:   log.New(file, "[ERROR] ", log.LstdFlags),
			}
			return
		}
		WarningLog.Println("打开日志文件失败, 使用标准输出: ", err)
	} else {
		WarningLog.Println("未配置日志文件，使用标准输出")
	}
	Logger = &WebLog{
		Info:    InfoLog,
		Warning: WarningLog,
		Error:   ErrorLog,
	}
}
