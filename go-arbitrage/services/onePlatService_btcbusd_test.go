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

// 取消挂单，重置资产
func Test_OnePlat_BTCBUSD_Porfit(t *testing.T) {
	symbol := "BTC/BUSD"
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
func Test_OnePlat_BTCBUSD_OpenOrders_Merge(t *testing.T) {
	var symbol = "BTC/BUSD"

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
func Test_OnePlat_BTCBUSD_OpenOrders(t *testing.T) {
	var symbol = "BTC/BUSD"

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
func Test_OnePlat_BTCBUSD_TryRescueOrderByLogs(t *testing.T) {
	var (
		logs          string
		loglines      []string
		rescuelogkey  = make(map[string]string)
		rescueloglist = make([]string, 0)

		p = plat.Get(plat.P_Binance)
	)

	var symbol = "BTC/BUSD"

	var findSymbol = NewOnePlatService(plat.P_Binance).getSymbol(symbol)

	logs = `
时间:2023-02-08 22:49:54 订单:0.04917(23096.18000000) 	SELL挂单 下降:1406.37(6.089%) 亏损:69.15121u 	rescuelog: 8675985432,23096.18,FILLED,23091.38,0.04917,15
时间:2023-02-08 22:51:45 订单:0.00065(23082.68000000) 	SELL挂单 下降:1392.87(6.034%) 亏损:0.90537u 	rescuelog: 8676034865,23082.68,FILLED,23077.88,0.00065,15
时间:2023-02-08 22:51:45 订单:0.00065(23082.68000000) 	SELL挂单 下降:1392.87(6.034%) 亏损:0.90537u 	rescuelog: 8676034867,23082.68,FILLED,23077.88,0.00065,15
时间:2023-02-09 00:51:08 订单:0.00064(23091.00000000) 	SELL挂单 下降:1401.19(6.068%) 亏损:0.89676u 	rescuelog: 8679042143,23091,FILLED,23086.2,0.00064,15
时间:2023-02-09 02:07:28 订单:0.00061(23032.61000000) 	SELL挂单 下降:1342.8(5.830%) 亏损:0.81911u 	rescuelog: 8680649470,23032.61,FILLED,23027.81,0.00061,15
时间:2023-02-09 09:48:40 订单:0.06813(22968.85000000) 	SELL挂单 下降:1279.04(5.569%) 亏损:87.141u 	rescuelog: 8687454063,22968.85,FILLED,22964.05,0.06813,15
时间:2023-02-09 10:32:45 订单:0.00065(22946.13000000) 	SELL挂单 下降:1256.32(5.475%) 亏损:0.81661u 	rescuelog: 8688074827,22946.13,FILLED,22941.33,0.00065,15
时间:2023-02-09 10:37:49 订单:0.00065(22930.32000000) 	SELL挂单 下降:1240.51(5.410%) 亏损:0.80633u 	rescuelog: 8688141453,22930.32,FILLED,22925.52,0.00065,15
时间:2023-02-09 10:38:01 订单:0.00065(22928.69000000) 	SELL挂单 下降:1238.88(5.403%) 亏损:0.80527u 	rescuelog: 8688144170,22928.69,FILLED,22923.89,0.00065,15
时间:2023-02-09 10:40:04 订单:0.00065(22914.01000000) 	SELL挂单 下降:1224.2(5.343%) 亏损:0.79573u 	rescuelog: 8688174900,22914.01,FILLED,22909.21,0.00065,15
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
		var buyPrice, _ = decimal.NewFromFloat(p1.Tick.Bids[0][0]-4).DivRound(decimal.NewFromFloat(1), findSymbol.PricePrecision).Float64()
		var sellPrice, _ = decimal.NewFromFloat(p1.Tick.Asks[0][0]+4).DivRound(decimal.NewFromFloat(1), findSymbol.PricePrecision).Float64()
		btime := time.Now().Unix()
		_tryRescueOrder(p, findSymbol, group, groupCount, symbol, buyPrice, sellPrice, rescuelog)
		if time.Now().Unix()-btime > 5 {
			//刚救完，过几秒在开
			time.Sleep(time.Second * 3)
		}
	}
}

// 根据挂单批量救单
func Test_OnePlat_BTCBUSD_Debug_RescueOrder(t *testing.T) {
	onePlat := NewOnePlatService(plat.P_Binance)
	onePlat.timerOpenOrders("BTC/BUSD")
}

func Test_OnePlat_BTCBUSD_TryRescueOrder1(t *testing.T) {
	var symbol = "BTC/BUSD"

	p := plat.Get(plat.P_Binance)

	var findSymbol = NewOnePlatService(plat.P_Binance).getSymbol(symbol)

	p1, _ := p.GetMarketDepth(symbol)
	buyPrice, _ := decimal.NewFromFloat(p1.Tick.Bids[0][0]-2).DivRound(decimal.NewFromFloat(1), findSymbol.PricePrecision).Float64()
	sellPrice, _ := decimal.NewFromFloat(p1.Tick.Asks[0][0]+2).DivRound(decimal.NewFromFloat(1), findSymbol.PricePrecision).Float64()

	_tryRescueOrder(p, findSymbol, 1, 1, symbol, buyPrice, sellPrice, "17857436103,22992.81,FILLED,22992,0.00065,15")
}

func Test_OnePlat_BTCBUSD_TryRescueOrder2(t *testing.T) {
	var symbol = "BTC/BUSD"

	p := plat.Get(plat.P_Binance)

	var findSymbol = NewOnePlatService(plat.P_Binance).getSymbol(symbol)

	p1, _ := p.GetMarketDepth(symbol)
	buyPrice, _ := decimal.NewFromFloat(p1.Tick.Bids[0][0]-2).DivRound(decimal.NewFromFloat(1), findSymbol.PricePrecision).Float64()
	sellPrice, _ := decimal.NewFromFloat(p1.Tick.Asks[0][0]+2).DivRound(decimal.NewFromFloat(1), findSymbol.PricePrecision).Float64()

	_tryRescueOrder(p, findSymbol, 1, 1, symbol, buyPrice, sellPrice, "17857436103,22992.81,FILLED,22992,0.00065,15")
}

func Test_OnePlat_BTCBUSD_Buy(t *testing.T) {
	plat := plat.Get(plat.P_Binance)
	//（待救 FILLED,23296.61,18193062707,23290.22,0.09235,2150.851817)
	var buyprice = "23300"
	var btcamount = "0"
	orderid, err := plat.OrdersPlace("BTC/BUSD", buyprice, btcamount, "BUY")
	fmt.Println(orderid, err)
	time.Sleep(time.Second)
}

func Test_OnePlat_BTCBUSD_Sell(t *testing.T) {
	plat := plat.Get(plat.P_Binance)

	var sellprice = "23240"
	var btcamount = "0.02535"
	orderid, err := plat.OrdersPlace("BTC/BUSD", sellprice, btcamount, "SELL")
	//余额不足
	commerr, ok := err.(*common.APIError)
	if ok && commerr.Code == -2010 {
		fmt.Println("余额不足", err)
	} else {
		fmt.Println(orderid)
	}
	time.Sleep(time.Second)
}
