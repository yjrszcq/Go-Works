package web_log

import "log"

var WebLogger *WebLog

type WebLog struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}
