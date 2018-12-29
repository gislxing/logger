package logger

// 常量-日志级别
const (
	DEBUG = iota
	TRACE
	INFO
	WARN
	ERROR
	FATAL
)

func getLogLevel(level int) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// 日志输出模式
const (
	ConsoleModel = iota
	FileModel
)

const (
	// 文件写入数据通道缓存大小
	chanCacheSize = 50000

	// 1MB
	MB int64 = 1048576

	// 默认日志切分大小
	splitFileSize = 100 * MB
)
