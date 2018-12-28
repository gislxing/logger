package logger

// 常量-日志级别
const (
	LogLevelDebug = iota
	LogLevelTrace
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

func getLogLevel(level int) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
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
