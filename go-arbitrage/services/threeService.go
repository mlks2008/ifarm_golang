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
	"components/message"
	"fmt"
	"github.com/shopspring/decimal"
	"goarbitrage/api"
	"goarbitrage/internal/model"
	"goarbitrage/internal/plat"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ThreeService struct {
	p              plat.Plat
	dayStartTime   int64
	dayStat        time.Time
	dayTotalCount  int64
	dayTotalProfit decimal.Decimal
	totalProfit    decimal.Decimal
	startTime      int64
	totalCount     int64
	timeoutCount   int64 //超时未完成次数
	//初始投入
	initInput map[string]model.Balance
	//启动时帐户余额
	run_bBalance map[string]model.Balance

	exchange_pending bool //usdt与busd兑换状态
}

func NewThreeService(platcode string) *ThreeService {

	ser := &ThreeService{
		p:            plat.Get(platcode),
		dayStartTime: time.Now().Unix(),
		dayStat:      time.Now(),
		startTime:    time.Now().Unix(),
	}
	ser.initInput = ser.p.GetInitialInput("BUSDUSDT")
	ser.run_bBalance, _ = ser.p.GetAccountBalance("BUSDUSDT", true)
	return ser
}

func (this *ThreeService) Start() {

	schedulingIndex, symbolsData := this.GetSchd()

	group := 1

	for {

		if api.OnePlatStop == true {
			message.PrintLog(api.PrintLog, fmt.Sprintf("every_stop:%v,%v", group, "BUSDUSDT"), fmt.Sprintf("%v", time.Now().Unix()/86400), fmt.Sprintf("机器人已暂停\n"))
			time.Sleep(time.Second * 5)
			continue
		}

		this.CccSched(schedulingIndex, symbolsData)

		group++

		time.Sleep(2 * time.Second)
	}
}

// 线上交易对可组成的三角套利
func (this *ThreeService) GetSchd() ([]model.SchedulingIndex, []model.SymbolsData) {
	//获取交易所所有交易对
	symbolsReturn, err := this.p.GetSymbols()
	if err != nil {
		panic(err)
	}

	scope, err := this.p.GetScope()
	if err != nil {
		panic(err)
	}

	//----仅搞BTCUSDT,BTCBUSD,BUSDUSDT
	if true {
		scope = map[string]string{"USDT": "", "BUSD": ""}

		var symbols = &model.SymbolsReturn{}
		symbols.Data = make([]model.SymbolsData, 0)
		for _, symbol := range symbolsReturn.Data {
			if symbol.Symbol == "BTCUSDT" || symbol.Symbol == "BTCBUSD" || symbol.Symbol == "BUSDUSDT" {
				symbols.Data = append(symbols.Data, symbol)
			}
		}
		symbolsReturn = symbols
	}

	//调度结构体，标记形成三角套利的每一步币种，在所有交易对数组的位置
	schedulingIndex := model.SchedulingIndex{}
	schedulingIndexArg := []model.SchedulingIndex{}
	//log.Println(symbolsReturn.Data[0].Symbol)
	for i := 0; i < len(symbolsReturn.Data); i++ {
		//log.Println("====",symbolsReturn.Data[i].State)
		//只要上线的交易对
		if symbolsReturn.Data[i].State {

			if _, ok := scope[symbolsReturn.Data[i].QuoteCurrency]; ok { //定套利母币，可根据需要调整
				//step_1:=symbolsReturn.Data[i].Symbol
				step_1_count := i
				//获取套利第二步交易对，与第一个交易对基础币相同，且计价币不同
				for j := 0; j < len(symbolsReturn.Data); j++ {

					if _, ok := scope[symbolsReturn.Data[j].QuoteCurrency]; ok { //定套利母币，可根据需要调整
						if symbolsReturn.Data[j].BaseCurrency == symbolsReturn.Data[step_1_count].BaseCurrency && j != step_1_count {
							//step_2:=symbolsReturn.Data[j].Symbol
							step_2_count := j
							//log.Println("=====",step_2,step_1)
							//获取套利第三步的交易对，取第一个交易对的计价币和第二个交易对的计价币，第三步交易对币种顺序无影响
							for k := 0; k < len(symbolsReturn.Data); k++ {
								if k != step_2_count {
									step_1_symbol := symbolsReturn.Data[step_1_count].QuoteCurrency + "" + symbolsReturn.Data[step_2_count].QuoteCurrency
									step_2_symbol := symbolsReturn.Data[step_2_count].QuoteCurrency + "" + symbolsReturn.Data[step_1_count].QuoteCurrency
									switch symbolsReturn.Data[k].Symbol {
									case step_1_symbol, step_2_symbol:
										//log.Println("k.....",k)
										//step_3:=symbolsReturn.Data[k].Symbol
										//log.Println("step_3",step_1,step_2,step_3,"step_3_count",step_1_count,step_2_count,k)
										schedulingIndex.Step1Index = step_1_count
										schedulingIndex.Step2Index = step_2_count
										schedulingIndex.Step3Index = k
										schedulingIndexArg = append(schedulingIndexArg, schedulingIndex)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return schedulingIndexArg, symbolsReturn.Data
}

func (this *ThreeService) CccSched(schedulingIndex []model.SchedulingIndex, symbolsData []model.SymbolsData) {

	for i := 0; i < len(schedulingIndex); i++ {
		time.Sleep(time.Millisecond * 2000)
		btime := time.Now()
		//log4go.Info(time.Now().Format("2006-01-02 15:04:05"))
		//获取123步交易对
		step_1_symbol := symbolsData[schedulingIndex[i].Step1Index].Symbol
		step_2_symbol := symbolsData[schedulingIndex[i].Step2Index].Symbol
		step_3_symbol := symbolsData[schedulingIndex[i].Step3Index].Symbol

		var isTrue bool
		if symbolsData[schedulingIndex[i].Step1Index].QuoteCurrency == symbolsData[schedulingIndex[i].Step3Index].QuoteCurrency {
			isTrue = true
		} else {
			isTrue = false
		}

		var group = 1
		var clientsymbol = "BUSDUSDT"
		var succ bool
		var usdtint float64
		var lotint float64
		var pricediff float64
		var sellPrice float64
		switch isTrue {
		//对比第一步和第三步交易计价币，例如：omgusdt omgbtc btcusdt  此时三角价格计算为sell buy buy （买 卖 卖)
		case true:
			//第一步深度价格和数量
			step_1_buy_data := this.PriceSched(step_1_symbol, "SELL")
			step_1_buy_price := step_1_buy_data[0]
			step_1_buy_count := step_1_buy_data[1]

			//第一步深度价格和数量
			step_2_sell_data := this.PriceSched(step_2_symbol, "BUY")
			step_2_sell_price := step_2_sell_data[0]
			step_2_sell_count := step_2_sell_data[1]

			//第一步深度价格和数量
			step_3_sell_data := this.PriceSched(step_3_symbol, "BUY")
			step_3_sell_price := 1.0 //step_3_sell_data[0]
			step_3_sell_count := step_3_sell_data[1]

			//套利空间计算
			v1_0 := decimal.NewFromFloat(step_1_buy_price)
			v2_0 := decimal.NewFromFloat(step_2_sell_price)
			v3_0 := decimal.NewFromFloat(step_3_sell_price)
			v2_add_v3 := v2_0.Mul(v3_0)
			v2_add_v3_sub_v1 := v2_add_v3.Sub(v1_0)
			rate := v2_add_v3_sub_v1.Div(v1_0)
			rate_float, _ := strconv.ParseFloat(rate.String(), 64)

			//下单策略
			lot := this.OrderCount(step_1_buy_count, step_2_sell_count, step_3_sell_count, step_2_sell_price)
			lotint, _ = strconv.ParseFloat(lot, 64)
			usdtint = step_1_buy_price * lotint
			if usdtint <= 10 {
				continue
			}
			//上限
			if usdtint > api.OnePlatUsdtAmount1 {
				usdtint = api.OnePlatUsdtAmount1
				lotint, _ = decimal.NewFromFloat(usdtint).DivRound(decimal.NewFromFloat(step_1_buy_price), 5).Float64()
			}
			pricediff, _ = decimal.NewFromFloat(step_2_sell_price).Sub(decimal.NewFromFloat(step_1_buy_price)).DivRound(decimal.NewFromFloat(1), 2).Float64()
			sellPrice = step_2_sell_price

			if math.Abs(pricediff) < 0.1 {
				continue
			}

			//log.Println(step_1_symbol,"=价格：=",step_1_buy_price,"=数量：=",step_1_buy_count,"=",step_2_symbol,"=价格：=",step_2_sell_price,"=数量：=",step_2_sell_count,"=",step_3_symbol,"=价格：=",step_3_sell_price,"=数量：=",step_2_sell_count,"套利空间：", rate_float,"=====可下单数量：",lot)

			if rate_float > 0 && lotint > 0 {
				log.Println("=============================================================套利空间:", rate_float)
				//log.Println(step_1_symbol, "step_1_buy_price：", step_1_buy_price)
				//log.Println(step_2_symbol, "step_2_sell_price：", step_2_sell_price)
				//log.Println(step_3_symbol, "step_3_sell_price：", step_3_sell_price)
				//log.Println("下单数量：", lotint)
				logmsg := fmt.Sprintf("group#%v B... %v diff:%vu(%v,%v) 运行挂单:%v", group, clientsymbol, pricediff, lotint, usdtint, 0)
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit:%v,%v", group, clientsymbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)

				var err error
				succ, err = this.arbitrageOrder(1, "BUSDUSDT", pricediff, step_1_symbol, step_1_buy_price, step_2_symbol, step_2_sell_price, lotint, usdtint)
				if err != nil {
					message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit00:%v,%v", "group", clientsymbol), fmt.Sprintf("%v", time.Now().UnixNano()), err.Error())
					continue
				}
				//this.p.OrdersPlace2("BUSDUSDT", step_1_symbol, fmt.Sprintf("%v", step_1_buy_price), lot, "BUY")
				//this.p.OrdersPlace2("BUSDUSDT", step_2_symbol, fmt.Sprintf("%v", step_2_sell_price), lot, "SELL")

				//this.provider.OrdersPlace(step_3_symbol, fmt.Sprintf("%v",step_3_sell_price), lot*(1-0.005)*(1-0.005), "SELL")
				log.Println("==============================================================")
			} else {
				continue
			}

		//对比第一步和第三步交易计价币，例如：omgeth omgbtc ethbtc  此时三角价格计算为sell buy sell (买 卖 买)
		default:

			//第一步深度价格和数量
			step_1_buy_data := this.PriceSched(step_1_symbol, "SELL")
			step_1_buy_price := step_1_buy_data[0]
			step_1_buy_count := step_1_buy_data[1]

			//第一步深度价格和数量
			step_2_sell_data := this.PriceSched(step_2_symbol, "BUY")
			step_2_sell_price := step_2_sell_data[0]
			step_2_sell_count := step_2_sell_data[1]

			//第一步深度价格和数量
			step_3_buy_data := this.PriceSched(step_3_symbol, "SELL")
			step_3_buy_price := 1.0 //step_3_buy_data[0]
			step_3_buy_count := step_3_buy_data[1]

			//套利空间计算
			v1_0 := decimal.NewFromFloat(step_1_buy_price)
			v2_0 := decimal.NewFromFloat(step_2_sell_price)
			v3_0 := decimal.NewFromFloat(step_3_buy_price)
			v2_add_v3 := v2_0.Div(v3_0)
			v2_add_v3_sub_v1 := v2_add_v3.Sub(v1_0)
			rate := v2_add_v3_sub_v1.Div(v1_0)
			rate_float, _ := strconv.ParseFloat(rate.String(), 64)

			lot := this.OrderCount(step_1_buy_count, step_2_sell_count, step_3_buy_count, step_2_sell_price)
			lotint, _ := strconv.ParseFloat(lot, 64)
			usdtint = step_1_buy_price * lotint
			if usdtint <= 10 {
				continue
			}
			//上限
			if usdtint > api.OnePlatUsdtAmount1 {
				usdtint = api.OnePlatUsdtAmount1
				lotint, _ = decimal.NewFromFloat(usdtint).DivRound(decimal.NewFromFloat(step_1_buy_price), 5).Float64()
			}

			pricediff, _ = decimal.NewFromFloat(step_2_sell_price).Sub(decimal.NewFromFloat(step_1_buy_price)).DivRound(decimal.NewFromFloat(1), 2).Float64()
			sellPrice = step_2_sell_price

			if math.Abs(pricediff) < 0.1 {
				continue
			}

			//下单策略
			//log.Println(step_1_symbol,"=价格：=",step_1_buy_price,"=数量：=",step_1_buy_count,"=",step_2_symbol,"=价格：=",step_2_sell_price,"=数量：=",step_2_sell_count,"=",step_3_symbol,"=价格：=",step_3_buy_price,"=数量：=",step_3_buy_count,"套利空间：", rate_float,"=====可下单数量：",lot)
			if rate_float > 0 && lotint > 0 {
				log.Println("=============================================================套利空间:", rate_float)
				//log.Println(step_1_symbol, "step_1_buy_price：", step_1_buy_price)
				//log.Println(step_2_symbol, "step_2_sell_price：", step_2_sell_price)
				//log.Println(step_3_symbol, "step_3_buy_price：", step_3_buy_price)
				//log.Println("下单数量：", lotint)
				logmsg := fmt.Sprintf("group#%v B... %v diff:%vu(%v,%v) 运行挂单:%v", group, clientsymbol, pricediff, lotint, usdtint, 0)
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit:%v,%v", group, clientsymbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)

				var err error
				succ, err = this.arbitrageOrder(1, "BUSDUSDT", pricediff, step_1_symbol, step_1_buy_price, step_2_symbol, step_2_sell_price, lotint, usdtint)
				if err != nil {
					message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit00:%v,%v", "group", clientsymbol), fmt.Sprintf("%v", time.Now().UnixNano()), err.Error())
					continue
				}
				//this.p.OrdersPlace2("BUSDUSDT", step_1_symbol, fmt.Sprintf("%v", step_1_buy_price), lot, "BUY")
				//this.p.OrdersPlace2("BUSDUSDT", step_2_symbol, fmt.Sprintf("%v", step_2_sell_price), lot, "SELL")
				log.Println("==============================================================")
			} else {
				continue
			}
		}

		//查询余额，计算收益
		for {
			bot_eBalance, err := this.p.GetAccountBalance(clientsymbol, false)
			if err != nil {
				continue
			}

			if succ == true {
				this.totalCount++
				var profit = decimal.NewFromFloat(pricediff).Mul(decimal.NewFromFloat(lotint))

				//当天
				if time.Now().Format("2006-01-02") != this.dayStat.Format("2006-01-02") {
					this.dayStartTime = time.Now().Unix() - 10
					this.dayStat = time.Now()
					this.dayTotalCount = 1
					this.dayTotalProfit = profit
				} else {
					this.dayTotalCount += 1
					this.dayTotalProfit = this.dayTotalProfit.Add(profit)
				}
				//总收益
				this.totalProfit = this.totalProfit.Add(profit)

				//运行累计收益
				runProfit, runChangeBase := this.profitFunc(sellPrice, this.run_bBalance, bot_eBalance, "USDT", "BUSD")
				//帐户累计收益
				accountProfit, accountChangeBase := this.profitFunc(sellPrice, this.initInput, bot_eBalance, "USDT", "BUSD")

				//运行时间
				runTotalUseTime := time.Now().Unix() - this.startTime
				runTotalFormatTime, _ := decimal.NewFromFloat(float64(runTotalUseTime)/86400.0).DivRound(decimal.NewFromInt(1), 2).Float64()
				runTotalFormatUnit := "d"
				if runTotalFormatTime < 1 {
					runTotalFormatTime, _ = decimal.NewFromFloat(float64(runTotalUseTime)/3600.0).DivRound(decimal.NewFromInt(1), 2).Float64()
					runTotalFormatUnit = "h"
				}

				//预估一天收益
				dayTotalUseTime := time.Now().Unix() - this.dayStartTime
				avgUseTime := dayTotalUseTime / this.dayTotalCount
				estimateDayCount := 86400 / avgUseTime
				estimateDayProfit := this.dayTotalProfit.DivRound(decimal.NewFromInt(this.dayTotalCount), 5).Mul(decimal.NewFromInt(estimateDayCount)).String()

				logmsg := fmt.Sprintf("group#%v E... %v(%v,%v),(%vu) diff:%vu   套利:%vu 运行累计:%vu(chanB:%v) 帐户累计:%vu(chanB:%v)   已运行:%v 挂单数:(%v,%v) {curprice}   预估Day:%vu/%vt usetime:%vs avgtime:%vs count:%v/%v\n",
					group, clientsymbol, bot_eBalance["USDT"].Free, bot_eBalance["BUSD"].Free, usdtint, pricediff, profit.String(), runProfit, runChangeBase, accountProfit.String(), accountChangeBase.String(),
					fmt.Sprintf("%v%v", runTotalFormatTime, runTotalFormatUnit), 0, 0,
					estimateDayProfit, estimateDayCount, int(time.Now().Sub(btime).Seconds()), avgUseTime, this.timeoutCount, this.totalCount)
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit:%v,%v", group, clientsymbol), fmt.Sprintf("%v", time.Now().UnixNano()), strings.Replace(logmsg, "{curprice}", "", 1))

				logmsg = strings.Replace(logmsg, "{curprice}", fmt.Sprintf("币价:%v", sellPrice), 1)
				logmsg = strings.Replace(logmsg, "套利", "\n套利", 1)
				logmsg = strings.Replace(logmsg, "已运行", "\n\n已运行", 1)
				logmsg = strings.Replace(logmsg, "预估Day", "\n预估Day", 1)
				message.SendDingTalkRobit(true, "oneplat", "every_profit"+clientsymbol, fmt.Sprintf("%v", time.Now().Unix()/(10*60)), logmsg)
				//日报
				message.SendDingTalkRobit(true, "report", "report_oneplat"+clientsymbol, fmt.Sprintf("%v", time.Now().Unix()/(24*60*60)), logmsg)
			} else {
				this.timeoutCount++
				logmsg := fmt.Sprintf("group#%v E... %v(%v,%v),(%vu) \t diff:%vu 挂单:(%v,%v) \t  \t usetime:%vs \t timeout ---------\n", group, clientsymbol, bot_eBalance["USDT"].Free, bot_eBalance["BUSD"].Free, usdtint, pricediff, 0, 0, int(time.Now().Sub(btime).Seconds()))
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit:%v,%v", group, clientsymbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
			}

			break
		}
	}
	//price_count:=PriceSched("","")
}

// 套利下单
func (this *ThreeService) arbitrageOrder(group int, clientsymbol string, pricediff float64, buysymbol string, buyprice float64, sellsymbol string, sellprice float64, btcamount, usdtamount float64) (succ bool, err error) {
	var sellorderid string
	var buyorderid string

	var btime = time.Now() //下单开始时间
	var waitGroup sync.WaitGroup

	makerfee, takerfee, err := this.p.GetTradeFee(buysymbol)
	if err != nil {
		return false, fmt.Errorf("\t\t%v:group#%v 下单失败:GetTradeFee------%v \n", clientsymbol, group, err.Error())
	}
	if makerfee.Cmp(decimal.NewFromInt(0)) > 0 || takerfee.Cmp(decimal.NewFromInt(0)) > 0 {
		return false, fmt.Errorf("\t\t%v:group#%v 下单失败:GetTradeFee不为0(%v,%v) \n", clientsymbol, group, makerfee.String(), takerfee.String())
	}

	makerfee, takerfee, err = this.p.GetTradeFee(sellsymbol)
	if err != nil {
		return false, fmt.Errorf("\t\t%v:group#%v 下单失败:GetTradeFee------%v \n", clientsymbol, group, err.Error())
	}
	if makerfee.Cmp(decimal.NewFromInt(0)) > 0 || takerfee.Cmp(decimal.NewFromInt(0)) > 0 {
		return false, fmt.Errorf("\t\t%v:group#%v 下单失败:GetTradeFee不为0(%v,%v) \n", clientsymbol, group, makerfee.String(), takerfee.String())
	}

	var insufficientBalance bool = true
	var balance map[string]model.Balance
	for i := 0; i < 1; i++ {
		var err error
		balance, err = this.p.GetAccountBalance(clientsymbol, false)
		if err != nil {
			fmt.Println("-----getbalance--------------------------------", err, time.Now().Format(time.Stamp))
		} else {
			if balance[strings.Replace(buysymbol, "BTC", "", 1)].Free.Cmp(decimal.NewFromFloat(usdtamount)) > 0 {
				insufficientBalance = false
			}
			break
		}
	}

	//余额不足
	if insufficientBalance == true {
		if balance == nil {
			return false, fmt.Errorf("\t\t%v:group#%v 下单失败:获取余额接口出错,无法完成下单 \n", clientsymbol, group)
		}

		logmsg := fmt.Sprintf("%v:group:%v 余额不足无法下单", clientsymbol, group)
		message.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_order2:%v", clientsymbol), fmt.Sprintf("%v", time.Now().Unix()/(60*60)), logmsg)
		message.PrintLog(api.PrintLog, fmt.Sprintf("every_order2:%v,%v", group, clientsymbol), fmt.Sprintf("%v", time.Now().Unix()/(60*60)), logmsg)

		//余额不足,刷新缓存余额
		this.p.GetAccountBalance(clientsymbol, true)

		if strings.Replace(buysymbol, "BTC", "", 1) == "USDT" {
			orderid, err := this.p.OrdersPlace2(clientsymbol, clientsymbol, "1", balance["BUSD"].Free.DivRound(decimal.NewFromFloat(2), 0).String(), "SELL")
			if err != nil {
				return false, err
			}
			for {
				ok, _ := this.p.GetOrderStatus2(clientsymbol, clientsymbol, orderid)
				if ok == true {
					break
				} else {
					time.Sleep(time.Second * 60)
				}
			}
		} else if strings.Replace(buysymbol, "BTC", "", 1) == "BUSD" {
			orderid, err := this.p.OrdersPlace2(clientsymbol, clientsymbol, "1", balance["USDT"].Free.DivRound(decimal.NewFromFloat(2), 0).String(), "BUY")
			if err != nil {
				return false, err
			}
			for {
				ok, _ := this.p.GetOrderStatus2(clientsymbol, clientsymbol, orderid)
				if ok == true {
					break
				} else {
					time.Sleep(time.Second * 60)
				}
			}
		}
		return false, fmt.Errorf("\t\t%v:group#%v 下单失败:余额不足无法下单 \n", clientsymbol, group)
	}

	//busd可用余额大于usdt可用余额2倍时，卖一半busd
	if balance["BUSD"].Free.Div(balance["USDT"].Free).Cmp(decimal.NewFromInt(2)) >= 0 && this.exchange_pending == false {
		this.exchange_pending = true
		go func() {
			orderid, err := this.p.OrdersPlace2(clientsymbol, clientsymbol, "1", balance["BUSD"].Free.DivRound(decimal.NewFromFloat(2), 0).String(), "SELL")
			if err != nil {
				this.exchange_pending = false
				return
			}
			for {
				ok, _ := this.p.GetOrderStatus2(clientsymbol, clientsymbol, orderid)
				if ok == true {
					this.exchange_pending = false
					break
				} else {
					time.Sleep(time.Second * 60)
				}
			}
		}()
	}

	var stoptime = int64(12 * 60)

	//下卖单
	waitGroup.Add(1)
	go func() {
		var err error
		for {
			sellorderid, err = this.p.OrdersPlace2(clientsymbol, sellsymbol, fmt.Sprintf("%v", sellprice), fmt.Sprintf("%v", btcamount), "SELL")
			if err == nil {
				break
			}

			//随机停止时间
			err_min := 1
			err_max := 3
			randSleepTime := time.Duration(rand.Intn(err_max-err_min) + err_min)
			time.Sleep(time.Second * randSleepTime)
		}
		if sellorderid != "FILLED" {
			for {
				var usetime = time.Now().Unix() - btime.Unix()
				//超时未完成，退出查询
				if usetime > stoptime {
					sellorderid = "TIMEOUT" + sellorderid
					time.Sleep(time.Millisecond * 200)
					waitGroup.Done()
					break
				}
				//查询订单
				orderstatus, err := this.p.GetOrderStatus2(clientsymbol, sellsymbol, sellorderid)
				if err != nil {
					//随机停止时间
					err_min := 5
					err_max := 10
					randSleepTime := time.Duration(rand.Intn(err_max-err_min) + err_min)
					time.Sleep(time.Second * randSleepTime)
					message.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_query_err,%v", clientsymbol), fmt.Sprintf("%v", time.Now().Unix()/(12*60*60)), err.Error())
					continue
				}
				if orderstatus == true && err == nil {
					sellorderid = "FILLED"
					time.Sleep(time.Millisecond * 200)
					waitGroup.Done()
					break
				}
				time.Sleep(time.Millisecond * 2000)
			}
		}
	}()

	//下买单
	waitGroup.Add(1)
	go func() {
		//defer utils.PrintPanicStack()
		var err error
		for {
			buyorderid, err = this.p.OrdersPlace2(clientsymbol, buysymbol, fmt.Sprintf("%v", buyprice), fmt.Sprintf("%v", btcamount), "BUY")
			if err == nil {
				break
			}

			//随机停止时间
			err_min := 1
			err_max := 3
			randSleepTime := time.Duration(rand.Intn(err_max-err_min) + err_min)
			time.Sleep(time.Second * randSleepTime)
		}
		if buyorderid != "FILLED" {
			for {
				var usetime = time.Now().Unix() - btime.Unix()
				//超时未完成，退出查询
				if usetime > stoptime {
					buyorderid = "TIMEOUT" + buyorderid
					time.Sleep(time.Millisecond * 200)
					waitGroup.Done()
					break
				}
				//查询订单
				orderstatus, err := this.p.GetOrderStatus2(clientsymbol, buysymbol, buyorderid)
				if err != nil {
					//随机停止时间
					err_min := 5
					err_max := 10
					randSleepTime := time.Duration(rand.Intn(err_max-err_min) + err_min)
					time.Sleep(time.Second * randSleepTime)
					message.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_query_err,%v", clientsymbol), fmt.Sprintf("%v", time.Now().Unix()/(12*60*60)), err.Error())
					continue
				}
				if orderstatus == true && err == nil {
					buyorderid = "FILLED"
					time.Sleep(time.Millisecond * 200)
					waitGroup.Done()
					break
				}
				time.Sleep(time.Millisecond * 2000)
			}
		}
	}()

	//超时，退出等待，进入下一组
	go func(sellorderId *string, buyorderId *string) {
		//defer utils.PrintPanicStack()
		for {
			var usetime = time.Now().Unix() - btime.Unix()
			if usetime > stoptime {
				logmsg := fmt.Sprintf("\t\t%v:group:%v 套利失败 超时:%v(s) rescuelog: %v,%v,%v,%v,%v,%v", clientsymbol, group, usetime, *sellorderId, sellprice, *buyorderId, buyprice, btcamount, usdtamount)
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_order6:%v,%v", group, clientsymbol), fmt.Sprintf("%v", time.Now().Unix()), logmsg)
				return
			} else {
				if *sellorderId == "FILLED" && *buyorderId == "FILLED" {
					return
				} else {
					//进行中
					logmsg := fmt.Sprintf("\t\t%v:group:%v .... 用时:%v(s) rescuelog: %v,%v,%v,%v,%v,%v  createtime:%v", clientsymbol, group, usetime, *sellorderId, sellprice, *buyorderId, buyprice, btcamount, usdtamount, btime.Format("2006-01-02 15:04:05"))
					message.PrintLog(api.PrintLog, fmt.Sprintf("every_order7:%v,%v", group, clientsymbol), fmt.Sprintf("%v", time.Now().Unix()/5), logmsg)
					time.Sleep(time.Millisecond * 1000)
				}
			}
		}
	}(&sellorderid, &buyorderid)

	waitGroup.Wait()

	if buyorderid == "FILLED" && sellorderid == "FILLED" {
		return true, nil
	} else {
		return false, nil
	}
}

func (this *ThreeService) profitFunc(sellPrice float64, iBalance map[string]model.Balance, eBalance map[string]model.Balance, fiat1, fiat2 string) (decimal.Decimal, decimal.Decimal) {
	iTotal := iBalance[fiat1].Total.Add(iBalance[fiat2].Total)

	//多出币换算为等值U
	changeBaseToQuote := eBalance["BTC"].Total.Sub(iBalance["BTC"].Total).Mul(decimal.NewFromFloat(sellPrice))

	//当前最新U余额
	eTotal := eBalance[fiat1].Total.Add(eBalance[fiat2].Total).Add(changeBaseToQuote)

	//收益
	profit := eTotal.Sub(iTotal).DivRound(decimal.NewFromFloat(1), 5)

	return profit, eBalance["BTC"].Total
}

func (this *ThreeService) PriceSched(step_symbol string, method string) []float64 {
	var err error
	var marketDepthReturn model.MarketDepthReturn
	for i := 0; i < 5; i++ {
		marketDepthReturn, err = this.p.GetMarketDepth(step_symbol)
		if err != nil {
			//log4go.Error("%v", err)
			time.Sleep(time.Second * 5)
			continue
		} else {
			break
		}
	}

	switch method {
	case "SELL":
		return marketDepthReturn.Tick.Asks[0]
	case "BUY":
		return marketDepthReturn.Tick.Bids[0]
	default:
		return []float64{-1, -1}
	}
}

func (this *ThreeService) OrderCount(lot1, lot2, lot3, price2 float64) string {
	if lot1 < 0 || lot2 < 0 || lot3 < 0 {
		return ""
	}

	lot1_0 := decimal.NewFromFloat(lot1)
	lot2_0 := decimal.NewFromFloat(lot2)
	lot3_0 := decimal.NewFromFloat(lot3)
	price2_0 := decimal.NewFromFloat(price2)

	if lot1 < lot2 {
		lot_x := lot1_0.Mul(price2_0)
		lot_x_float, _ := strconv.ParseFloat(lot_x.String(), 64)
		if lot_x_float < lot3 {
			return lot1_0.String()
			//log.Println("此次套利下单量为", lot1_0.String())
		} else {
			lot_2 := lot3_0.Div(price2_0)
			return lot_2.String()
			//log.Println("此次套利下单量为", lot_2.String())
		}
	} else {
		lot_x := lot2_0.Mul(price2_0)
		lot_x_float, _ := strconv.ParseFloat(lot_x.String(), 64)
		if lot_x_float < lot3 {
			return lot2_0.String()
			//log.Println("此次套利下单量为", lot2_0.String())
		} else {
			lot_2 := lot3_0.Div(price2_0)
			return lot_2.String()
			//log.Println("此次套利下单量为", lot_2.String())
		}
	}
	return ""
}
