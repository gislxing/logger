package logger

// 日志接口
type logInterface interface {
	// 设置日志级别，如果设置错误或者不设置则默认使用 info 级别
	setLevel(level int)

	// 设置日志文件最大大小值，超过该值则切分
	// 默认大小 100MB
	setLogFileMaxSize(size int64)

	// 获得日志文件的配置信息
	getLogParam() (logLevel int, logFileMaxSize int64)

	debug(format string, args ...interface{})
	trace(format string, args ...interface{})
	info(format string, args ...interface{})
	warn(format string, args ...interface{})
	error(format string, args ...interface{})
	fatal(format string, args ...interface{})
	close()
}
