package logger

import "os"

type baseLogger struct {
	level          int   // 日志级别
	logFileMaxSize int64 // 默认日志文件最大值（单位: 字节）
	logTotalSize   int64 // 日志总大小，大于该值则清理 30%
}

type consoleLogger struct {
	baseLogger
}

func newConsoleLogger(level int) (logInterface, error) {
	log := &consoleLogger{
		baseLogger: baseLogger{
			logFileMaxSize: splitFileSize,
		},
	}

	log.setLevel(level)
	return log, nil
}

func (c *consoleLogger) setLogTotalSize(size int64) {
	c.logTotalSize = size
}

func (c *consoleLogger) getLogParam() (logLevel int, logFileMaxSize int64, logTotalSize int64) {
	logLevel, logFileMaxSize, logTotalSize = c.level, c.logFileMaxSize, c.logTotalSize
	return
}

func (c *consoleLogger) setLogFileMaxSize(size int64) {
	if size <= 0 {
		size = splitFileSize
	}

	c.logFileMaxSize = size
}

func (c *consoleLogger) setLevel(level int) {
	if level < DEBUG || level > FATAL {
		c.level = INFO
	}

	c.level = level
}

func (c *consoleLogger) debug(format string, args ...interface{}) {
	if c.level > DEBUG {
		return
	}

	output(os.Stdout, nil, DEBUG, format, args...)
}

func (c *consoleLogger) trace(format string, args ...interface{}) {
	if c.level > TRACE {
		return
	}

	output(os.Stdout, nil, TRACE, format, args...)
}

func (c *consoleLogger) info(format string, args ...interface{}) {
	if c.level > INFO {
		return
	}

	output(os.Stdout, nil, INFO, format, args...)
}

func (c *consoleLogger) warn(format string, args ...interface{}) {
	if c.level > WARN {
		return
	}

	output(os.Stdout, nil, WARN, format, args...)
}

func (c *consoleLogger) error(format string, args ...interface{}) {
	if c.level > ERROR {
		return
	}

	output(os.Stdout, nil, ERROR, format, args...)
}

func (c *consoleLogger) fatal(format string, args ...interface{}) {
	if c.level > FATAL {
		return
	}

	output(os.Stdout, nil, FATAL, format, args...)
}

func (c *consoleLogger) close() {

}
