package main

import (
	"components/common/global"
	clog "components/log/zaplogger"
	"context"
	"flag"
	"fmt"
	"gamefi_equipment/internal/conf"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "dev", "config path, eg: -conf dev")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(fmt.Sprintf("configs/%v", flagconf)),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	logger := initLogger(bc.Log)

	//使用配制name
	Name = bc.Log.Name
	if bc.Log.Name != "" {
		conf.Prefix = bc.Log.Name
	}

	// 初始化参数
	conf.FlagConf = flagconf
	// 初始化global
	global.Init(logger)

	//// register
	//serviceHttpPort, _ := strconv.ParseInt(strings.Split(bc.Server.Http.GetAddr(), ":")[1], 10, 64)
	//nacosClient := nacos.NewClient(bc.Nacos.Ip, bc.Nacos.Port, bc.Nacos.UserName, bc.Nacos.Pwd, bc.Nacos.NamespaceId, bc.Nacos.TimeoutMs, bc.Log.Name, bc.Log.Name, bc.Log.Dir, true)
	//_, err := nacosClient.Register(serviceHttpPort)
	//if err != nil {
	//	panic(err)
	//}
	//defer nacosClient.Stop()

	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.Kafka, bc.Client, bc.Sys, bc.Log, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func initLogger(c *conf.Log) log.Logger {
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
	return log.With(clog.NewLoggerWithName(c.Dir, c.Name, "", c.Debug),
		"ts", log.DefaultTimestamp,
		"caller", caller(4),
		"service.id", id,
		"service.name", c.Name,
		"service.version", Version,
	)
}
