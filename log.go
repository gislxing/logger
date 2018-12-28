package logger

import (
	"reflect"
)

var (
	log logInterface
)

func init() {
	log, _ = newConsoleLogger(LogLevelDebug)
}

// 设置日志输出级别，如果不设置则默认为 info 级别
func SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		log.setLevel(LogLevelDebug)
	}

	log.setLevel(level)
}

// 设置日志输出模式（文件模式 FileModel、控制台模式 ConsoleModel）
func SetLogModel(model int) {
	switch model {
	case ConsoleModel:
		if reflect.TypeOf(log) != reflect.TypeOf(&consoleLogger{}) {
			log, _ = newConsoleLogger(LogLevelDebug)
		}
	case FileModel:
		if reflect.TypeOf(log) != reflect.TypeOf(&fileLogger{}) {
			log, _ = newFileLogger(LogLevelDebug)
		}
	default:
		if reflect.TypeOf(log) != reflect.TypeOf(&consoleLogger{}) {
			log, _ = newConsoleLogger(LogLevelDebug)
		}
	}
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
