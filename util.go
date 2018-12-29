package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex

// 获取程序运行过程中日志调用的文件名称、函数名称你、行号
func getFuncCallInfo() (fileName, funcName string, lineNo int) {
	mu.Lock()
	defer mu.Unlock()

	pc, file, line, ok := runtime.Caller(4)
	if ok {
		fileName = file
		funcName = runtime.FuncForPC(pc).Name()
		lineNo = line
	}

	return
}

// 输出日志输出
// 输出格式: 时间 [日志级别] 文件名称:方法名称 行号 - 用户日志
// 如果 dataChan = nil 则直接输出到 file, 否则将日志发送到 dataChan
func output(file *os.File, dataChan chan *fileLogData, level int, format string, args ...interface{}) {
	timeStr := time.Now().Format("2006-01-02 15:04:05.999")
	levelStr := getLogLevel(level)
	fileName, funcName, lineNo := getFuncCallInfo()
	fileName = filepath.Base(fileName)
	funcName = filepath.Base(funcName)

	msg := fmt.Sprintf(format, args...)

	logMsg := fmt.Sprintf("%s [%5s] %s:%s %d - %s\n", timeStr, levelStr, fileName, funcName, lineNo, msg)
	logData := &fileLogData{
		log:  logMsg,
		file: file,
	}

	if dataChan != nil {
		select {
		case dataChan <- logData:
		default:
		}
	} else {
		fmt.Fprintf(file, logMsg)
	}

}

// 获取项目名称
func getProjectName() string {
	name := os.Args[0]
	name = filepath.Base(name)
	index := strings.LastIndex(name, ".")
	return name[:index]
}

// 获得当前日志的保存路径,
// isLogRootDir 是否要获取完整日志目录,
// true：返回日志根目录，false：返回完整日志路径
func getCurrentLogPath(isLogRootDir bool) string {
	now := time.Now()
	projectName := getProjectName()

	if isLogRootDir {
		return string(filepath.Separator) + filepath.Join("var", "log", projectName)
	}

	return string(filepath.Separator) +
		filepath.Join("var", "log", projectName, strconv.Itoa(now.Year()), fmt.Sprintf("%d", now.Month()))
}

// 检查文件或者路径是否存在
// true：存在，false：不存在
func isExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// 获取路径或者文件大小（单位：字节）
func getPathSize(path string) (size int64) {
	if path == "" || !isExists(path) {
		return
	}

	infos, err := ioutil.ReadDir(path)
	if err != nil {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return
		}

		return fileInfo.Size()
	}

	for _, value := range infos {
		if value.IsDir() {
			size += getPathSize(filepath.Join(path, value.Name()))
		} else {
			size += value.Size()
		}
	}

	return
}
