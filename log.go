package logger

var (
	log      LogInterface
	logLevel int
	logModel int
)

func init() {
	switch logModel {
	case consoleModel:
		log, _ = NewConsoleLogger(logLevel)
	case fileModel:
		//log, _ = NewFileLogger(logLevel)
	}
}

// 设置日志输出级别，如果不设置则默认为 info 级别
func SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		logLevel = LogLevelInfo
	}

	logLevel = level
}

// 设置日志输出模式（文件模式、控制台模式）
func SetLogModel(model int) {
	if model < consoleModel || model > fileModel {
		logModel = consoleModel
	}

	logModel = model
}

func Debug(format string, args ...interface{}) {
	log.Debug(format, args...)
}
