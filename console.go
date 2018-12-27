package logger

import "os"

type BaseLogger struct {
	level int
}

type ConsoleLogger struct {
	BaseLogger
}

func NewConsoleLogger(level int) (LogInterface, error) {
	log := &ConsoleLogger{}

	log.SetLevel(level)
	log.init()
	return log, nil
}

func (c *ConsoleLogger) init() {

}

func (c *ConsoleLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		c.level = LogLevelInfo
	}

	c.level = level
}

func (c *ConsoleLogger) Debug(format string, args ...interface{}) {
	if c.level > LogLevelDebug {
		return
	}

	output(os.Stdout, LogLevelDebug, format, args...)
}

func (c *ConsoleLogger) Trace(format string, args ...interface{}) {
	if c.level > LogLevelTrace {
		return
	}

	output(os.Stdout, LogLevelTrace, format, args...)
}

func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	if c.level > LogLevelInfo {
		return
	}

	output(os.Stdout, LogLevelInfo, format, args...)
}

func (c *ConsoleLogger) Warn(format string, args ...interface{}) {
	if c.level > LogLevelWarn {
		return
	}

	output(os.Stdout, LogLevelWarn, format, args...)
}

func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	if c.level > LogLevelError {
		return
	}

	output(os.Stdout, LogLevelError, format, args...)
}

func (c *ConsoleLogger) Fatal(format string, args ...interface{}) {
	if c.level > LogLevelFatal {
		return
	}

	output(os.Stdout, LogLevelFatal, format, args...)
}

func (c *ConsoleLogger) Close() {

}
