package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// 获取程序运行过程中日志调用的文件名称、函数名称你、行号
func getFuncCallInfo() (fileName, funcName string, lineNo int) {
	pc, file, line, ok := runtime.Caller(3)
	if ok {
		fileName = file
		funcName = runtime.FuncForPC(pc).Name()
		lineNo = line
	}

	return
}

// 将日志输出到文件
// 输出格式: 时间 [日志级别] 文件名称:方法名称 行号 - 用户日志
func output(file *os.File, level int, format string, args ...interface{}) {
	timeStr := time.Now().Format("2006-01-02 15:04:05.999")
	levelStr := getLogLevel(level)
	fileName, funcName, lineNo := getFuncCallInfo()
	fileName = filepath.Base(fileName)
	funcName = filepath.Base(funcName)

	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(file, "%s [%5s] %s:%s %d - %s\n", timeStr, levelStr, fileName, funcName, lineNo, msg)
}
