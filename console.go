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
	return log, nil
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

	output(os.Stdout, nil, LogLevelDebug, format, args...)
}

func (c *consoleLogger) trace(format string, args ...interface{}) {
	if c.level > LogLevelTrace {
		return
	}

	output(os.Stdout, nil, LogLevelTrace, format, args...)
}

func (c *consoleLogger) info(format string, args ...interface{}) {
	if c.level > LogLevelInfo {
		return
	}

	output(os.Stdout, nil, LogLevelInfo, format, args...)
}

func (c *consoleLogger) warn(format string, args ...interface{}) {
	if c.level > LogLevelWarn {
		return
	}

	output(os.Stdout, nil, LogLevelWarn, format, args...)
}

func (c *consoleLogger) error(format string, args ...interface{}) {
	if c.level > LogLevelError {
		return
	}

	output(os.Stdout, nil, LogLevelError, format, args...)
}

func (c *consoleLogger) fatal(format string, args ...interface{}) {
	if c.level > LogLevelFatal {
		return
	}

	output(os.Stdout, nil, LogLevelFatal, format, args...)
}

func (c *consoleLogger) close() {

}
