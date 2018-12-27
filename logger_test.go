package logger

import (
	"fmt"
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
		err := os.RemoveAll(logFilePath)
		if err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		t.Errorf("创建日志对象失败: %v", err)
	}

	log.Debug("test debug level log")
	log.Trace("test trace level log")
	log.Info("test info level log")
	log.Warn("test warn level log")
	log.Error("test error level log")
	log.Fatal("test fatal level log")

	filePath := filepath.Join(logFilePath, logFileName) + ".log"
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("创建日志文件失败: %v", err)
	}

	if fileInfo.Size() <= 0 {
		t.Errorf("写入日志失败")
	}
}
