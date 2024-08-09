package services

import (
	"components/log/log"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	_ "goarbitrage/api/router"
	"goarbitrage/internal/model"
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

func Test_OnePlat_btc300(t *testing.T) {
	var binance = plat.Get(plat.P_Binance)
	var symbol = "BTC/FDUSD"

	/*
	* 捕单
	* 每0.2%挂一单,每单金额增长1.2
	* 10000
	* 10000*1.2
	* 10000*1.2 * 1.2
	 */
	run_bBalance, err0 := binance.GetAccountBalance(symbol, true)
	if err0 != nil {
		panic(err0)
	}
	var basePrice, err = GetBuyPrice(binance, symbol)
	if err != nil {
		panic(err)
	}
	var basePriceDown = 0.0025
	var baseUsdt = decimal.NewFromFloat(5000.0)
	var baseUsdtInc = 1.2

	var totalCost = decimal.Zero
	var totalAmount = decimal.Zero
	for i := 0; i < 6; i++ {
		var price = basePrice.Mul(decimal.NewFromFloat(1-basePriceDown*float64(i))).DivRound(decimal.NewFromFloat(1), 2)
		var cost = baseUsdt.Mul(decimal.NewFromFloat(power(baseUsdtInc, i)))
		var amount = cost.Div(price)
		amount = amount.DivRound(decimal.NewFromInt(1), 5)
		for {
			_, err := binance.OrdersPlace(symbol, fmt.Sprintf("%v", price), fmt.Sprintf("%v", amount), "BUY")
			if err == nil {
				break
			} else {
				time.Sleep(time.Second * 5)
			}
		}
		totalCost = totalCost.Add(cost)
		totalAmount = totalAmount.Add(amount)
		fmt.Println(i, "\t", price.DivRound(decimal.NewFromFloat(1), 2), "\t", cost, "\t", amount, fmt.Sprintf("\t总成本:%v,\t均价:%v", totalCost.String(), totalCost.Div(totalAmount).String()))
	}

	/*
	* 计算收益
	* 当前价值 - 初始余额 > 0.3%
	 */
	for {
		time.Sleep(time.Second * 30)
		bot_eBalance, err0 := binance.GetAccountBalance(symbol, true)
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
		profit, profitRate, _ := profitFunc(sellPriceFloat64, run_bBalance, bot_eBalance, symbolBase, symbolQuote)
		fmt.Println("profit", profit, profitRate)

		if profitRate.GreaterThanOrEqual(decimal.NewFromFloat(0.003)) {
			binance.CancelOpenOrders(symbol)
			time.Sleep(time.Second * 3)

			sellbtcamount := bot_eBalance[symbolBase].Total.DivRound(decimal.NewFromInt(1), 5).String()
			_, err = binance.OrdersPlace(symbol, fmt.Sprintf("%v", sellPriceFloat64), fmt.Sprintf("%v", sellbtcamount), "SELL")

			fmt.Println("已挂卖单，待成交")
		}
	}
}

func profitFunc(sellPrice float64, iBalance map[string]model.Balance, eBalance map[string]model.Balance, base, quote string) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	//多出币换算为等值U
	changeBase := eBalance[base].Total.Sub(iBalance[base].Total)
	changeBaseToQuote := changeBase.Mul(decimal.NewFromFloat(sellPrice))
	//当前最新U余额
	nowQuoteBalance := eBalance[quote].Total.Add(changeBaseToQuote)
	//收益
	profit := nowQuoteBalance.Sub(iBalance[quote].Total).DivRound(decimal.NewFromFloat(1), 4)
	//收益率
	profitRate := nowQuoteBalance.Sub(iBalance[quote].Total).Div(iBalance[quote].Total)
	return profit, profitRate, changeBase
}

func power(base float64, exponent int) float64 {
	if exponent == 0 {
		return 1
	}
	if exponent < 0 {
		return 1 / power(base, -exponent)
	}
	return base * power(base, exponent-1)
}

func GetBuyPrice(binance plat.Plat, symbol string) (decimal.Decimal, error) {
	p1, err := binance.GetMarketDepth(symbol)
	if err != nil {
		return decimal.Zero, err
	}
	if true {
		//挂单
		if p1.Tick.Bids[0][0] == 0 || len(p1.Tick.Bids) < 1 {
			return decimal.Zero, fmt.Errorf("bids is invalid")
		}
		return decimal.NewFromFloat(p1.Tick.Bids[0][0]), nil
	} else {
		//吃单
		if p1.Tick.Asks[0][0] == 0 || len(p1.Tick.Asks) < 1 {
			return decimal.Zero, fmt.Errorf("asks is invalid")
		}
		return decimal.NewFromFloat(p1.Tick.Asks[0][0]), nil
	}
}

func GetSellPrice(binance plat.Plat, symbol string) (decimal.Decimal, error) {
	p1, err := binance.GetMarketDepth(symbol)
	if err != nil {
		return decimal.Zero, err
	}

	if true {
		//挂单
		if p1.Tick.Asks[0][0] == 0 || len(p1.Tick.Asks) < 1 {
			return decimal.Zero, fmt.Errorf("asks is invalid")
		}
		return decimal.NewFromFloat(p1.Tick.Asks[0][0]), nil
	} else {
		//吃单
		if p1.Tick.Bids[0][0] == 0 || len(p1.Tick.Bids) < 1 {
			return decimal.Zero, fmt.Errorf("bids is invalid")
		}
		return decimal.NewFromFloat(p1.Tick.Bids[0][0]), nil
	}
}
