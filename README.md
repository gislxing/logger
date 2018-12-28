# logger
Go语言日志


## Quick Start
### Download and install
```
go get github.com/gislxing/logger
```

### Usage
```go
import "github.com/gislxing/logger"

# 默认只输出到控制台
# 默认日志级别: debug
logger.Debug("test log %s", "hello")

# 设置日志级别
logger.SetLevel(logger.LogLevelInfo)

# 设置输出日志到文件
logger.SetLogModel(logger.FileModel)
```
