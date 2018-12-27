package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileLogger struct {
	level    int
	logPath  string
	logName  string
	infoFile *os.File
	errFile  *os.File
}

func NewFileLogger(level int, logPath, logName string) (LogInterface, error) {
	log := &FileLogger{
		logPath: logPath,
		logName: logName,
	}

	log.SetLevel(level)
	log.init()
	return log, nil
}

func (f *FileLogger) init() {
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

func (f *FileLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		f.level = LogLevelInfo
	}

	f.level = level
}

func (f *FileLogger) Close() {
	f.infoFile.Close()
	f.errFile.Close()
}

func (f *FileLogger) Debug(format string, args ...interface{}) {
	if f.level > LogLevelDebug {
		return
	}

	writeLog(f.infoFile, LogLevelDebug, format, args...)
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	if f.level > LogLevelTrace {
		return
	}

	writeLog(f.infoFile, LogLevelTrace, format, args...)
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}

	writeLog(f.infoFile, LogLevelInfo, format, args...)
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	if f.level > LogLevelWarn {
		return
	}

	writeLog(f.infoFile, LogLevelWarn, format, args...)
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	if f.level > LogLevelError {
		return
	}

	writeLog(f.errFile, LogLevelError, format, args...)
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	if f.level > LogLevelFatal {
		return
	}

	writeLog(f.errFile, LogLevelFatal, format, args...)
}
