package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type fileLogger struct {
	baseLogger
	infoFile *os.File          // info及以下级别写入的文件
	errFile  *os.File          // error及以上级别写入的文件
	dataChan chan *fileLogData // 写入文件通道
}

type fileLogData struct {
	log  string   // 写入日志文件的内容
	file *os.File // 写入的文件
}

func newFileLogger(level int) (logInterface, error) {
	log := &fileLogger{
		dataChan: make(chan *fileLogData, chanCacheSize),
		baseLogger: baseLogger{
			logFileMaxSize: splitFileSize,
		},
	}

	log.setLevel(level)
	log.createLogFile()

	// 启动写入协程
	go log.writerToFile()

	return log, nil
}

func (f *fileLogger) setLogTotalSize(size int64) {
	f.logTotalSize = size
}

func (f *fileLogger) getLogParam() (logLevel int, logFileMaxSize int64, logTotalSize int64) {
	logLevel, logFileMaxSize, logTotalSize = f.level, f.logFileMaxSize, f.logTotalSize
	return
}

func (f *fileLogger) setLogFileMaxSize(size int64) {
	if size <= 0 {
		size = splitFileSize
	}

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
	logPath := getCurrentLogPath(false)

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
		fmt.Printf("open log file failed: %v", err)
		return
	}

	f.infoFile = file

	// 打开error级别日志文件(写入error、fatal日志)
	logFilePath = filepath.Join(logPath, projectName)
	logFilePath = fmt.Sprintf("%s%s", logFilePath, "-error.log")

	// 日志 error 文件不存在则创建
	file, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("open log file failed: %v", err)
		return
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

	if currentLogPath != getCurrentLogPath(false) {
		f.createLogFile()
		return
	}

	// 按照文件大小切分日志
	currentFileSize := fileInfo.Size()
	if f.logFileMaxSize <= currentFileSize {
		timeStr := time.Now().Format("20060102150405")
		newPath := strings.Replace(oldPath, ".log", "-"+timeStr+".log", -1)

		for i := 0; i < 5; i++ {
			data.file.Close()
			err := os.Rename(oldPath, newPath)
			if err != nil {
				// 重命名文件失败，则重试
				continue
			}

			break
		}

		// 重新创建日志文件
		f.createLogFile()

		// 清理日志
		if f.logTotalSize > 0 {
			go f.clearLogFile()
		}
	}

}

// 清理日志
// 日志总大小大于设定值后开始清理旧日志（清理30%总量最大值）
func (f *fileLogger) clearLogFile() {
	// 获取日志总大小
	rootPath := getCurrentLogPath(true)
	totalSize := getPathSize(rootPath)

	// 当前日志小于临界值不清理
	if totalSize < f.logTotalSize {
		return
	}

	// 从最早的日志开始清理
	delSizeFloat := float64(f.logTotalSize) * perDeleteLog
	if delSize, err := strconv.Atoi(fmt.Sprintf("%.0f", delSizeFloat)); err == nil {
		delLog(rootPath, int64(delSize), f.logFileMaxSize)
	}
}
