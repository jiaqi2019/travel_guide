package logger

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	// 创建日志文件
	logPath := "logs"
	if err := os.MkdirAll(logPath, 0755); err != nil {
		log.Fatal("Failed to create log directory:", err)
	}

	currentTime := time.Now().Format("2006-01-02")
	logFile, err := os.OpenFile(
		"logs/"+currentTime+".log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// 创建多重写入器，同时写入文件和标准输出
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// 设置日志格式
	InfoLogger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
