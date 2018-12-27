package logger

// 日志接口
type LogInterface interface {
	// 初始化函数
	init()

	// 设置日志级别，如果设置错误或者不设置则默认使用 info 级别
	SetLevel(level int)
	Debug(format string, args ...interface{})
	Trace(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Close()
}
