# logger
Go语言日志

## 特性
> 开箱即用
>
> 支持日志切分
> 
> 支持日志自动清理
> 
> 支持日志压缩

## Quick Start
### Download and install
```
go get github.com/gislxing/logger
```

### 使用
```go
import "github.com/gislxing/logger"

// 默认只输出到控制台
// 默认日志级别: debug
// 日志输出路径: /var/log/项目名称/年/月/项目名称.log
logger.Debug("test log %s", "hello")
```

### 参数设定
```go
// 设置日志级别，不设置默认日志级别：debug
logger.SetLevel(logger.DEBUG)

// 设置输出日志到文件，不设置则默认只输出到控制台
logger.SetLogModel(logger.FileModel)

// 设置日志切分大小（MB），不设置默认按照 100MB 切分
logger.SetLogFileMaxSize(50 * logger.MB)

```

#### 开启日志自动清理
```go
// 此处设定日志最大总大小为 100MB，大于该值则自动清理最大值的 30%
// 不设置或者设置为小于等于的值，则关闭
// 自动清理功能默认关闭
logger.SetLogTotalSize(logger.MB * 100)
```

1、日志路径优化: /var/log/项目名称/年/月/项目名称.log（完成）

2、日志压缩：日志切分后进行

3、日志清理：日志总大小大于设定值后开始清理旧日志（清理30%总量最大值）（完成）

4、修正需要先设定日志模式后设置其他参数(完成)