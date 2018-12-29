package logger

import (
	"reflect"
)

var (
	log logInterface
)

func init() {
	log, _ = newConsoleLogger(DEBUG)
}

// 设置日志输出级别，如果不设置则默认为 info 级别
func SetLogLevel(level int) {
	if level < DEBUG || level > FATAL {
		log.setLevel(DEBUG)
	}

	log.setLevel(level)
}

// 设置日志输出模式（文件模式 FileModel、控制台模式 ConsoleModel）
func SetLogModel(model int) {
	switch model {
	case ConsoleModel:
		if reflect.TypeOf(log) != reflect.TypeOf(&consoleLogger{}) {
			logLevel, logFileMaxSize := log.getLogParam()
			log, _ = newConsoleLogger(DEBUG)
			log.setLevel(logLevel)
			log.setLogFileMaxSize(logFileMaxSize)
		}
	case FileModel:
		if reflect.TypeOf(log) != reflect.TypeOf(&fileLogger{}) {
			logLevel, logFileMaxSize := log.getLogParam()
			log, _ = newFileLogger(DEBUG)
			log.setLevel(logLevel)
			log.setLogFileMaxSize(logFileMaxSize)
		}
	default:
		if reflect.TypeOf(log) != reflect.TypeOf(&consoleLogger{}) {
			logLevel, logFileMaxSize := log.getLogParam()
			log, _ = newConsoleLogger(DEBUG)
			log.setLevel(logLevel)
			log.setLogFileMaxSize(logFileMaxSize)
		}
	}
}

// 设置日志文件最大大小（字节）,如果超过这个值则切分日志
// 如果不设置，默认日志切分大小是 100MB
func SetLogFileMaxSize(size int64) {
	log.setLogFileMaxSize(size)
}

func Debug(format string, args ...interface{}) {
	log.debug(format, args...)
}

func Trace(format string, args ...interface{}) {
	log.trace(format, args...)
}

func Info(format string, args ...interface{}) {
	log.info(format, args...)
}

func Warn(format string, args ...interface{}) {
	log.warn(format, args...)
}

func Error(format string, args ...interface{}) {
	log.error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	log.fatal(format, args...)
}
