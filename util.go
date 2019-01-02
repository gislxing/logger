package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
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

// 按照年份从小到大，月份从小到大的顺序删除日志，
// 如果不是日志模块生成的文件或者目录则全部删除
func delLog(path string, delSize int64, maxFileSize int64) {
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	// 如果是空目录则直接删除
	if len(infos) == 0 {
		os.Remove(path)
		return
	}

	// 遍历年的目录
	years := make([]int, 0)
	for _, value := range infos {
		if year, err := strconv.Atoi(value.Name()); err == nil {
			years = append(years, year)
		} else {
			// 不是日志目录则全部删除
			os.RemoveAll(filepath.Join(path, value.Name()))
		}
	}

	// 按照从小到大排序
	sort.Ints(years)

	// 遍历每个年份的月份
	for _, value := range years {
		yearStr := strconv.Itoa(value)
		isDel := delLogMonth(path, yearStr, &delSize, maxFileSize)
		if !isDel {
			return
		}

		// 所有文件都删除了，则删除目录
		os.RemoveAll(filepath.Join(path, yearStr))
	}
}

// 按照月份从小到大删除
// 返回是否需要接着删除，true：是，false：否
func delLogMonth(rootPath, yearStr string, delSize *int64, maxFileSize int64) bool {
	yearPath := filepath.Join(rootPath, yearStr)
	fileInfos, err := ioutil.ReadDir(yearPath)
	if err != nil {
		return true
	}

	months := make([]int, 0)
	for _, monthStr := range fileInfos {
		if month, err := strconv.Atoi(monthStr.Name()); err == nil {
			months = append(months, month)
		} else {
			// 不是日志模块目录则删除
			os.RemoveAll(filepath.Join(yearPath, monthStr.Name()))
		}
	}

	sort.Ints(months)

	for _, month := range months {
		monthStr := strconv.Itoa(month)
		isDel := delLogFile(yearPath, monthStr, delSize, maxFileSize)
		if !isDel {
			return false
		}

		// 所有文件都删除了，则删除目录
		os.RemoveAll(filepath.Join(yearPath, monthStr))
	}

	return true
}

// 按照日志生成的时间删除
// 返回是否需要接着删除，true：是，false：否
func delLogFile(yearPath, month string, delSize *int64, maxFileSize int64) bool {
	infos, err := ioutil.ReadDir(filepath.Join(yearPath, month))
	if err != nil {
		return true
	}

	// 提取日志文件名中的时间字符串，按照从小到大的顺序删除
	// 如果没有则最后删除
	logTimeMap := make(map[int]string)
	lastRemoveFile := make([]string, 0)
	timeRep := regexp.MustCompile(`\d{14,}`)
	for _, value := range infos {
		s := timeRep.FindString(value.Name())
		path := filepath.Join(yearPath, month, value.Name())
		if tmp, err := strconv.Atoi(s); err == nil {
			logTimeMap[tmp] = path
		} else {
			lastRemoveFile = append(lastRemoveFile, path)
		}
	}

	// 对key进行排序
	keys := make([]int, 0)
	for key := range logTimeMap {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	// 按照时间大小删除日志
	delCount := 0
	for _, key := range keys {
		path := logTimeMap[key]

		// 获取当前文件的大小
		if info, err := os.Stat(path); err == nil {
			if (info.Size() <= *delSize || info.Size() >= maxFileSize) && *delSize > 0 {
				// 当前文件大小小于等于总共删除的文件大小则删除
				os.Remove(path)
				*delSize -= info.Size()
				delCount++
			}
		}
	}

	if delCount != len(keys) {
		return false
	}

	// 删除其余文件
	for _, value := range lastRemoveFile {
		if info, err := os.Stat(value); err == nil {
			if (info.Size() <= *delSize || info.Size() >= maxFileSize) && *delSize > 0 {
				// 当前文件大小小于等于总共删除的文件大小则删除
				os.Remove(value)
				*delSize -= info.Size()
			}
		}
	}

	return true
}
