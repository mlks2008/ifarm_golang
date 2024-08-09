package services

import (
	"components/log/log"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	_ "goarbitrage/api/router"
	"goarbitrage/internal/plat"
	"strings"
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

func Test_OnePlat_doge(t *testing.T) {
	var binance = plat.Get(plat.P_Binance)
	var symbol = "DOGE/FDUSD"

	begin_Balance, err0 := binance.GetAccountBalance(symbol, true)
	if err0 != nil {
		panic(err0)
	}
	var basePrice, err = GetBuyPrice(binance, symbol)
	if err != nil {
		panic(err)
	}
	var basePriceUp = 0.002
	var baseDoge = decimal.NewFromFloat(10000.0)
	var baseUsdtInc = 1.0

	var totalCost = decimal.Zero
	var totalAmount = decimal.Zero
	for i := 0; i < 6; i++ {
		var price = basePrice.Mul(decimal.NewFromFloat(1+basePriceUp*float64(i))).DivRound(decimal.NewFromFloat(1), 2)
		var amount = baseDoge.Mul(decimal.NewFromFloat(power(baseUsdtInc, i)))
		amount = amount.DivRound(decimal.NewFromInt(1), 0)
		for {
			_, err := binance.OrdersPlace(symbol, fmt.Sprintf("%v", price), fmt.Sprintf("%v", amount), "SELL")
			if err == nil {
				break
			} else {
				time.Sleep(time.Second * 5)
			}
		}
		totalCost = totalCost.Add(amount.Mul(price))
		totalAmount = totalAmount.Add(amount)
		fmt.Println(i, "\t", price.DivRound(decimal.NewFromFloat(1), 2), "\t", amount, "\t", amount, fmt.Sprintf("\t总成本:%v,\t均价:%v", totalCost.String(), totalCost.Div(totalAmount).String()))
	}

	/*
	* 计算收益
	* 当前价值 - 初始余额 > 0.3%
	 */
	for {
		time.Sleep(time.Second * 15)
		end_Balance, err0 := binance.GetAccountBalance(symbol, true)
		if err0 != nil {
			log.Logger.Error(err0)
			time.Sleep(time.Second * 10)
			continue
		}
		var sellPrice, err = GetSellPrice(binance, symbol)
		if err != nil {
			panic(err)
		}
		var sellPriceFloat64, _ = sellPrice.Float64()
		symbolBase := strings.Split(symbol, "/")[0]
		symbolQuote := strings.Split(symbol, "/")[1]
		profit, profitRate, _ := profitFunc(sellPriceFloat64, begin_Balance, end_Balance, symbolBase, symbolQuote)
		fmt.Println("profit", profit, profitRate)

		if profitRate.GreaterThanOrEqual(decimal.NewFromFloat(0.003)) {
			binance.CancelOpenOrders(symbol)
			time.Sleep(time.Second * 3)

			sellbtcamount := end_Balance[symbolBase].Total.DivRound(decimal.NewFromInt(1), 5).String()
			_, err = binance.OrdersPlace(symbol, fmt.Sprintf("%v", sellPriceFloat64), fmt.Sprintf("%v", sellbtcamount), "SELL")

			fmt.Println("已挂卖单，待成交")
		}
	}
}
