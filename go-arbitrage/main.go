package main

import (
	"components/message"
	"flag"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	_ "goarbitrage/api/router"
	"goarbitrage/internal/gvar"
	"goarbitrage/internal/plat"
	"goarbitrage/services"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

var configFile string
var cmd string

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	flag.StringVar(&configFile, "config", "configs/hotbit_dev.yaml", "")
	flag.StringVar(&cmd, "cmd", "", "")
	flag.Parse()

	message.SendDingTalkRobit(true, "all", "typ", "txid", "content")

	if configFile != "configs/empty.yaml" {
		//viper解析配置文件
		ext := filepath.Ext(configFile)
		if strings.ToLower(ext) != ".yaml" {
			panic("config仅支持yaml格式")
		}
		viper.SetConfigType(strings.ToLower(ext)[1:])
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&gvar.G_Config); err != nil {
			panic(err)
		}
	}

	//l4g.LoadConfiguration(gvar.G_Config.LogerXml)
	//defer l4g.Close()

	//注册服务
	plat.Register()

	if cmd == "arb" {
		dcaService := services.NewDCAService(plat.P_Binance, "ARB/USDT")

		for i := 1; i <= 10; i++ {
			buyPoint, profitTargetPoint := dcaService.GetBuyPoint(i, time.Now().Unix())
			fmt.Println("ARB/USDT", i, ":", buyPoint, profitTargetPoint)
		}
		_, _, firstBuyUsdt := dcaService.GetBuyTimes(decimal.NewFromInt(1), decimal.NewFromInt(1))
		fmt.Println("ARB/USDT", "firstBuyUsdt", ":", firstBuyUsdt)

		dcaService.Start()
	} else if cmd == "bnb" {
		dcaService := services.NewDCAService(plat.P_Binance, "BNB/FDUSD")

		for i := 1; i <= 10; i++ {
			buyPoint, profitTargetPoint := dcaService.GetBuyPoint(i, time.Now().Unix())
			fmt.Println("BNB/FDUSD", i, ":", buyPoint, profitTargetPoint)
		}
		_, _, firstBuyUsdt := dcaService.GetBuyTimes(decimal.NewFromInt(1), decimal.NewFromInt(1))
		fmt.Println("BNB/FDUSD", "firstBuyUsdt", ":", firstBuyUsdt)

		dcaService.Start()
	} else if cmd == "oneplat" {
		onePlatService := services.NewOnePlatService(plat.P_Binance)
		go onePlatService.Start("DOGE/FDUSD", 1, 15*60, 5*60)
	} else {
		panic(fmt.Sprintf("无效cmd:%v", cmd))
	}

	////单交易所快速买卖策略
	//onePlatService1 := services.NewOnePlatService(plat.P_Binance)
	//go onePlatService1.Start("BTC/USDT", 5, 5*60, 5*60)
	//
	//onePlatService2 := services.NewOnePlatService(plat.P_Binance)
	//go onePlatService2.Start("BTC/BUSD", 3, 15*60, 5*60)

	//onePlatService3 := services.NewOnePlatService(plat.P_Binance)
	//go onePlatService3.Start("ETH/USDT", 5, 5*60, 10*60)
	//
	//onePlatService4 := services.NewOnePlatService(plat.P_Binance)
	//go onePlatService4.Start("BNBDOWN/USDT", 5, 10*60, 5*60)

	//三角套利
	//threeService := services.NewThreeService(plat.P_KuCoin)
	//go threeService.Start()

	////多交易所搬砖套利
	//brickMoverService := services.NewBrickMoverService(plat.P_Binance, plat.P_KuCoin)
	//go brickMoverService.Start()

	////关闭printlog
	//api.PrintLog = false

	//if os.Getenv("TEST") == "1" {
	//	httpserver.SetHandlerChan(10)
	//	httpserver.SetHttpFormat("jsonrpc")
	//	httpserver.RunApiServer(10, 10, 5, ":16868")
	//} else {
	//	httpserver.SetHandlerChan(10)
	//	httpserver.SetHttpFormat("jsonrpc")
	//	httpserver.RunApiServer(10, 10, 5, ":6868")
	//}
	select {}
}
