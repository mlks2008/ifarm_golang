package zaplogger

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

type (
	Config struct {
		Level           string `json:"level"`             // 输出日志等级
		StackLevel      string `json:"stack_level"`       // 堆栈输出日志等级
		EnableConsole   bool   `json:"enable_console"`    // 是否控制台输出
		EnableWriteFile bool   `json:"enable_write_file"` // 是否输出文件(必需配置FilePath)
		MaxAge          int    `json:"max_age"`           // 最大保留天数(达到限制，则会被清理)
		TimeFormat      string `json:"time_format"`       // 打印时间输出格式
		PrintCaller     bool   `json:"print_caller"`      // 是否打印调用函数
		RotationTime    int    `json:"rotation_time"`     // 日期分割时间(秒)
		FileLinkPath    string `json:"file_link_path"`    // 日志文件连接路径
		FilePathFormat  string `json:"file_path_format"`  // 日志文件路径格式
		IncludeStdout   bool   `json:"include_stdout"`    // 是否包含os.stdout输出
		IncludeStderr   bool   `json:"include_stderr"`    // 是否包含os.stderr输出
	}
)

func DefaultConfig(path, serviceName, refLoggerName string) *Config {
	config := &Config{
		Level:           "info",
		StackLevel:      "error",
		EnableConsole:   false,
		EnableWriteFile: true,
		MaxAge:          7,
		TimeFormat:      "15:04:05.000", //2006-01-02 15:04:05.000
		PrintCaller:     false,
		RotationTime:    86400,
		FileLinkPath:    GetFileLinkPath(path, serviceName, refLoggerName),
		FilePathFormat:  GetFilePathFormat(path, serviceName, refLoggerName),
		IncludeStdout:   false,
		IncludeStderr:   false,
	}
	return config
}

func DefaultConsoleConfig() *Config {
	config := &Config{
		Level:           "debug",
		StackLevel:      "error",
		EnableConsole:   true,
		EnableWriteFile: false,
		MaxAge:          7,
		TimeFormat:      "15:04:05.000", //2006-01-02 15:04:05.000
		PrintCaller:     true,
		RotationTime:    86400,
		FileLinkPath:    "./logs/logstash/logstash.log",
		FilePathFormat:  "./logs/logstash/%Y%m%d_logstash.log",
		IncludeStdout:   false,
		IncludeStderr:   false,
	}
	return config
}

func (c *Config) TimeEncoder() zapcore.TimeEncoder {
	return func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format(c.TimeFormat))
	}
}

func GetLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

func getFilePath(path, serviceName, loggerName string) string {
	if path == "" {
		path = "."
	}

	return fmt.Sprintf("%s/logs/%s/%s", path, serviceName, loggerName)
}

func GetFileLinkPath(path, serviceName, loggerName string) string {
	return fmt.Sprintf("%s/%s.log", getFilePath(path, serviceName, loggerName), loggerName)
}

func GetFilePathFormat(path, serviceName, loggerName string) string {
	return fmt.Sprintf("%s/%s_%s.log", getFilePath(path, serviceName, loggerName), "%Y%m%d", loggerName)
}
