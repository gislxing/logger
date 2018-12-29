package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type fileLogger struct {
	baseLogger
	//logPath        string            // 日志文件路径
	//logName        string            // 日志问价名称
	infoFile       *os.File          // info及以下级别写入的文件
	errFile        *os.File          // error及以上级别写入的文件
	dataChan       chan *fileLogData // 写入文件通道
	logFileMaxSize int64             // 默认日志文件最大值（单位: 字节）
}

type fileLogData struct {
	log  string   // 写入日志文件的内容
	file *os.File // 写入的文件
}

func newFileLogger(level int) (logInterface, error) {
	log := &fileLogger{
		dataChan:       make(chan *fileLogData, chanCacheSize),
		logFileMaxSize: splitFileSize,
	}

	log.setLevel(level)
	log.createLogFile()

	// 启动写入协程
	go log.writerToFile()

	return log, nil
}

func (f *fileLogger) setLogFileMaxSize(size int64) {
	f.logFileMaxSize = size
}

// 写入文件
func (f *fileLogger) writerToFile() {
	for data := range f.dataChan {
		// 切分日志
		f.splitLogFile(data)

		fmt.Fprintf(data.file, data.log)
	}
}

// 创建日志文件
func (f *fileLogger) createLogFile() {
	projectName := getProjectName()
	logPath := getCurrentLogPath()

	// 路径不存在则创建
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.MkdirAll(logPath, os.ModePerm)
	}

	// 打开info级别日志文件（写入debug、trace、info、warn日志）
	logFilePath := filepath.Join(logPath, projectName)
	logFilePath = fmt.Sprintf("%s%s", logFilePath, ".log")

	// 日志 info 文件不存在则创建
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("open log file failed: %v", err))
	}

	f.infoFile = file

	// 打开error级别日志文件(写入error、fatal日志)
	logFilePath = filepath.Join(logPath, projectName)
	logFilePath = fmt.Sprintf("%s%s", logFilePath, "-error.log")

	// 日志 error 文件不存在则创建
	file, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("open log file failed: %v", err))
	}

	f.errFile = file
}

func (f *fileLogger) setLevel(level int) {
	if level < DEBUG || level > FATAL {
		f.level = INFO
	}

	f.level = level
}

func (f *fileLogger) close() {
	f.infoFile.Close()
	f.errFile.Close()
}

func (f *fileLogger) debug(format string, args ...interface{}) {
	if f.level > DEBUG {
		return
	}

	output(f.infoFile, f.dataChan, DEBUG, format, args...)
}

func (f *fileLogger) trace(format string, args ...interface{}) {
	if f.level > TRACE {
		return
	}

	output(f.infoFile, f.dataChan, TRACE, format, args...)
}

func (f *fileLogger) info(format string, args ...interface{}) {
	if f.level > INFO {
		return
	}

	output(f.infoFile, f.dataChan, INFO, format, args...)
}

func (f *fileLogger) warn(format string, args ...interface{}) {
	if f.level > WARN {
		return
	}

	output(f.infoFile, f.dataChan, WARN, format, args...)
}

func (f *fileLogger) error(format string, args ...interface{}) {
	if f.level > ERROR {
		return
	}

	output(f.errFile, f.dataChan, ERROR, format, args...)
}

func (f *fileLogger) fatal(format string, args ...interface{}) {
	if f.level > FATAL {
		return
	}

	output(f.errFile, f.dataChan, FATAL, format, args...)
}

func (f *fileLogger) splitLogFile(data *fileLogData) {
	// 检查文件大小
	fileInfo, err := data.file.Stat()
	if err != nil {
		f.createLogFile()
		return
	}

	// 检查是否需要重新创建日志目录和文件
	oldPath := data.file.Name()
	currentLogPath := filepath.Dir(oldPath)

	if currentLogPath != getCurrentLogPath() {
		f.createLogFile()
		return
	}

	// 按照文件大小切分日志
	currentFileSize := fileInfo.Size()
	if f.logFileMaxSize <= currentFileSize {
		timeStr := time.Now().Format("20060102150405")
		newPath := strings.Replace(oldPath, ".log", "-"+timeStr+".log", -1)
		f.close()
		os.Rename(oldPath, newPath)

		// 重新创建日志文件
		f.createLogFile()
	}

}
