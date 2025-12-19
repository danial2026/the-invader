package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	logFile *os.File
	logger  *log.Logger
)

func Init() error {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join("logs", fmt.Sprintf("game_%s.log", timestamp))

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	logFile = f
	logger = log.New(logFile, "", log.LstdFlags)

	Info("Logger initialized. Writing to %s", filename)
	return nil
}

func Info(format string, v ...interface{}) {
	if logger != nil {
		msg := fmt.Sprintf("INFO: "+format, v...)
		logger.Println(msg)
	}
}

func Error(format string, v ...interface{}) {
	if logger != nil {
		msg := fmt.Sprintf("ERROR: "+format, v...)
		logger.Println(msg)
	}
}

func Close() {
	if logFile != nil {
		Info("Closing logger.")
		logFile.Close()
	}
}

