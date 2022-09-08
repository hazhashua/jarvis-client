package utils

import (
	"fmt"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	if Logger == nil {
		fmt.Println("日志对象为空，创建日志对象...")
		if logFile, err := os.OpenFile("./exporter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			Logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
		}
	}
}
