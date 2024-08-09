/**
 * ThreeService.go
 * ============================================================================
 * 三角套利
 * ============================================================================
 * author: peter.wang
 * createtime: 2020-07-07 22:35
 */

package services

import (
	"fmt"
	"github.com/shopspring/decimal"
	"goarbitrage/api"
	_ "goarbitrage/api/router"
	"goarbitrage/internal/plat"
	"goarbitrage/pkg/httpserver"
	"goarbitrage/pkg/utils"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_Three_BUSDUSDT_Start(t *testing.T) {
	api.OnePlatUsdtAmount1 = 30
	threeService := NewThreeService(plat.P_Binance)
	go threeService.Start()

	if os.Getenv("TEST") == "1" {
		httpserver.SetHandlerChan(10)
		httpserver.SetHttpFormat("jsonrpc")
		httpserver.RunApiServer(10, 10, 5, ":16868")
	} else {
		httpserver.SetHandlerChan(10)
		httpserver.SetHttpFormat("jsonrpc")
		httpserver.RunApiServer(10, 10, 5, ":6868")
	}
}

// 当前挂单列表
func Test_OnePlat_BUSDUSDT_OpenOrders(t *testing.T) {
	var symbol = ""

	plat := plat.Get(plat.P_Binance)
	openorders, _ := plat.OpenOrders2("BUSD/USDT", symbol)

	for i, open := range openorders {
		price, err := decimal.NewFromString(open.Price)
		if err != nil {
			panic(err)
		}
		origQuantity, err := decimal.NewFromString(open.OrigQuantity)
		if err != nil {
			panic(err)
		}

		rescuelog := fmt.Sprintf("时间:%v 订单:%v(%v) \t%v挂单 \trescuelog: %v,%v,%v,%v,%v,%v",
			time.Unix(open.Time/1000, 0).Format("2006-01-02 15:04:05"), origQuantity.String(), open.Price, open.Side,
			open.OrderID, price.String(), "FILLED", price.Sub(decimal.NewFromFloat(4.8)).String(), origQuantity.String(), price.Mul(origQuantity).DivRound(decimal.NewFromInt(1), 2).String())
		fmt.Println(i+1, rescuelog)
	}

	fmt.Println("")
	fmt.Println("订单数:", len(openorders))
	fmt.Println("")
}

// USDC/USDT稳定币套利
// 1及以上挂卖，1以下挂买
func Test_OnePlat_USDCUSDT_BuySell(t *testing.T) {
	plat := plat.Get(plat.P_Binance)

	var symbol = "USDC/USDT"
	var amount = 200

	//plat.OrdersPlace(symbol, fmt.Sprintf("%v", 1.0001), fmt.Sprintf("%v", amount), "SELL")

	//计算日化，年化率
	go func() {
		for {
			initBal := plat.GetInitialInput(symbol)
			nowBal, _ := plat.GetAccountBalance(symbol, true)
			nowDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).Unix()
			runDay := (nowDay - initBal["USDT"].AddTime) / 86400

			initTotal := initBal["USDT"].Total.Add(initBal["USDC"].Total)
			nowTotal := nowBal["USDT"].Total.Add(nowBal["USDC"].Total)
			totalProfit := nowTotal.Sub(initTotal)
			avgDayProfit := totalProfit.Div(decimal.NewFromInt(runDay))
			dayRate := avgDayProfit.Div(initTotal).Mul(decimal.NewFromInt(100))
			yearRate := dayRate.Mul(decimal.NewFromInt(365))
			fmt.Println(fmt.Sprintf("运行天数:%v, 日化:%v%%, 年化:%v%%", runDay, dayRate, yearRate))

			time.Sleep(time.Second * 12 * 3600)
		}
	}()

	for {
		depth, err := plat.GetMarketDepth(symbol)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		openorders, err := plat.OpenOrders(symbol)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		//6个买单，3个卖单
		for iprice := 10002; iprice >= 9994; {
			price := decimal.NewFromInt(int64(iprice)).Div(decimal.New(1, 4))

			//确认当前price是否有挂单
			var exist_price_order bool
			for _, open := range openorders {
				open_price, _ := decimal.NewFromString(open.Price)
				if open_price.Cmp(price) == 0 {
					exist_price_order = true
					break
				}
			}

			//不存在此价格挂单时
			if exist_price_order == false {
				//价格小于1时，挂买单
				if price.Cmp(decimal.NewFromInt(1)) < 0 && price.Cmp(decimal.NewFromFloat(depth.Tick.Asks[0][0])) < 0 {
					//买
					_, err := plat.OrdersPlace(symbol, fmt.Sprintf("%v", price), fmt.Sprintf("%v", amount), "BUY")
					if err == nil {
						fmt.Println("已挂买单", price)
					}
					//fmt.Println("将挂买单", price, "当前卖一", depth.Tick.Asks[0][0])
					time.Sleep(time.Second * 1)
				}
				//价格大于等于1时，挂卖单
				if price.Cmp(decimal.NewFromInt(1)) >= 0 && price.Cmp(decimal.NewFromFloat(depth.Tick.Bids[0][0])) >= 0 {
					//卖
					_, err := plat.OrdersPlace(symbol, fmt.Sprintf("%v", price), fmt.Sprintf("%v", amount), "SELL")
					if err == nil {
						fmt.Println("已挂卖单", price)
					}
					_, err = plat.OrdersPlace(symbol, fmt.Sprintf("%v", price), fmt.Sprintf("%v", amount), "SELL")
					if err == nil {
						fmt.Println("已挂卖单", price)
					}
					//fmt.Println("将挂卖单", price, "当前买一", depth.Tick.Bids[0][0])
					time.Sleep(time.Second * 1)
				}
			}

			iprice -= 1
		}

		time.Sleep(time.Second * 60 * 5)
	}
}

func Test_Three_BUSDUSDT_Buy(t *testing.T) {
	plat := plat.Get(plat.P_Binance)

	if true {
		symbol := "FIL/USDT"

		var buyprice = "7.6"
		var btcamount = "630"
		orderid, err := plat.OrdersPlace2("BUSD/USDT", symbol, buyprice, btcamount, "BUY")
		fmt.Println(orderid, err)
		time.Sleep(time.Second)
	} else {
		//var buyprice = "0.9998"
		//var btcamount = "0"
		//orderid, err := plat.OrdersPlace2("BUSD/USDT", "BUSD/USDT", buyprice, btcamount, "BUY")
		//fmt.Println(orderid, err)
		//time.Sleep(time.Second)
	}
}

func Test_Three_BUSDUSDT_Sell(t *testing.T) {
	plat := plat.Get(plat.P_Binance)
	balance, _ := plat.GetAccountBalance("BUSD/USDT", true)

	if true {
		symbol := "FIL/USDT"
		symbolBase := strings.Split(symbol, "/")[0]
		onePlatService := NewOnePlatService(plat.PlatCode())
		findSymbol := onePlatService.getSymbol(symbol)

		var sellprice = "8"
		var btcamount = utils.DivRound(balance[symbolBase].Free, findSymbol.AmountPrecision)
		orderid, err := plat.OrdersPlace2("BUSD/USDT", symbol, sellprice, btcamount, "SELL")
		fmt.Println(orderid, err)
		time.Sleep(time.Second)
	} else {
		//var sellprice = "1"
		//var btcamount = "0"
		//orderid, err := plat.OrdersPlace2("BUSD/USDT", "BUSD/USDT", sellprice, btcamount, "SELL")
		//fmt.Println(orderid, err)
		//time.Sleep(time.Second)
	}
}

func Test_Three_BUSDUSDT_BuySell(t *testing.T) {
	plat := plat.Get(plat.P_Binance)

	var symbol = "BTC/USDT"
	var buyprice = "22230"
	var sellprice = "22330"
	var btcamount = "0"

	orderid, err := plat.OrdersPlace2("BUSD/USDT", symbol, buyprice, btcamount, "BUY")
	if err != nil {
		panic(err)
	}
	for {
		status, _ := plat.GetOrderStatus2("BUSD/USDT", symbol, orderid)
		if status == true {
			orderid, err := plat.OrdersPlace2("BUSD/USDT", symbol, sellprice, btcamount, "SELL")
			if err != nil {
				panic(err)
			}
			for {
				status, _ := plat.GetOrderStatus2("BUSD/USDT", symbol, orderid)
				if status == true {
					fmt.Println("买卖完成")
					return
				}
			}
		}
	}
}
