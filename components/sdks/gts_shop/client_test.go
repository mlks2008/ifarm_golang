package gts_shop

import (
	clog "components/log/zaplogger"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func getClient() *Client {

	return NewClient("http://10.40.10.10:30011", 5*time.Second, "sdks", initLogger())
}

func initLogger() log.Logger {
	caller := func(depth int) log.Valuer {
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
	return log.With(clog.NewLoggerWithName("", "sdks", "", true),
		"ts", log.DefaultTimestamp,
		"caller", caller(4),
		"service.id", "serviceid",
		"service.name", "sdks",
		"service.version", "serviceversion",
	)
}

func Test_AuthCode(t *testing.T) {
	client := getClient()

	res, err := client.AuthCode(context.Background(), "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1aWQiOiIxNDk4MDUyMTkxOTkzOTI1IiwiaXNHdWVzdCI6ZmFsc2UsImV4cCI6MTcwMzgyMTQ0MH0.n6vdGGVkbxfjYG58cPI-WDG69bEAy2_9HIjoNttebrQ")
	fmt.Println("Res:", res, err)
	fmt.Println("")
}

func Test_VerifyCode(t *testing.T) {
	client := getClient()

	res, err := client.VerifyCode(context.Background(), "1498052191993925", "284143")
	fmt.Println("Res:", res, err)
	fmt.Println("")
}
