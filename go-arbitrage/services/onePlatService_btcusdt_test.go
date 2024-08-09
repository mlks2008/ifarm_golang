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
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	_ "goarbitrage/api/router"
	"goarbitrage/internal/model"
	"goarbitrage/internal/plat"
	"strconv"
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

// 交易所内同一时间挂买一卖一
func Test_OnePlat_BTCUSDT_Start(t *testing.T) {
	//oneplatService := NewOnePlatService(plat.P_Binance)
	//go oneplatService.Start(symbol, 1, 3*60, 5*60)
	//
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

// 取消挂单，重置资产
func Test_OnePlat_BTCUSDT_Porfit(t *testing.T) {
	symbol := "BTC/USDT"
	symbolBase := strings.Split(symbol, "/")[0]
	symbolQuote := strings.Split(symbol, "/")[1]

	onePlatService := NewOnePlatService(plat.P_Binance)
	//初始投入
	initInput := onePlatService.p.GetInitialInput(symbol)
	//当前余额
	run_bBalance, _ := onePlatService.p.GetAccountBalance(symbol, true)

	p1, _ := onePlatService.p.GetMarketDepth(symbol)
	var price, _ = decimal.NewFromFloat(p1.Tick.Bids[0][0]).DivRound(decimal.NewFromFloat(1), onePlatService.getSymbol(symbol).PricePrecision).Float64()
	//price = 23355.0

	accountProfit, accountChangeBase := onePlatService.profitFunc(price, initInput, run_bBalance, symbolBase, symbolQuote)
	fmt.Println(price, accountProfit, accountChangeBase)

	if accountChangeBase.Cmp(decimal.NewFromFloat(0)) > 0 {
		//orderid, err := onePlatService.p.OrdersPlace(symbol, fmt.Sprintf("%v", price), accountChangeBase.String(), "SELL")
		//fmt.Println(orderid, err)

	} else if accountChangeBase.Cmp(decimal.NewFromFloat(0)) < 0 {
		//orderid, err := onePlatService.p.OrdersPlace(symbol, fmt.Sprintf("%v", price), decimal.NewFromFloat(0).Sub(accountChangeBase).String(), "BUY")
		//fmt.Println(orderid, err)
	}
	time.Sleep(time.Second)
}

// 当前挂单每单均价
func Test_OnePlat_BTCUSDT_OpenOrders_Merge(t *testing.T) {
	var symbol = "BTC/USDT"

	plat := plat.Get(plat.P_Binance)
	openorders, _ := plat.OpenOrders(symbol)

	var cancelsell_SumUSDT decimal.Decimal
	var cancelsell_SumQuantity decimal.Decimal
	var cancelsell_SumCount int32
	var cancelbuy_SumUSDT decimal.Decimal
	var cancelbuy_SumQuantity decimal.Decimal
	var cancelbuy_SumCount int32
	for i, open := range openorders {
		price, err := decimal.NewFromString(open.Price)
		if err != nil {
			panic(err)
		}
		origQuantity, err := decimal.NewFromString(open.OrigQuantity)
		if err != nil {
			panic(err)
		}

		if i < 200 && open.Side == "SELL" {
			//要取消的挂单
			if false {
				if _, err := plat.CancelOrder(symbol, fmt.Sprintf("%v", open.OrderID)); err == nil {
					cancelsell_SumUSDT = cancelsell_SumUSDT.Add(price.Mul(origQuantity))
					cancelsell_SumQuantity = cancelsell_SumQuantity.Add(origQuantity)
					cancelsell_SumCount++
					fmt.Println(i+1, "取消数:", cancelsell_SumCount, "总USDT:", cancelsell_SumUSDT.String(), "总BTC:", cancelsell_SumQuantity.String(), "合并订单价:", cancelsell_SumUSDT.DivRound(cancelsell_SumQuantity, 2))
				} else {
					fmt.Println(i+1, "取消失败", err)
				}
			} else {
				cancelsell_SumUSDT = cancelsell_SumUSDT.Add(price.Mul(origQuantity))
				cancelsell_SumQuantity = cancelsell_SumQuantity.Add(origQuantity)
				cancelsell_SumCount++
				//fmt.Println(i+1, "预计取消数:", cancelsell_SumCount, "总USDT:", cancelsell_SumUSDT.String(), "总BTC:", cancelsell_SumQuantity.String(), "合并订单价:", cancelsell_SumUSDT.DivRound(cancelsell_SumQuantity, 2))
			}
		}
		if i < 200 && open.Side == "BUY" {
			//要取消的挂单
			if false {
				if _, err := plat.CancelOrder(symbol, fmt.Sprintf("%v", open.OrderID)); err == nil {
					cancelbuy_SumUSDT = cancelbuy_SumUSDT.Add(price.Mul(origQuantity))
					cancelbuy_SumQuantity = cancelbuy_SumQuantity.Add(origQuantity)
					cancelbuy_SumCount++
					fmt.Println(i+1, "取消数:", cancelbuy_SumCount, "总USDT:", cancelbuy_SumUSDT.String(), "总BTC:", cancelbuy_SumQuantity.String(), "合并订单价:", cancelbuy_SumUSDT.DivRound(cancelbuy_SumQuantity, 2))
				} else {
					fmt.Println(i+1, "取消失败", err)
				}
			} else {
				cancelbuy_SumUSDT = cancelbuy_SumUSDT.Add(price.Mul(origQuantity))
				cancelbuy_SumQuantity = cancelbuy_SumQuantity.Add(origQuantity)
				cancelbuy_SumCount++
				//fmt.Println(i+1, "预计取消数:", cancelbuy_SumCount, "总USDT:", cancelbuy_SumUSDT.String(), "总BTC:", cancelbuy_SumQuantity.String(), "合并订单价:", cancelbuy_SumUSDT.DivRound(cancelbuy_SumQuantity, 2))
			}
		}
	}

	fmt.Println("订单数:", len(openorders))
	fmt.Println("Sell预计取消数:", cancelsell_SumCount, "订单平均价值:", cancelsell_SumUSDT.DivRound(decimal.NewFromInt(int64(cancelsell_SumCount)), 2).String(), "总USDT:", cancelsell_SumUSDT.String(), "总BTC:", cancelsell_SumQuantity.String(), "合并订单价:", cancelsell_SumUSDT.DivRound(cancelsell_SumQuantity, 2))
	fmt.Println("Buy预计取消数:", cancelbuy_SumCount, "订单平均价值:", cancelbuy_SumUSDT.DivRound(decimal.NewFromInt(int64(cancelbuy_SumCount)), 2).String(), "总USDT:", cancelbuy_SumUSDT.String(), "总BTC:", cancelbuy_SumQuantity.String(), "合并订单价:", cancelbuy_SumUSDT.DivRound(cancelbuy_SumQuantity, 2))
	fmt.Println("")
}

// 当前挂单亏损统计
func Test_OnePlat_BTCUSDT_OpenOrders(t *testing.T) {
	var symbol = "BTC/USDT"

	plat := plat.Get(plat.P_Binance)
	openorders, _ := plat.OpenOrders(symbol)
	fmt.Println("订单数:", len(openorders))

	p1, err := plat.GetMarketDepth(symbol)
	if err != nil {
		panic(err)
	}
	var curPrice = decimal.NewFromFloat(p1.Tick.Bids[0][0])
	var totalLossUsdt decimal.Decimal
	for i, open := range openorders {
		price, err := decimal.NewFromString(open.Price)
		if err != nil {
			panic(err)
		}
		origQuantity, err := decimal.NewFromString(open.OrigQuantity)
		if err != nil {
			panic(err)
		}

		//SELL
		if open.Side == "SELL" && price.Cmp(decimal.NewFromFloat(0*10000)) > 0 {
			downPrice := price.Sub(curPrice).DivRound(decimal.NewFromFloat(1), 2)
			downRate, _ := downPrice.DivRound(price, 5).Mul(decimal.NewFromFloat(100)).Float64()
			lossusdtamount := downPrice.Mul(origQuantity).DivRound(decimal.NewFromFloat(1), 5)
			totalLossUsdt = totalLossUsdt.Add(lossusdtamount)
			rescuelog := fmt.Sprintf("时间:%v 订单:%v(%v) \t%v挂单 下降:%v(%.3f%%) 亏损:%vu \trescuelog: %v,%v,%v,%v,%v,%v",
				time.Unix(open.Time/1000, 0).Format("2006-01-02 15:04:05"), origQuantity.String(), open.Price, open.Side,
				downPrice, downRate, lossusdtamount.String(),
				open.OrderID, price.String(), "FILLED", price.Sub(decimal.NewFromFloat(4.8)).String(), origQuantity.String(), price.Mul(origQuantity).DivRound(decimal.NewFromInt(1), 2).String())
			fmt.Println(i+1, rescuelog)
		}
		//BUY
		if open.Side == "BUY" && price.Cmp(decimal.NewFromFloat(10*10000)) < 0 {
			upPrice := curPrice.Sub(price).DivRound(decimal.NewFromFloat(1), 2)
			upRate, _ := upPrice.DivRound(price, 5).Mul(decimal.NewFromFloat(100)).Float64()
			lossusdtamount := upPrice.Mul(origQuantity).DivRound(decimal.NewFromFloat(1), 5)
			totalLossUsdt = totalLossUsdt.Add(lossusdtamount)
			rescuelog := fmt.Sprintf("时间:%v 订单:%v(%v) \t%v挂单  上升:%v(%.3f%%) 亏损:%vu \trescuelog: %v,%v,%v,%v,%v,%v",
				time.Unix(open.Time/1000, 0).Format("2006-01-02 15:04:05"), origQuantity.String(), open.Price, open.Side,
				upPrice, upRate, lossusdtamount.String(),
				"FILLED", price.Add(decimal.NewFromFloat(4.8)).String(), open.OrderID, price.String(), origQuantity.String(), price.Mul(origQuantity).DivRound(decimal.NewFromInt(1), 2).String())
			fmt.Println(i+1, rescuelog)
		}
	}

	fmt.Println("")
	fmt.Println("订单数:", len(openorders), "总亏损:", totalLossUsdt.String())
	fmt.Println("")
}

// 根据日志批量救单
func Test_OnePlat_BTCUSDT_TryRescueOrderByLogs(t *testing.T) {
	var (
		logs          string
		loglines      []string
		rescuelogkey  = make(map[string]string)
		rescueloglist = make([]string, 0)

		p = plat.Get(plat.P_Binance)
	)

	var symbol = "BTC/USDT"

	var findSymbol = NewOnePlatService(plat.P_Binance).getSymbol(symbol)

	logs = `
时间:2023-02-08 10:02:37 订单:0.00064(23306.10000000) 	SELL挂单 下降:1499.28(6.433%) 亏损:0.95954u 	rescuelog: 18351855501,23306.1,FILLED,23301.3,0.00064,15
时间:2023-02-08 10:02:47 订单:0.00064(23306.38000000) 	SELL挂单 下降:1499.56(6.434%) 亏损:0.95972u 	rescuelog: 18351860159,23306.38,FILLED,23301.58,0.00064,15
时间:2023-02-08 10:02:52 订单:0.00064(23306.20000000) 	SELL挂单 下降:1499.38(6.433%) 亏损:0.9596u 	rescuelog: 18351862223,23306.2,FILLED,23301.4,0.00064,15
时间:2023-02-08 10:28:15 订单:0.00064(23298.94000000) 	SELL挂单 下降:1492.12(6.404%) 亏损:0.95496u 	rescuelog: 18352577228,23298.94,FILLED,23294.14,0.00064,15
时间:2023-02-08 10:28:16 订单:0.00064(23298.97000000) 	SELL挂单 下降:1492.15(6.404%) 亏损:0.95498u 	rescuelog: 18352577490,23298.97,FILLED,23294.17,0.00064,15
时间:2023-02-08 10:28:16 订单:0.00064(23298.97000000) 	SELL挂单 下降:1492.15(6.404%) 亏损:0.95498u 	rescuelog: 18352577491,23298.97,FILLED,23294.17,0.00064,15
时间:2023-02-08 10:28:42 订单:0.00064(23298.72000000) 	SELL挂单 下降:1491.9(6.403%) 亏损:0.95482u 	rescuelog: 18352590123,23298.72,FILLED,23293.92,0.00064,15
时间:2023-02-08 10:40:31 订单:0.00064(23294.42000000) 	SELL挂单 下降:1487.6(6.386%) 亏损:0.95206u 	rescuelog: 18352908160,23294.42,FILLED,23289.62,0.00064,15
时间:2023-02-08 12:10:47 订单:0.00064(23281.63000000) 	SELL挂单 下降:1474.81(6.335%) 亏损:0.94388u 	rescuelog: 18355352845,23281.63,FILLED,23276.83,0.00064,15
时间:2023-02-08 12:10:47 订单:0.00064(23281.63000000) 	SELL挂单 下降:1474.81(6.335%) 亏损:0.94388u 	rescuelog: 18355352849,23281.63,FILLED,23276.83,0.00064,15

`
	loglines = strings.Split(logs, "\n")
	//找出要救订单
	for _, line := range loglines {
		if strings.Index(line, "(rescuelog1:") > 0 {
			tmp1 := strings.TrimSpace(strings.Split(line, "(rescuelog1:")[1])
			rescuelog1 := strings.Split(tmp1, " ")[0]
			if _, exist := rescuelogkey[rescuelog1]; exist == false {
				rescuelogkey[rescuelog1] = ""
				rescueloglist = append(rescueloglist, rescuelog1)
			}
		} else if strings.Index(line, "rescuelog:") > 0 {
			tmp1 := strings.TrimSpace(strings.Split(line, "rescuelog:")[1])
			rescuelog := strings.Split(tmp1, " ")[0]
			if _, exist := rescuelogkey[rescuelog]; exist == false {
				rescuelogkey[rescuelog] = ""
				rescueloglist = append(rescueloglist, rescuelog)
			}
		}
	}
	//开始救单
	var group int
	var groupCount int = len(rescueloglist)
	for _, rescuelog := range rescueloglist {
		group++
		p1, err := p.GetMarketDepth(symbol)
		if err != nil {
			time.Sleep(time.Millisecond * 1000)
			continue
		}
		var buyPrice, _ = decimal.NewFromFloat(p1.Tick.Bids[0][0]-5.5).DivRound(decimal.NewFromFloat(1), findSymbol.PricePrecision).Float64()
		var sellPrice, _ = decimal.NewFromFloat(p1.Tick.Asks[0][0]+4).DivRound(decimal.NewFromFloat(1), findSymbol.PricePrecision).Float64()
		btime := time.Now().Unix()
		_tryRescueOrder(p, findSymbol, group, groupCount, symbol, buyPrice, sellPrice, rescuelog)
		if time.Now().Unix()-btime > 5 {
			//刚救完，过几秒在开
			time.Sleep(time.Second * 3)
		}
	}
}

func Test_OnePlat_BTCUSDT_TryRescueOrder1(t *testing.T) {
	var symbol = "BTC/USDT"

	p := plat.Get(plat.P_Binance)

	var findSymbol = NewOnePlatService(plat.P_Binance).getSymbol(symbol)

	p1, _ := p.GetMarketDepth(symbol)
	buyPrice, _ := decimal.NewFromFloat(p1.Tick.Bids[0][0]-2).DivRound(decimal.NewFromFloat(1), findSymbol.PricePrecision).Float64()
	sellPrice, _ := decimal.NewFromFloat(p1.Tick.Asks[0][0]+2).DivRound(decimal.NewFromFloat(1), findSymbol.PricePrecision).Float64()

	_tryRescueOrder(p, findSymbol, 1, 1, symbol, buyPrice, sellPrice, "17857436103,22992.81,FILLED,22992,0.00065,15")
}

func Test_OnePlat_BTCUSDT_TryRescueOrder2(t *testing.T) {
	var symbol = "BTC/USDT"

	p := plat.Get(plat.P_Binance)

	var findSymbol = NewOnePlatService(plat.P_Binance).getSymbol(symbol)

	p1, _ := p.GetMarketDepth(symbol)
	buyPrice, _ := decimal.NewFromFloat(p1.Tick.Bids[0][0]-2).DivRound(decimal.NewFromFloat(1), findSymbol.PricePrecision).Float64()
	sellPrice, _ := decimal.NewFromFloat(p1.Tick.Asks[0][0]+2).DivRound(decimal.NewFromFloat(1), findSymbol.PricePrecision).Float64()

	_tryRescueOrder(p, findSymbol, 1, 1, symbol, buyPrice, sellPrice, "17857436103,22992.81,FILLED,22992,0.00065,15")
}

func Test_OnePlat_BTCUSDT_Order(t *testing.T) {
	plat := plat.Get(plat.P_Binance)
	fmt.Println(plat.CancelOrder("BTC/USDT", "17889274608"))
}

func Test_OnePlat_BTCUSDT_Buy(t *testing.T) {
	plat := plat.Get(plat.P_Binance)
	//todo Buy预计取消数: 20 订单平均价值: 14.65 总USDT: 293.0687079 总BTC: 0.01334 合并订单价: 21969.17
	var buyprice = "21900"
	var btcamount = "0.01334"
	orderid, err := plat.OrdersPlace("BTC/USDT", buyprice, btcamount, "BUY")
	fmt.Println(orderid, err)
	time.Sleep(time.Second)
}

func Test_OnePlat_BTCUSDT_Sell(t *testing.T) {
	plat := plat.Get(plat.P_Binance)

	var sellprice = "23259.55"
	var btcamount = "0"
	orderid, err := plat.OrdersPlace("BTC/USDT", sellprice, btcamount, "SELL")
	//余额不足
	commerr, ok := err.(*common.APIError)
	if ok && commerr.Code == -2010 {
		fmt.Println("余额不足", err)
	} else {
		fmt.Println(orderid)
	}
	time.Sleep(time.Second)
}

func _tryRescueOrder(p plat.Plat, findSymbol model.SymbolsData, group int, groupCount int, symbol string, buyprice, sellprice float64, rescuelog string) {
	oneplatService := NewOnePlatService(p.PlatCode())
	oneplatService.manual_price_enable = true      //开启
	oneplatService.manual_curbuyprice = buyprice   //人工救单卖买着价
	oneplatService.manual_cursellprice = sellprice //人工救单卖买着价

	//sellOrderId, sellPrice, buyOrderId, buyPrice, btcAmount, usdtAmount
	orderinfo := strings.Split(rescuelog, ",")

	var sellorderid = strings.Replace(orderinfo[0], "TIMEOUT", "", 1)
	var lastsellprice, _ = strconv.ParseFloat(orderinfo[1], 64)
	var buyorderid = strings.Replace(orderinfo[2], "TIMEOUT", "", 1)
	var lastbuyprice, _ = strconv.ParseFloat(orderinfo[3], 64)
	var lastbtcamount, _ = strconv.ParseFloat(orderinfo[4], 64)
	var lastusdtamount, _ = strconv.ParseFloat(orderinfo[5], 64)

	if oneplatService.manual_curbuyprice >= oneplatService.manual_cursellprice {
		panic("买价大于卖价")
		return
	}
	//oneplatService.rescueMap.Store(oneplatService.genRotId(group, symbol, 3*1000-1), time.Now().Unix())
	oneplatService.orderPrice.Store(oneplatService.genRotId(group, symbol, 3*1000-1), &OrderInfo{SellOrderId: &sellorderid, SellPrice: oneplatService.manual_cursellprice, BuyOrderId: &buyorderid, BuyPrice: oneplatService.manual_curbuyprice})

	fmt.Println("****", fmt.Sprintf("%v/%v", group, groupCount), "开始救单", rescuelog, time.Now().Format("2006-01-02 T15:04:05"))
	ok, _ := oneplatService.rescueLastOrder(3*1000, group, 24*60*60, symbol, findSymbol, &sellorderid, lastsellprice, &buyorderid, lastbuyprice, lastbtcamount, lastusdtamount)
	if ok {
		fmt.Println(ok, fmt.Sprintf("%v/%v", group, groupCount), "救单成功", rescuelog, time.Now().Format("2006-01-02 T15:04:05"))
	} else {
		fmt.Println(ok, fmt.Sprintf("%v/%v", group, groupCount), "救单失败", rescuelog, time.Now().Format("2006-01-02 T15:04:05"))
		fmt.Println("")
		fmt.Println("")
	}
	time.Sleep(time.Millisecond * 200)
}

func Test_OnePlat_BTCUSDT_Buy2(t *testing.T) {
	plat := plat.Get(plat.P_Binance)

	var symbol = "FDUSD/USDT"
	var buyprice = "0.9998"
	var amount = "12879"

	orderid, err := plat.OrdersPlace2("ARB/FDUSD", symbol, buyprice, amount, "BUY")
	if err != nil {
		panic(err)
	} else {
		t.Log(orderid)
	}
}
