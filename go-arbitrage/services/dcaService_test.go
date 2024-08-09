/**
 * onePlatService.go
 * ============================================================================
 * 3、单交易所快速买卖策略
 * ============================================================================
 * author: peter.wang
 */

package services

import (
	"fmt"
	"github.com/adshao/go-binance/v2/common"
	"github.com/spf13/viper"
	_ "goarbitrage/api/router"
	"goarbitrage/internal/plat"
	"testing"
	"time"
)

func init() {
	viper.SetConfigFile("../configs/hotbit_dev.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	plat.Register()
}

// 交易所BTC dca买卖交易
func Test_DCA_BTCUSDT_Start(t *testing.T) {
	////2023-10-1开始运行
	//dcaService := NewDCAService(plat.P_Binance, "BTC/USDT")
	//dcaService.Start()

	//makerfee, takerfee, err := dcaService.p.GetTradeFee("BTC/USDT")
	//fmt.Println("交易费", makerfee, takerfee, err)

	//if os.Getenv("TEST") == "1" {
	//	httpserver.SetHandlerChan(10)
	//	httpserver.SetHttpFormat("jsonrpc")
	//	httpserver.RunApiServer(10, 10, 5, ":16868")
	//} else {
	//	httpserver.SetHandlerChan(10)
	//	httpserver.SetHttpFormat("jsonrpc")
	//	httpserver.RunApiServer(10, 10, 5, ":6868")
	//}
}

// 交易所ARB dca买卖交易
func Test_DCA_ARBUSDT_Start(t *testing.T) {
	//2023-10-1开始运行
	dcaService := NewDCAService(plat.P_Binance, "ARB/USDT")
	dcaService.Start()
}

// 交易所BNB dca买卖交易
func Test_DCA_BNBUSDT_Start(t *testing.T) {
	//2023-11-09开始运行
	dcaService := NewDCAService(plat.P_Binance, "BNB/FDUSD")
	dcaService.Start()
}

func Test_DCA_Buy(t *testing.T) {
	plat := plat.Get(plat.P_Binance)

	var buyprice = "555"
	var btcamount = "1.15"
	orderid, err := plat.OrdersPlace("BNB/FDUSD", buyprice, btcamount, "BUY")
	//余额不足
	commerr, ok := err.(*common.APIError)
	if ok && commerr.Code == -2010 {
		fmt.Println("余额不足", err)
	} else {
		fmt.Println(orderid)
	}
	time.Sleep(time.Second)
}

func Test_DCA_Sell(t *testing.T) {
	plat := plat.Get(plat.P_Binance)

	var buyprice = "577.5"
	var btcamount = "23"
	orderid, err := plat.OrdersPlace("BNB/FDUSD", buyprice, btcamount, "SELL")
	//余额不足
	commerr, ok := err.(*common.APIError)
	if ok && commerr.Code == -2010 {
		fmt.Println("余额不足", err)
	} else {
		fmt.Println(orderid)
	}
	time.Sleep(time.Second)
}
