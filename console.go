package logger

import "os"

type baseLogger struct {
	level int
}

type consoleLogger struct {
	baseLogger
}

func newConsoleLogger(level int) (logInterface, error) {
	log := &consoleLogger{}

	log.setLevel(level)
	log.init()
	return log, nil
}

func (c *consoleLogger) init() {

}

func (c *consoleLogger) setLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		c.level = LogLevelInfo
	}

	c.level = level
}

func (c *consoleLogger) debug(format string, args ...interface{}) {
	if c.level > LogLevelDebug {
		return
	}

	output(os.Stdout, LogLevelDebug, format, args...)
}

func (c *consoleLogger) trace(format string, args ...interface{}) {
	if c.level > LogLevelTrace {
		return
	}

	output(os.Stdout, LogLevelTrace, format, args...)
}

func (c *consoleLogger) info(format string, args ...interface{}) {
	if c.level > LogLevelInfo {
		return
	}

	output(os.Stdout, LogLevelInfo, format, args...)
}

func (c *consoleLogger) warn(format string, args ...interface{}) {
	if c.level > LogLevelWarn {
		return
	}

	output(os.Stdout, LogLevelWarn, format, args...)
}

func (c *consoleLogger) error(format string, args ...interface{}) {
	if c.level > LogLevelError {
		return
	}

	output(os.Stdout, LogLevelError, format, args...)
}

func (c *consoleLogger) fatal(format string, args ...interface{}) {
	if c.level > LogLevelFatal {
		return
	}

	output(os.Stdout, LogLevelFatal, format, args...)
}

func (c *consoleLogger) close() {

}
