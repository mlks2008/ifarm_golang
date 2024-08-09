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
	"github.com/shopspring/decimal"
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

func Test_DCA_FILUSDT(t *testing.T) {
	if true {
		dcaService := NewDCAService(plat.P_Binance, "FIL/USDT")
		fmt.Println(dcaService.Buy(decimal.NewFromFloat(0), decimal.NewFromFloat(3)))
		time.Sleep(time.Second)
	} else {
		dcaService := NewDCAService(plat.P_Binance, "FIL/USDT")
		fmt.Println(dcaService.Sell(decimal.NewFromFloat(0), decimal.NewFromFloat(10)))
		time.Sleep(time.Second)
	}
}

func Test_DCA_FDUSDUSDT(t *testing.T) {
	//dcaService := NewDCAService(plat.P_Binance, "FDUSD/USDT")
	//fmt.Println(dcaService.Sell(decimal.NewFromFloat(10000), decimal.NewFromFloat(0.9996)))
	//time.Sleep(time.Second)

	dcaService := NewDCAService(plat.P_Binance, "FDUSD/USDT")
	fmt.Println(dcaService.Buy(decimal.NewFromFloat(47800), decimal.NewFromFloat(0)))
	time.Sleep(time.Second)
}

func Test_DCA_MEMEUSDT(t *testing.T) {
	dcaService := NewDCAService(plat.P_Binance, "MEME/USDT")
	fmt.Println(dcaService.Sell(decimal.NewFromFloat(538), decimal.NewFromFloat(0.0292)))
	time.Sleep(time.Second)
}

func Test_DCA_ARBUSDT(t *testing.T) {
	if true {
		dcaService := NewDCAService(plat.P_Binance, "ARB/USDT")
		fmt.Println(dcaService.Buy(decimal.NewFromFloat(13552), decimal.NewFromFloat(1.106)))
		time.Sleep(time.Second)
	} else {
		dcaService := NewDCAService(plat.P_Binance, "ARB/USDT")
		fmt.Println(dcaService.Sell(decimal.NewFromFloat(0), decimal.NewFromFloat(0)))
		time.Sleep(time.Second)
	}
}
