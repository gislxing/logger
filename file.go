package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

type fileLogger struct {
	baseLogger
	logPath  string
	logName  string
	infoFile *os.File
	errFile  *os.File
}

func newFileLogger(level int) (logInterface, error) {
	projectName := getProjectName()
	logPath := filepath.Join("/var/log/", projectName)
	log := &fileLogger{
		logPath: logPath,
		logName: projectName,
	}

	log.setLevel(level)
	log.init()
	return log, nil
}

func (f *fileLogger) init() {
	// 路径不存在则创建
	if _, err := os.Stat(f.logPath); os.IsNotExist(err) {
		os.MkdirAll(f.logPath, os.ModePerm)
	}

	// 打开info级别日志文件（写入debug、trace、info、warn日志）
	logFilePath := filepath.Join(f.logPath, f.logName)
	logFilePath = fmt.Sprintf("%s%s", logFilePath, ".log")
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("open log file failed: %v", err))
	}

	f.infoFile = file

	// 打开error级别日志文件(写入error、fatal日志)
	logFilePath = filepath.Join(f.logPath, f.logName)
	logFilePath = fmt.Sprintf("%s%s", logFilePath, "-error.log")
	file, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("open log file failed: %v", err))
	}

	f.errFile = file
}

func (f *fileLogger) setLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		f.level = LogLevelInfo
	}

	f.level = level
}

func (f *fileLogger) close() {
	f.infoFile.Close()
	f.errFile.Close()
}

func (f *fileLogger) debug(format string, args ...interface{}) {
	if f.level > LogLevelDebug {
		return
	}

	output(f.infoFile, LogLevelDebug, format, args...)
}

func (f *fileLogger) trace(format string, args ...interface{}) {
	if f.level > LogLevelTrace {
		return
	}

	output(f.infoFile, LogLevelTrace, format, args...)
}

func (f *fileLogger) info(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}

	output(f.infoFile, LogLevelInfo, format, args...)
}

func (f *fileLogger) warn(format string, args ...interface{}) {
	if f.level > LogLevelWarn {
		return
	}

	output(f.infoFile, LogLevelWarn, format, args...)
}

func (f *fileLogger) error(format string, args ...interface{}) {
	if f.level > LogLevelError {
		return
	}

	output(f.errFile, LogLevelError, format, args...)
}

func (f *fileLogger) fatal(format string, args ...interface{}) {
	if f.level > LogLevelFatal {
		return
	}

	output(f.errFile, LogLevelFatal, format, args...)
}
