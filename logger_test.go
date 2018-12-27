package logger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileLogger(t *testing.T) {
	logFilePath := "/var/log/logger"
	logFileName := "logger"
	log, err := NewFileLogger(LogLevelDebug, logFilePath, logFileName)
	defer func() {
		log.Close()
		os.RemoveAll(logFilePath)
	}()
	if err != nil {
		t.Errorf("创建日志对象失败: %v", err)
	}

	log.Debug("test file debug level log")
	log.Trace("test file trace level log")
	log.Info("test file info level log")
	log.Warn("test file warn level log")
	log.Error("test file error level log")
	log.Fatal("test file fatal level log")

	filePath := filepath.Join(logFilePath, logFileName) + ".log"
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("创建日志文件失败: %v", err)
	}

	if fileInfo.Size() <= 0 {
		t.Errorf("写入日志失败")
	}
}

func TestConsoleLogger(t *testing.T) {
	log, err := NewConsoleLogger(LogLevelDebug)
	if err != nil {
		t.Errorf("创建控制台日志错误: %v", err)
	}

	log.Debug("test debug level log")
	log.Trace("test trace level log")
	log.Info("test info level log")
	log.Warn("test warn level log")
	log.Error("test error level log")
	log.Fatal("test fatal level log")
}
