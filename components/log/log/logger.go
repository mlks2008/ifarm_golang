package log

import (
	clog "components/log/zaplogger"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	id, _ = os.Hostname()
	once  sync.Once
)

func caller(depth int) log.Valuer {
	return func(context.Context) interface{} {
		_, file, line, _ := runtime.Caller(depth)
		idx := strings.LastIndexByte(file, '/')
		if idx == -1 {
			return file[idx+1:] + ":" + strconv.Itoa(line)
		}
		for i := 0; i < 2; i++ {
			idx = strings.LastIndexByte(file[:idx], '/')
		}
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

var Logger *log.Helper

func init() {
	var logDir = ""
	var serviceName = "default"
	Logger = log.NewHelper(log.With(clog.NewLoggerWithName(logDir, serviceName, "", true),
		"ts", log.DefaultTimestamp,
		"caller", caller(4),
		"service.id", id,
		"service.name", serviceName,
		"service.version", "1.0.0",
	))
}

func InitLogger(logDir, serviceName string, debug bool) *log.Helper {
	once.Do(func() {
		Logger = log.NewHelper(log.With(clog.NewLoggerWithName(logDir, serviceName, "", debug),
			"ts", log.DefaultTimestamp,
			"caller", caller(4),
			//"service.id", id,
			//"service.name", serviceName,
			//"service.version", "1.0.0",
		))
	})
	return Logger
}
