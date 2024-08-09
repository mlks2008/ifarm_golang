package gtscenter

import (
	clog "components/log/zaplogger"
	"components/sdks/gamefi_platform"
	"components/sdks/gts_shop"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func getGtsCenter() *GtsCenter {
	return NewGtsCenter(
		gts_shop.NewClient("http://10.40.10.10:30011", 5*time.Second, "sdks", initLogger()),
		gamefi_platform.NewClient("http://10.40.10.10:30005", 5*time.Second, "sdks", initLogger()),
		initLogger(),
	)
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

func Test_Gp(t *testing.T) {
	gtscenter := getGtsCenter()

	var id = time.Now().Unix()
	err := gtscenter.FreezeAsset(context.Background(), "1481130751615045", "gp", "1", id)
	fmt.Println("Res:", err)
	fmt.Println("")

	err = gtscenter.ReturnAsset(context.Background(), "1481130751615045", "gp", "1", id)
	fmt.Println("Res:", err)
	fmt.Println("")
}

func Test_Dust(t *testing.T) {
	gtscenter := getGtsCenter()

	var id = time.Now().Unix()
	err := gtscenter.FreezeAsset(context.Background(), "1481130751615045", "dust", "1", id)
	fmt.Println("Res:", err)
	fmt.Println("")

	err = gtscenter.ReturnAsset(context.Background(), "1481130751615045", "dust", "1", id)
	fmt.Println("Res:", err)
	fmt.Println("")

}
