package logger

var (
	log      logInterface
	logLevel int
	logModel int
)

func init() {
	createLog()
}

// 创建日志对象
func createLog() {
	switch logModel {
	case ConsoleModel:
		log, _ = newConsoleLogger(logLevel)
	case FileModel:
		log, _ = newFileLogger(logLevel)
	}
}

// 设置日志输出级别，如果不设置则默认为 info 级别
func SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		log.setLevel(LogLevelDebug)
	}

	log.setLevel(level)
}

// 设置日志输出模式（文件模式、控制台模式）
func SetLogModel(model int) {
	if model < ConsoleModel || model > FileModel {
		logModel = ConsoleModel
	}

	logModel = model
	createLog()
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
