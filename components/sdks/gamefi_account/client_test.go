package gamefi_account

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

//func getClient() *Client {
//	fmt.Println(filepath.Abs(fmt.Sprintf("../../../../configs/dev")))
//
//	c := config.New(
//		config.WithSource(
//			file.NewSource(fmt.Sprintf("../../../../configs/dev")),
//		),
//	)
//	defer c.Close()
//
//	if err := c.Load(); err != nil {
//		panic(err)
//	}
//
//	var bc conf.Bootstrap
//	if err := c.Scan(&bc); err != nil {
//		panic(err)
//	}
//
//	return NewClient(bc.Client, initLogger(bc.Log))
//}

func getClient() *Client {

	return NewClient("http://10.40.10.10:30006", 5*time.Second, "sdks", initLogger())
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

func Test_SendEmail(t *testing.T) {
	client := getClient()

	res, err := client.SendEmail(context.Background(), "test@163.com", "subject", "content")
	fmt.Println("Res:", res, err)
	fmt.Println("")
}
