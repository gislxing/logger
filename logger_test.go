package logger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileLogger(t *testing.T) {
	projectName := getProjectName()
	logPath := filepath.Join("/var/log/", projectName)

	log, err := newFileLogger(DEBUG)
	defer func() {
		log.close()
		os.RemoveAll(logPath)
	}()
	if err != nil {
		t.Errorf("创建日志对象失败: %v", err)
	}

	log.debug("test file debug level log")
	log.trace("test file trace level log")
	log.info("test file info level log")
	log.warn("test file warn level log")
	log.error("test file error level log")
	log.fatal("test file fatal level log")

	filePath := filepath.Join(logPath, projectName) + ".log"
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("创建日志文件失败: %v", err)
	}

	if fileInfo.Size() <= 0 {
		t.Errorf("写入日志失败")
	}
}

func TestConsoleLogger(t *testing.T) {
	log, err := newConsoleLogger(DEBUG)
	if err != nil {
		t.Errorf("创建控制台日志错误: %v", err)
	}

	log.debug("test debug level log")
	log.trace("test trace level log")
	log.info("test info level log")
	log.warn("test warn level log")
	log.error("test error level log")
	log.fatal("test fatal level log")
}
