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

func (c *consoleLogger) setLogFileMaxSize(size int64) {

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
