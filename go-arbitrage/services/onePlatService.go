/**
 * onePlatService.go
 * ============================================================================
 * 交易所内买一卖一快速交易策略
 * ============================================================================
 * author: peter.wang
 */

package services

import (
	"components/message"
	"errors"
	"fmt"
	"github.com/adshao/go-binance/v2/common"
	"github.com/shopspring/decimal"
	"goarbitrage/api"
	"goarbitrage/internal/model"
	"goarbitrage/internal/plat"
	"goarbitrage/pkg/utils"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type OnePlatService struct {
	p              plat.Plat
	orderPlaceLock sync.Mutex //声明一个全局互斥锁

	apiopenordercount   int //当前挂单数
	findSymbol          model.SymbolsData
	makerFee            float64 //挂单手续费
	takerFee            float64 //吃单手续费
	firstRescueBeginSec int     //n秒后未完成套利开启首次救单
	timerRescueSucc     int64   //累计救单成功数

	orderPrice          sync.Map //预防自成交(key:group value:model.OrderInfo)
	rescueMap           sync.Map //救单机器人列表(包含运行和未开启)
	workingRescueRots   int      //已开启救单机器人数
	manual_price_enable bool     //开启人工设买卖价救单
	manual_curbuyprice  float64  //人工救单买价
	manual_cursellprice float64  //人工救单卖价
}

func NewOnePlatService(platcode string) *OnePlatService {
	ser := &OnePlatService{
		p:                   plat.Get(platcode),
		firstRescueBeginSec: 5 * 60, //默认5分钟后未完成套利开启首次救单
	}
	return ser
}

func (this *OnePlatService) getSymbol(symbol string) model.SymbolsData {
	platSymbols, err := this.p.GetSymbols()
	if err != nil {
		panic(err)
	}
	//交易对信息
	var fdSymbol model.SymbolsData
	var find bool
	for _, tmp := range platSymbols.Data {
		if tmp.Symbol == this.p.FormatSymbol(symbol) {
			fdSymbol = tmp
			find = true
			break
		}
	}
	if find == true {
		return fdSymbol
	} else {
		panic(fmt.Sprintf("not find symbol %v", symbol))
	}
}

func (this *OnePlatService) profitFunc(sellPrice float64, iBalance map[string]model.Balance, eBalance map[string]model.Balance, base, quote string) (decimal.Decimal, decimal.Decimal) {
	//多出币换算为等值U
	changeBase := eBalance[base].Total.Sub(iBalance[base].Total)
	changeBaseToQuote := changeBase.Mul(decimal.NewFromFloat(sellPrice))
	//当前最新U余额
	nowQuoteBalance := eBalance[quote].Total.Add(changeBaseToQuote)
	//收益
	profit := nowQuoteBalance.Sub(iBalance[quote].Total).DivRound(decimal.NewFromFloat(1), 4)
	return profit, changeBase
}

func (this *OnePlatService) Start(symbol string, botcount int, newchildsec int64, firstrescuebeginsec int) {
	//并发机器人数
	var symbols = make([]string, botcount)
	for i := 0; i < botcount; i++ {
		symbols[i] = symbol
	}

	this.firstRescueBeginSec = firstrescuebeginsec
	this.findSymbol = this.getSymbol(symbol)

	//初始投入
	initInput := this.p.GetInitialInput(symbol)

	//启动时帐户余额
	run_bBalance, err := this.p.GetAccountBalance(symbol, true)
	if err != nil {
		panic(err)
	}

	var stopMainBots = make([]string, 0)
	var childBots sync.Map
	var timeoutCount int64          //超时未完成次数
	var totalCount int64            //总套利次数
	var totalProfit decimal.Decimal //总套利额
	var startTime = time.Now().Unix()
	//var continuousLossCount int64   //连续亏损次数

	var dayStartTime = time.Now().Unix()
	var dayStat = time.Now()
	var dayTotalCount int64
	var dayTotalProfit decimal.Decimal

	var arbitrage_bot = func(group int, symbol string, mainbot bool) {
		for {
			defer utils.PrintPanicStack()

			//已暂停交易
			if api.OnePlatStop == true {
				time.Sleep(time.Second * 5)
				//utils.PrintLog(api.PrintLog, fmt.Sprintf("every_stop_0:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Day()), fmt.Sprintf("%v号机器人暂停\n", group))

				if group < 1000 {
					var exist bool
					for _, g := range stopMainBots {
						if g == fmt.Sprintf("%v", group) {
							exist = true
							break
						}
					}
					if exist == false {
						stopMainBots = append(stopMainBots, fmt.Sprintf("%v", group))
						//stopbotStr := strings.Join(stopMainBots, ",")

						existRescueRot := this.existRescueRot(0, "")
						if existRescueRot == true {
							message.PrintLog(api.PrintLog, fmt.Sprintf("every_stop:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/86400), fmt.Sprintf("%v共%v个,已有%v个机器人暂停,还有救单机器人在运行\n", symbol, len(symbols), len(stopMainBots)))
						} else {
							message.PrintLog(api.PrintLog, fmt.Sprintf("every_stop:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/86400), fmt.Sprintf("%v共%v个,已有%v个机器人暂停\n", symbol, len(symbols), len(stopMainBots)))
						}
					}
				} else {
					existRescueRot := this.existRescueRot(0, "")
					if existRescueRot == true {
						message.PrintLog(api.PrintLog, fmt.Sprintf("every_stop:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/86400), fmt.Sprintf("%v共%v个,已有%v个机器人暂停,还有救单机器人在运行\n", symbol, len(symbols), len(stopMainBots)))
					} else {
						message.PrintLog(api.PrintLog, fmt.Sprintf("every_stop:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/86400), fmt.Sprintf("%v共%v个,已有%v个机器人暂停\n", symbol, len(symbols), len(stopMainBots)))
					}
				}
				continue
			} else {
				stopMainBots = make([]string, 0)
			}

			//主机器人恢复后退出子机器人
			if mainbot == false {
				parantid := group - 1000
				if _, ok := this.rescueMap.Load(fmt.Sprintf("%v_%v_%v", parantid, symbol, 1)); ok == false {
					logmsg := fmt.Sprintf("%v:group:%v, %v号机器人恢复,我已退出\n", symbol, group, parantid)
					message.PrintLog(api.PrintLog, fmt.Sprintf("every_child:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/10), logmsg)
					//删除子机器人启动标记
					childBots.Delete(group)
					time.Sleep(time.Second * 5)
					return
				}
			}

			//随机停止时间
			err_min := 5
			err_max := 15
			randSleepTime := time.Duration(rand.Intn(err_max-err_min) + err_min)

			//交易手续费判断
			makerfee, takerfee, err := this.p.GetTradeFee(symbol)
			if err != nil {
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_error:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/10), err.Error())
				time.Sleep(time.Millisecond * 1000)
				continue
			}
			this.makerFee, _ = makerfee.Float64()
			this.takerFee, _ = takerfee.Float64()
			//if makerfee.Cmp(decimal.NewFromInt(0)) > 0 || takerfee.Cmp(decimal.NewFromInt(0)) > 0 {
			//	logmsg := fmt.Sprintf("\t\tgroup#%v, 交易手续费不为0，不进行套利(makerfee:%v,takerfee:%v) \n", group, makerfee, takerfee)
			//	utils.PrintLog(fmt.Sprintf("every_error:%v", group), fmt.Sprintf("%v", time.Now().Unix()/60), logmsg)
			//	utils.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_stop:%v", group), fmt.Sprintf("%v", time.Now().Unix()/3600), logmsg)
			//	time.Sleep(time.Second * 60)
			//	continue
			//}

			symbolBase := strings.Split(symbol, "/")[0]
			symbolQuote := strings.Split(symbol, "/")[1]

			//计算本次用时
			btime := time.Now()
			////计算本次收益
			//bot_bBalance, err := this.p.GetAccountBalance()
			//if err != nil {
			//	utils.PrintLog(fmt.Sprintf("every_error:%v", group), fmt.Sprintf("%v", time.Now().Unix()/10), err.Error())
			//	time.Sleep(time.Millisecond * 1000)
			//	continue
			//}

			p1, err := this.p.GetMarketDepth(symbol)
			if err != nil {
				logmsg := fmt.Sprintf("\t\t%vgroup#%v, -----err----- %v \n", symbol, group, err.Error())
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()), logmsg)
				time.Sleep(time.Second * randSleepTime)
				continue
			}
			if p1.Tick.Bids[0][0] == 0 || len(p1.Tick.Bids) < 1 {
				time.Sleep(time.Millisecond * 1000)
				continue
			}
			if p1.Tick.Asks[0][0] == 0 || len(p1.Tick.Asks) < 1 {
				time.Sleep(time.Millisecond * 1000)
				continue
			}

			//买一卖一对应报价币(usdt)数量
			usdtBid := decimal.NewFromFloat(p1.Tick.Bids[0][0]).Mul(decimal.NewFromFloat(p1.Tick.Bids[0][1])).DivRound(decimal.NewFromInt(1), 0)
			usdtAsks := decimal.NewFromFloat(p1.Tick.Asks[0][0]).Mul(decimal.NewFromFloat(p1.Tick.Asks[0][1])).DivRound(decimal.NewFromInt(1), 0)
			//取报价币最小数量，一次性成交
			var usdtAmount float64
			if usdtBid.Cmp(usdtAsks) < 0 {
				usdtAmount, _ = usdtBid.Float64()
			} else {
				usdtAmount, _ = usdtAsks.Float64()
			}

			var buyPrice, _ = decimal.NewFromFloat(p1.Tick.Bids[0][0]-0).DivRound(decimal.NewFromFloat(1), this.findSymbol.PricePrecision).Float64()
			var sellPrice, _ = decimal.NewFromFloat(p1.Tick.Asks[0][0]+0).DivRound(decimal.NewFromFloat(1), this.findSymbol.PricePrecision).Float64()

			//价差
			pricediff, _ := decimal.NewFromFloat(sellPrice-buyPrice).DivRound(decimal.NewFromFloat(1), this.findSymbol.PricePrecision).Float64()
			//最小价差:makerfee+takerfee+万二市价
			minpricediff, _ := decimal.NewFromFloat(sellPrice).Mul(makerfee.Add(takerfee).Add(decimal.NewFromFloat(1/10000.0))).DivRound(decimal.NewFromFloat(1), this.findSymbol.PricePrecision).Float64()
			if pricediff < minpricediff {
				half := (minpricediff - pricediff) / 2
				buyPrice, _ = decimal.NewFromFloat(buyPrice-half).DivRound(decimal.NewFromFloat(1), this.findSymbol.PricePrecision).Float64()   //降低买价
				sellPrice, _ = decimal.NewFromFloat(sellPrice+half).DivRound(decimal.NewFromFloat(1), this.findSymbol.PricePrecision).Float64() //提高卖价
				//重新算价差
				pricediff, _ = decimal.NewFromFloat(sellPrice-buyPrice).DivRound(decimal.NewFromFloat(1), this.findSymbol.PricePrecision).Float64()
			}

			//报价币(usdt)数量(限定最大值)
			if usdtAmount > api.OnePlatUsdtAmount1 {
				usdtAmount = api.OnePlatUsdtAmount1
			}
			//计算基础币(btc)下单量
			var btcAmount, _ = decimal.NewFromFloat(usdtAmount).DivRound(decimal.NewFromFloat(buyPrice), this.findSymbol.AmountPrecision).Float64()

			//--- 价差大可能波动比较大,不适合当前策略
			if pricediff < 5 {
				//todo
			} else {
				logmsg := fmt.Sprintf("\t\tgroup#%v, ------pricediff too big------%v \n", group, pricediff)
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/10), logmsg)
				time.Sleep(200 * time.Millisecond)
				continue
			}

			//--- 小于最小下单量
			if usdtAmount <= this.findSymbol.MinQuoteOrder {
				//logmsg := fmt.Sprintf("\t\tgroup#%v, ------minQuoteOrder------%v \n", group, usdtAmount)
				//utils.PrintLog(fmt.Sprintf("every_profit:%v", group), fmt.Sprintf("%v", time.Now().Unix()/10), logmsg)
				time.Sleep(200 * time.Millisecond)
				continue
			}

			//todo debuglog
			logmsg := fmt.Sprintf("group#%v B... %v go", group, symbol)
			message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit_start:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
			////todo debuglog
			//logmsg1 := fmt.Sprintf("group#%v B... %v/%v(%v,%v),(%v,%v) diff:%vu(%v,%v)", group, symbolBase, symbolQuote, "*", "*", btcAmount, usdtAmount, pricediff, sellPrice, buyPrice)
			//utils.PrintLog(api.PrintLog, fmt.Sprintf("every_profit:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg1)

			timeout, err := this.arbitrageOrder(0, group, symbol, this.findSymbol, pricediff, sellPrice, buyPrice, btcAmount, usdtAmount)
			if err != nil {
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit00:%v,%v", "group", symbol), fmt.Sprintf("%v", time.Now().UnixNano()), err.Error())
				time.Sleep(time.Second)
				continue
			}

			//查询余额，计算收益
			for {
				bot_eBalance, err := this.p.GetAccountBalance(symbol, false)
				if err != nil {
					continue
				}

				//挂单数
				var runopenorders = 0
				this.orderPrice.Range(func(key, value any) bool {
					runopenorders++
					return true
				})

				if timeout == false {
					totalCount++
					var profit = decimal.NewFromFloat(pricediff).Mul(decimal.NewFromFloat(btcAmount))

					//当天
					if time.Now().Format("2006-01-02") != dayStat.Format("2006-01-02") {
						dayStartTime = time.Now().Unix() - 10
						dayStat = time.Now()
						dayTotalCount = 1
						dayTotalProfit = profit
					} else {
						dayTotalCount += 1
						dayTotalProfit = dayTotalProfit.Add(profit)
					}
					//总收益
					totalProfit = totalProfit.Add(profit)

					//运行累计收益
					runProfit, runChangeBase := this.profitFunc(sellPrice, run_bBalance, bot_eBalance, symbolBase, symbolQuote)
					//帐户累计收益
					accountProfit, accountChangeBase := this.profitFunc(sellPrice, initInput, bot_eBalance, symbolBase, symbolQuote)

					//运行时间
					runTotalUseTime := time.Now().Unix() - startTime
					runTotalFormatTime, _ := decimal.NewFromFloat(float64(runTotalUseTime)/86400.0).DivRound(decimal.NewFromInt(1), 2).Float64()
					runTotalFormatUnit := "d"
					if runTotalFormatTime < 1 {
						runTotalFormatTime, _ = decimal.NewFromFloat(float64(runTotalUseTime)/3600.0).DivRound(decimal.NewFromInt(1), 2).Float64()
						runTotalFormatUnit = "h"
					}

					//预估一天收益
					dayTotalUseTime := time.Now().Unix() - dayStartTime
					avgUseTime := dayTotalUseTime / dayTotalCount
					estimateDayCount := 86400 / avgUseTime
					estimateDayProfit := dayTotalProfit.DivRound(decimal.NewFromInt(dayTotalCount), 4).Mul(decimal.NewFromInt(estimateDayCount)).String()

					////计算本次收益(多协程原因，有时不准）
					//eValue := bot_eBalance[symbolQuote].Total
					//bValue := bot_bBalance[symbolQuote].Total
					//profit := eValue.Sub(bValue).DivRound(decimal.NewFromFloat(1), 5)
					////统计连续亏损次数
					//if profit.Cmp(decimal.NewFromFloat(0)) >= 0 {
					//	//盈
					//	continuousLossCount = 0
					//} else {
					//	//亏
					//	continuousLossCount++
					//}
					////连续亏损说明存在问题
					//if continuousLossCount >= int64(5*len(symbols)) {
					//	logmsg := fmt.Sprintf("group#%v STOP... %v:%v,%v:%v continuousLossCount:%v count:%v/%v", group, symbolBase, bot_eBalance[symbolBase].Total, symbolQuote, bot_eBalance[symbolQuote].Total, continuousLossCount, timeoutCount, totalCount)
					//	utils.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_stop:%v", group), fmt.Sprintf("%v", time.Now().Unix()), logmsg)
					//	utils.PrintLog(fmt.Sprintf("every_profit:%v", group), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
					//	api.OnePlatStop = true
					//	time.Sleep(time.Hour * 24 * 365)
					//}

					logmsg := fmt.Sprintf("group#%v E... %v/%v(%v,%v),(%v,%v) diff:%vu   套利:%vu 运行累计:%vu(chanB:%v) 帐户累计:%vu(chanB:%v)   已运行:%v 挂单数:(%v,%v) {curprice}   预估Day:%vu/%vt usetime:%vs avgtime:%vs count:%v/%v/%v\n",
						group, symbolBase, symbolQuote, bot_eBalance[symbolBase].Free, utils.DivRound(bot_eBalance[symbolQuote].Free, 2), btcAmount, usdtAmount, pricediff, utils.DivRound(profit, 4), runProfit, runChangeBase, accountProfit.String(), accountChangeBase.String(),
						fmt.Sprintf("%v%v", runTotalFormatTime, runTotalFormatUnit), runopenorders, this.apiopenordercount,
						estimateDayProfit, estimateDayCount, int(time.Now().Sub(btime).Seconds()), avgUseTime, this.timerRescueSucc, timeoutCount, totalCount)
					message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), strings.Replace(logmsg, "{curprice}", "", 1))

					logmsg = strings.Replace(logmsg, "{curprice}", fmt.Sprintf("币价:%v", sellPrice), 1)
					logmsg = strings.Replace(logmsg, "套利", "\n套利", 1)
					logmsg = strings.Replace(logmsg, "已运行", "\n\n已运行", 1)
					logmsg = strings.Replace(logmsg, "预估Day", "\n预估Day", 1)
					message.SendDingTalkRobit(true, "oneplat", "every_profit"+symbol, fmt.Sprintf("%v", time.Now().Unix()/(10*60)), logmsg)
					//日报
					message.SendDingTalkRobit(true, "report", "report_oneplat"+symbol, fmt.Sprintf("%v", time.Now().Unix()/(24*60*60)), logmsg)
				} else {
					timeoutCount++
					pricediff := decimal.NewFromFloat(sellPrice-buyPrice).DivRound(decimal.NewFromFloat(1), 8).String()
					logmsg := fmt.Sprintf("group#%v E... %v/%v(%v,%v),(%vu) \t diff:%vu 挂单:(%v,%v) \t buyprice:%v sellprice:%v \t usetime:%vs \t timeout ---------\n",
						group, symbolBase, symbolQuote, bot_eBalance[symbolBase].Free, utils.DivRound(bot_eBalance[symbolQuote].Free, 2), usdtAmount, pricediff, runopenorders, this.apiopenordercount,
						buyPrice, sellPrice, int(time.Now().Sub(btime).Seconds()))
					//utils.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_profit:%v", group), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
					message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
				}

				//下一单随机间隔n秒
				time.Sleep(time.Second * randSleepTime)
				break
			}
		}
	}

	//启动主机器人
	go func() {
		for r_index, r_symbol := range symbols {
			r_index++
			go arbitrage_bot(r_index, r_symbol, true)
			time.Sleep(time.Second * 5)
		}
	}()

	//启动子机器人
	go func() {
		for {
			if api.OnePlatStop == false {
				this.orderPrice.Range(func(key, value any) bool {
					order := value.(*OrderInfo)
					btime := order.CreateTime
					group, symbol, rescuecount := this.parseRotId(key.(string))

					childid := group + 1000
					if _, ok := childBots.Load(childid); !ok {
						//主机器人超n分钟未完成，新启一个子机器人
						if rescuecount == 0 && time.Now().Unix()-btime > newchildsec {
							logmsg := fmt.Sprintf("%v:group:%v, %v号机器人超时未完成,我已开启\n", symbol, childid, group)
							message.PrintLog(api.PrintLog, fmt.Sprintf("every_child:%v,%v", childid, symbol), fmt.Sprintf("%v", time.Now().Unix()/10), logmsg)
							//添加子机器人启动标记
							childBots.Store(childid, "1")
							go arbitrage_bot(childid, symbol, false)
							time.Sleep(time.Second * 10)
						}
					}
					return true
				})
			}
			time.Sleep(time.Second * 5)
		}
	}()

	//启动拯救机器人
	go func() {
		this.timerOpenOrders(symbol)
	}()
}

// 检测市价是否在平均价周围
func (this *OnePlatService) CheckPrice(symbol string) bool {
	p1, err := this.p.GetMarketDepth(symbol)
	if err != nil {
		time.Sleep(time.Second * 5)
		return false
	}
	curprice := decimal.NewFromFloat(p1.Tick.Bids[0][0])

	avgprice, err := this.p.AveragePrice(symbol)
	if err != nil {
		time.Sleep(time.Second * 5)
		return false
	}

	//todo 波动中，暂时不救
	diff, _ := curprice.Sub(avgprice).Float64()
	if math.Abs(diff) >= 3 {
		time.Sleep(time.Second * 60)
		return false
	}
	return true
}

// 拯救挂单列表
func (this *OnePlatService) timerOpenOrders(symbol string) {
	this.findSymbol = this.getSymbol(symbol)

	for {
		defer utils.PrintPanicStack()

		openorders, err := this.p.OpenOrders(symbol)
		if err != nil {
			time.Sleep(time.Minute)
			continue
		}

		var succ int
		for i, open := range openorders {
			price, err := decimal.NewFromString(open.Price)
			if err != nil {
				panic(err)
			}
			origQuantity, err := decimal.NewFromString(open.OrigQuantity)
			if err != nil {
				panic(err)
			}

			p1, err := this.p.GetMarketDepth(symbol)
			if err != nil {
				panic(err)
			}
			var curPrice = decimal.NewFromFloat(p1.Tick.Bids[0][0])

			var orderlog string
			var rescuelog string

			//SELL
			if open.Side == "SELL" {
				downPrice := price.Sub(curPrice).DivRound(decimal.NewFromFloat(1), 2)
				downRate, _ := downPrice.DivRound(price, 5).Mul(decimal.NewFromFloat(100)).Float64()
				lossusdtamount := downPrice.Mul(origQuantity).DivRound(decimal.NewFromFloat(1), 5)
				orderlog = fmt.Sprintf("%v,%v,%v,%v,%v,%v", open.OrderID, price.String(), "FILLED", price.Sub(decimal.NewFromFloat(2)).String(), origQuantity.String(), price.Mul(origQuantity).DivRound(decimal.NewFromInt(1), 2).String())
				rescuelog = fmt.Sprintf("时间:%v 订单:%v(%v) \t%v挂单 下降:%v(%.3f%%) 亏损:%vu \trescuelog: %v",
					time.Unix(open.Time/1000, 0).Format("2006-01-02 15:04:05"), origQuantity.String(), open.Price, open.Side,
					downPrice, downRate, lossusdtamount.String(), orderlog)
				if downRate < 0.4 {
					continue
				}
			}
			//BUY
			if open.Side == "BUY" {
				upPrice := curPrice.Sub(price).DivRound(decimal.NewFromFloat(1), 2)
				upRate, _ := upPrice.DivRound(price, 5).Mul(decimal.NewFromFloat(100)).Float64()
				lossusdtamount := upPrice.Mul(origQuantity).DivRound(decimal.NewFromFloat(1), 5)
				orderlog = fmt.Sprintf("%v,%v,%v,%v,%v,%v", "FILLED", price.Add(decimal.NewFromFloat(2)).String(), open.OrderID, price.String(), origQuantity.String(), price.Mul(origQuantity).DivRound(decimal.NewFromInt(1), 2).String())
				rescuelog = fmt.Sprintf("时间:%v 订单:%v(%v) \t%v挂单  上升:%v(%.3f%%) 亏损:%vu \trescuelog: %v",
					time.Unix(open.Time/1000, 0).Format("2006-01-02 15:04:05"), origQuantity.String(), open.Price, open.Side,
					upPrice, upRate, lossusdtamount.String(), orderlog)
				if upRate < 0.4 {
					continue
				}
			}

			for {
				if this.CheckPrice(symbol) == true {
					break
				}
				time.Sleep(time.Second * 60)
			}

			p1, err = this.p.GetMarketDepth(symbol)
			if err != nil {
				panic(err)
			}
			buyPrice, _ := decimal.NewFromFloat(p1.Tick.Bids[0][0]-2).DivRound(decimal.NewFromFloat(1), this.findSymbol.PricePrecision).Float64()
			sellPrice, _ := decimal.NewFromFloat(p1.Tick.Asks[0][0]+2).DivRound(decimal.NewFromFloat(1), this.findSymbol.PricePrecision).Float64()

			logmsg := fmt.Sprintf("%v %v/%v %v %v", symbol, i+1, len(openorders), "开始救单", rescuelog)
			message.PrintLog(api.PrintLog, fmt.Sprintf("rescueOpenOrders:%v,%v", i+1, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)

			if this.timerRescueOrder(this.findSymbol, i+1, len(openorders), symbol, buyPrice, sellPrice, orderlog) == true {
				succ++
				this.timerRescueSucc++
				time.Sleep(time.Second * 60 * 5)
			}
		}

		if succ > 0 {
			time.Sleep(time.Second * 3600 * 4)
		} else {
			time.Sleep(time.Second * 60)
		}
	}
}

func (this *OnePlatService) timerRescueOrder(findSymbol model.SymbolsData, group int, groupCount int, symbol string, buyprice, sellprice float64, rescuelog string) bool {
	oneplatService := NewOnePlatService(this.p.PlatCode())
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
		return false
	}

	//oneplatService.rescueMap.Store(oneplatService.genRotId(group, symbol, 3*1000-1), time.Now().Unix())
	oneplatService.orderPrice.Store(oneplatService.genRotId(group, symbol, 3*1000-1), &OrderInfo{SellOrderId: &sellorderid, SellPrice: oneplatService.manual_cursellprice, BuyOrderId: &buyorderid, BuyPrice: oneplatService.manual_curbuyprice})
	defer func() {
		oneplatService.orderPrice.Delete(oneplatService.genRotId(group, symbol, 3*1000-1))
	}()

	//logmsg := fmt.Sprintf("%v %v/%v %v %v", symbol, group, groupCount, "开始救单", rescuelog)
	//utils.PrintLog(api.PrintLog, fmt.Sprintf("rescueOpenOrders:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)

	btime := time.Now().Unix()
	ok, _ := oneplatService.rescueLastOrder(3*1000, group, 24*60*60, symbol, findSymbol, &sellorderid, lastsellprice, &buyorderid, lastbuyprice, lastbtcamount, lastusdtamount)
	if ok {
		logmsg := fmt.Sprintf("%v %v/%v %v %v %v", symbol, group, groupCount, "救单成功", fmt.Sprintf("用时:%v(s)", time.Now().Unix()-btime), rescuelog)
		message.PrintLog(api.PrintLog, fmt.Sprintf("rescueOpenOrders:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
		return true
	} else {
		logmsg := fmt.Sprintf("%v %v/%v %v %v", symbol, group, groupCount, "救单失败", rescuelog)
		message.PrintLog(api.PrintLog, fmt.Sprintf("rescueOpenOrders:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
		return false
	}
	//time.Sleep(time.Millisecond * 200)
}

// 套利下单
func (this *OnePlatService) arbitrageOrder(rescuecount int, group int, symbol string, symboldata model.SymbolsData, pricediff float64, sellprice float64, buyprice float64, btcamount, usdtamount float64) (timeout bool, err error) {
	var sellorderid string
	var buyorderid string
	var rescue_end bool
	var rescue_limithour_1 float64 = 1.5 //1级救单超时时间
	var rescue_limithour_2 float64 = 0.5 //2级救单超时时间

	if rescuecount == 0 && api.OnePlatStop == true {
		return false, fmt.Errorf("\t\t%v:group#%v (%v次) 下单失败:已暂停------%v \n", symbol, group, rescuecount, usdtamount)
	}

	if this.findSymbol.Symbol == "" {
		this.findSymbol = this.getSymbol(symbol)
		makerfee, takerfee, err := this.p.GetTradeFee(symbol)
		if err != nil {
			return false, fmt.Errorf("\t\t%v:group#%v (%v次) 下单失败:GetTradeFee------%v \n", symbol, group, rescuecount, err.Error())
		}
		this.makerFee, _ = makerfee.Float64()
		this.takerFee, _ = takerfee.Float64()
	}

	if len(strings.Split(symbol, "/")) != 2 {
		panic(fmt.Sprintf("symbol值没有/分隔符:%v", symbol))
	}
	symbolBase := strings.Split(symbol, "/")[0]
	symbolQuote := strings.Split(symbol, "/")[1]

	// 小于最小下单量
	if usdtamount <= symboldata.MinQuoteOrder {
		//logmsg := fmt.Sprintf("\t\tgroup:%v (%v次) minQuoteOrder------%v \n", group, rescuecount, usdtamount)
		//utils.PrintLog(fmt.Sprintf("every_order0:%v", group), fmt.Sprintf("%v", time.Now().Unix()), logmsg)
		//time.Sleep(1000 * time.Millisecond)
		return false, fmt.Errorf("\t\tgroup#%v (%v次) 下单失败:minQuoteOrder------%v \n", group, rescuecount, usdtamount)
	}

	var insufficientBalance bool = true
	var balance map[string]model.Balance
	for i := 0; i < 1; i++ {
		var err error
		balance, err = this.p.GetAccountBalance(symbol, false)
		if err != nil {
			fmt.Println("-----getbalance--------------------------------", err, time.Now().Format(time.Stamp))
		} else {
			if balance[symbolBase].Free.Cmp(decimal.NewFromFloat(btcamount)) > 0 && balance[symbolQuote].Free.Cmp(decimal.NewFromFloat(usdtamount)) > 0 {
				insufficientBalance = false
			}
			break
		}
	}
	//余额不足
	if api.CheckBalance == true && insufficientBalance == true {
		if balance == nil {
			//utils.PrintLog(fmt.Sprintf("every_order1:%v", group), fmt.Sprintf("%v", time.Now().Unix()/(60)), fmt.Sprintf("group:%v 获取余额接口出错,无法完成下单", group))
			//time.Sleep(time.Second)
			return false, fmt.Errorf("\t\tgroup#%v (%v次) 下单失败:获取余额接口出错,无法完成下单 \n", group, rescuecount)
		}

		logmsg := fmt.Sprintf("%v:group:%v (%v次):余额不足无法下单 %v/%v(%v,%v),(%v,%v)", symbol, group, rescuecount, symbolBase, symbolQuote, balance[symbolBase].Free, balance[symbolQuote].Free, btcamount, usdtamount)
		message.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_order2:%v", symbol), fmt.Sprintf("%v", time.Now().Unix()/(60*60)), logmsg)
		message.PrintLog(api.PrintLog, fmt.Sprintf("every_order2:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/(60*60)), logmsg)

		//余额不足,刷新缓存余额
		this.p.GetAccountBalance(symbol, true)
		time.Sleep(time.Second * 60)

		return false, fmt.Errorf("\t\t%v:group#%v (%v次) 下单失败:余额不足无法下单 \n", symbol, group, rescuecount)
	}

	//是否存在自成交
	var selfTradeOrderInfo string
	var selfTradeOrderCount int
	var selfTrade bool
	var openorders int
	this.orderPrice.Range(func(key, value any) bool {
		openorders++
		orderprice := value.(*OrderInfo)
		//之前的买单未成交之前，新卖价过小会发生自交易
		if *orderprice.BuyOrderId != "" && *orderprice.BuyOrderId != "FILLED" && sellprice <= orderprice.BuyPrice {
			selfTrade = true
			selfTradeOrderInfo = fmt.Sprintf("(%v,%v)/%v", orderprice.SellPrice, orderprice.BuyPrice, key)
			selfTradeOrderCount++
		}
		//之前的卖单未成交之前，新买价过大会发生自交易
		if *orderprice.SellOrderId != "" && *orderprice.SellOrderId != "FILLED" && buyprice >= orderprice.SellPrice {
			selfTrade = true
			selfTradeOrderInfo = fmt.Sprintf("(%v,%v)/%v", orderprice.SellPrice, orderprice.BuyPrice, key)
			selfTradeOrderCount++
		}
		return true
	})

	if rescuecount == 0 {
		// --开新单

		//同时只能有一个协程下单
		this.orderPlaceLock.Lock()
		defer func() {
			//解锁
			this.orderPlaceLock.TryLock()
			this.orderPlaceLock.Unlock()
		}()

		// 阻止自成交订单(非救单任务)
		if selfTrade == true {
			logmsg := fmt.Sprintf("\t\tgroup#%v (%v次) 下单失败:possible self-trading:%v \t %v(%v,%v)<--->%v\n", group, rescuecount, selfTradeOrderCount, usdtamount, sellprice, buyprice, selfTradeOrderInfo)
			return false, fmt.Errorf("%v", logmsg)
		}

		// 订单数上限判断
		apiopenorders, err := this.p.OpenOrders(symbol)
		if err != nil {
			//随机停止时间
			err_min := 5
			err_max := 15
			randSleepTime := time.Duration(rand.Intn(err_max-err_min) + err_min)

			logmsg := fmt.Sprintf("\t\t%vgroup#%v, -----err----- %v \n", symbol, group, err.Error())
			message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()), logmsg)
			time.Sleep(time.Second * randSleepTime)
			return false, errors.New(logmsg)
		}

		this.apiopenordercount = len(apiopenorders)

		if this.apiopenordercount >= int(this.findSymbol.MaxNumOrders-20) {
			logmsg := fmt.Sprintf("\t\t%v:group#%v 已超maxNumOrders------%v \n", symbol, group, this.apiopenordercount)
			message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit_stop:%v,%v", "group", symbol), fmt.Sprintf("%v", time.Now().Unix()/3600), logmsg)
			message.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_profit_stop:%v,%v", "group", symbol), fmt.Sprintf("%v", time.Now().Unix()/3600), logmsg)
			time.Sleep(time.Second * 60)
			return false, errors.New(logmsg)
		}

		// 阻止自成交订单:对挂单列表在进行一次自成交验证
		selfTradeOrderInfo = ""
		selfTradeOrderCount = 0
		for _, open := range apiopenorders {
			price, _ := decimal.NewFromString(open.Price)
			fprice, _ := price.Float64()
			//之前的买单未成交之前，新卖价过小会发生自交易
			if open.Side == "BUY" && sellprice <= fprice {
				selfTrade = true
				selfTradeOrderInfo = fmt.Sprintf("挂单买:%v 新单卖:%v ", open.Price, sellprice)
				selfTradeOrderCount++
			}
			//之前的卖单未成交之前，新买价过大会发生自交易
			if open.Side == "SELL" && buyprice >= fprice {
				selfTrade = true
				selfTradeOrderInfo = fmt.Sprintf("挂单卖:%v 新单买:%v ", open.Price, buyprice)
				selfTradeOrderCount++
			}
		}

		if selfTrade == true {
			//阻止自成交订单(非救单任务)
			logmsg := fmt.Sprintf("\t\tgroup#%v (%v次) 下单失败:possible2 self-trading:%v \t %v(%v,%v)<--->%v\n", group, rescuecount, selfTradeOrderCount, usdtamount, sellprice, buyprice, selfTradeOrderInfo)
			return false, fmt.Errorf("%v", logmsg)
		}

		// 铺单买价最小间隔:买价十万分之3
		var minIntervalPrice, _ = decimal.NewFromFloat(buyprice).Mul(decimal.NewFromFloat(3 / (10 * 10000.0))).Float64()
		var minIntervalCount int
		for _, open := range apiopenorders {
			price, err := decimal.NewFromString(open.Price)
			curBuyIntervalPrice, _ := decimal.NewFromFloat(buyprice).Sub(price).Float64()
			if err == nil && math.Abs(curBuyIntervalPrice) < minIntervalPrice {
				minIntervalCount++
				continue
			}
			curSellIntervalPrice, _ := decimal.NewFromFloat(buyprice).Sub(price).Float64()
			if err == nil && math.Abs(curSellIntervalPrice) < minIntervalPrice {
				minIntervalCount++
				continue
			}
		}
		//最多2单
		if minIntervalCount > 0 {
			logmsg := fmt.Sprintf("\t\t%v:group#%v 价格太近------(%v,%v)\n", symbol, group, sellprice, buyprice)
			//utils.PrintLog(api.PrintLog, fmt.Sprintf("every_profit:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()), logmsg)
			return false, errors.New(logmsg)
		}

		// 最后在与实时价进行对比，差价大不下单（可能过期价或振幅大）
		p1, _ := this.p.GetMarketDepth(symbol)
		if err != nil {
			return false, err
		}
		if p1.Tick.Bids[0][0] == 0 || len(p1.Tick.Bids) < 1 {
			return false, errors.New("bids is err")
		}
		if p1.Tick.Asks[0][0] == 0 || len(p1.Tick.Asks) < 1 {
			return false, errors.New("asks is err")
		}
		var buyPrice, _ = decimal.NewFromFloat(p1.Tick.Bids[0][0]-0).DivRound(decimal.NewFromFloat(1), this.findSymbol.PricePrecision).Float64()
		var sellPrice, _ = decimal.NewFromFloat(p1.Tick.Asks[0][0]+0).DivRound(decimal.NewFromFloat(1), this.findSymbol.PricePrecision).Float64()
		if math.Abs(buyPrice-buyprice) > 5 {
			return false, errors.New(fmt.Sprintf("buyprice is timeout %v %v %v", buyPrice, buyprice, math.Abs(buyPrice-buyprice)))
		}
		if math.Abs(sellPrice-sellprice) > 5 {
			return false, errors.New(fmt.Sprintf("sellprice is timeout %v %v %v", sellPrice, sellprice, math.Abs(sellPrice-sellprice)))
		}

		//解锁
		this.orderPlaceLock.TryLock()
		this.orderPlaceLock.Unlock()

	} else {
		// --救单任务

		if selfTrade == true {
			//救单任务打印日志
			logmsg := fmt.Sprintf("\t\tgroup#%v (%v次) 救单自成交log:possible self-trading:%v \t %v(%v,%v)<--->%v\n", group, rescuecount, selfTradeOrderCount, usdtamount, sellprice, buyprice, selfTradeOrderInfo)
			message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit_self-trading:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
		}
	}

	logmsg := fmt.Sprintf("group#%v B... %v/%v(%v,%v),(%v,%v) diff:%vu(%v,%v) 运行挂单:%v", group, symbolBase, symbolQuote, balance[symbolBase].Free, utils.DivRound(balance[symbolQuote].Free, 2), btcamount, usdtamount, pricediff, sellprice, buyprice, openorders+1)
	message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
	//logmsg := fmt.Sprintf("group#%v 下单成功", group)
	//utils.PrintLog(fmt.Sprintf("every_profit:%v", group), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)

	//开始下单
	var self_rescuecount = rescuecount
	this.orderPrice.Store(this.genRotId(group, symbol, self_rescuecount), &OrderInfo{CreateTime: time.Now().Unix(), SellOrderId: &sellorderid, SellPrice: sellprice, BuyOrderId: &buyorderid, BuyPrice: buyprice})
	defer func() {
		//清除已成交订单
		this.orderPrice.Delete(this.genRotId(group, symbol, self_rescuecount))
	}()

	var btime = time.Now() //下单开始时间
	var waitGroup sync.WaitGroup

	//下卖单
	waitGroup.Add(1)
	go func() {
		//defer utils.PrintPanicStack()
		var err error
		for {
			//买到币扣掉手续费后为可卖数
			sellbtcamount := decimal.NewFromFloat(btcamount).Mul(decimal.NewFromFloat(1-this.makerFee)).DivRound(decimal.NewFromInt(1), this.findSymbol.AmountPrecision).String()
			sellorderid, err = this.p.OrdersPlace(symbol, fmt.Sprintf("%v", sellprice), fmt.Sprintf("%v", sellbtcamount), "SELL")
			if err == nil {
				break
			}

			//余额不足,刷新缓存余额
			commerr, ok := err.(*common.APIError)
			if ok && commerr.Code == -2010 {
				this.p.GetAccountBalance(symbol, true)
				time.Sleep(time.Second * 60)
			}

			//随机停止时间
			err_min := 5
			err_max := 15
			randSleepTime := time.Duration(rand.Intn(err_max-err_min) + err_min)
			time.Sleep(time.Second * randSleepTime)
		}
		if sellorderid != "FILLED" {
			for {
				var usetime = time.Now().Unix() - btime.Unix()
				var stoptime = int64((rescue_limithour_1 + rescue_limithour_2 + 0.003) * 60 * 60)
				//var stoptime = int64(36 * 60 * 60)
				//超时未完成，退出查询
				if usetime > stoptime+60 || rescue_end == true {
					sellorderid = "TIMEOUT" + sellorderid
					time.Sleep(time.Millisecond * 200)
					waitGroup.Done()
					break
				}
				//查询订单
				orderstatus, err := this.p.GetOrderStatus(symbol, sellorderid)
				if err != nil {
					//随机停止时间
					err_min := 5
					err_max := 20
					randSleepTime := time.Duration(rand.Intn(err_max-err_min) + err_min)
					time.Sleep(time.Second * randSleepTime)
					message.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_query_err,%v", symbol), fmt.Sprintf("%v", time.Now().Unix()/(12*60*60)), err.Error())
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
			buyorderid, err = this.p.OrdersPlace(symbol, fmt.Sprintf("%v", buyprice), fmt.Sprintf("%v", btcamount), "BUY")
			if err == nil {
				break
			}

			//余额不足,刷新缓存余额
			commerr, ok := err.(*common.APIError)
			if ok && commerr.Code == -2010 {
				this.p.GetAccountBalance(symbol, true)
				time.Sleep(time.Second * 60)
			}

			//随机停止时间
			err_min := 5
			err_max := 15
			randSleepTime := time.Duration(rand.Intn(err_max-err_min) + err_min)
			time.Sleep(time.Second * randSleepTime)
		}
		if buyorderid != "FILLED" {
			for {
				var usetime = time.Now().Unix() - btime.Unix()
				var stoptime = int64((rescue_limithour_1 + rescue_limithour_2 + 0.003) * 60 * 60)
				//var stoptime = int64(36 * 60 * 60)
				//超时未完成，退出查询
				if usetime > stoptime+60 || rescue_end == true {
					buyorderid = "TIMEOUT" + buyorderid
					time.Sleep(time.Millisecond * 200)
					waitGroup.Done()
					break
				}
				//查询订单
				orderstatus, err := this.p.GetOrderStatus(symbol, buyorderid)
				if err != nil {
					//随机停止时间
					err_min := 5
					err_max := 20
					randSleepTime := time.Duration(rand.Intn(err_max-err_min) + err_min)
					time.Sleep(time.Second * randSleepTime)
					message.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_query_err,%v", symbol), fmt.Sprintf("%v", time.Now().Unix()/(12*60*60)), err.Error())
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
			var timeout_second = int64(20*rescuecount+this.firstRescueBeginSec/60) * 60
			var usetime = time.Now().Unix() - btime.Unix()
			var remainingtime = timeout_second - usetime
			if usetime > timeout_second {
				if rescuecount >= 2 && rescuecount < 1000 {
					timeout = true
					logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次):救单超时 用时:%v(s) rescuelog: %v,%v,%v,%v,%v,%v", symbol, group, rescuecount, usetime, *sellorderId, sellprice, *buyorderId, buyprice, btcamount, usdtamount)
					message.PrintLog(api.PrintLog, fmt.Sprintf("every_order6:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()), logmsg)
					return
				}

				//开始救单
				rescuecount++

				//救单时间
				var rescue_limittime int64
				if rescuecount == 1 {
					//拯救n小时
					rescue_limittime = int64(rescue_limithour_1 * 60 * 60)
				} else {
					//拯救n小时
					rescue_limittime = int64(rescue_limithour_2 * 60 * 60)
				}

				//救单开始
				this.rescueMap.Store(this.genRotId(group, symbol, rescuecount), time.Now().Unix())
				rescuebtime := time.Now().Unix()
				rescue, craterescueorder := this.rescueLastOrder(rescuecount, group, rescue_limittime, symbol, symboldata, sellorderId, sellprice, buyorderId, buyprice, btcamount, usdtamount)
				rescueetime := time.Now().Unix()
				//救单结束
				timeout = !rescue
				this.rescueMap.Delete(this.genRotId(group, symbol, rescuecount))
				rescue_end = true

				var logmsg string
				if rescue == true {
					//创建了救单属于救单成功，否则属于自动完成或已超时
					if craterescueorder == true {
						logmsg = fmt.Sprintf("\t\t%v:group:%v (%v次):救单成功 用时:%v(s) rescuelog: %v,%v,%v,%v,%v,%v", symbol, group, rescuecount, rescueetime-rescuebtime, *sellorderId, sellprice, *buyorderId, buyprice, btcamount, usdtamount)
					} else {
						if *sellorderId == "FILLED" && *buyorderId == "FILLED" {
							logmsg = fmt.Sprintf("\t\t%v:group:%v (%v次):订单已自动成交 用时:%v(s) rescuelog: %v,%v,%v,%v,%v,%v", symbol, group, rescuecount, rescueetime-rescuebtime, *sellorderId, sellprice, *buyorderId, buyprice, btcamount, usdtamount)
						} else {
							logmsg = fmt.Sprintf("\t\t%v:group:%v (%v次):超时未成交,已取消买卖挂单 用时:%v(s) rescuelog: %v,%v,%v,%v,%v,%v", symbol, group, rescuecount, rescueetime-rescuebtime, *sellorderId, sellprice, *buyorderId, buyprice, btcamount, usdtamount)
						}
					}
				} else {
					if *sellorderId == "FILLED" && *buyorderId == "FILLED" {
						logmsg = fmt.Sprintf("\t\t%v:group:%v (%v次):订单已自动成交 用时:%v(s) rescuelog: %v,%v,%v,%v,%v,%v", symbol, group, rescuecount, rescueetime-rescuebtime, *sellorderId, sellprice, *buyorderId, buyprice, btcamount, usdtamount)
					} else {
						logmsg = fmt.Sprintf("\t\t%v:group:%v (%v次):救单超时 用时:%v(s) rescuelog: %v,%v,%v,%v,%v,%v  createtime:%v", symbol, group, rescuecount, rescueetime-rescuebtime, *sellorderId, sellprice, *buyorderId, buyprice, btcamount, usdtamount, btime.Format("2006-01-02 15:04:05"))
						//utils.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_timeout:%v", group), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
					}
				}
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_order6:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg+"\n\n")
				return
			} else {
				if *sellorderId == "FILLED" && *buyorderId == "FILLED" {
					return
				} else {
					//进行中
					if rescuecount == 0 {
						logmsg := fmt.Sprintf("\t\t%v:group:%v .... 用时:%v(s) rescuelog: %v,%v,%v,%v,%v,%v  createtime:%v", symbol, group, usetime, *sellorderId, sellprice, *buyorderId, buyprice, btcamount, usdtamount, btime.Format("2006-01-02 15:04:05"))
						message.PrintLog(api.PrintLog, fmt.Sprintf("every_order7:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/(60*5)), logmsg)
					} else {
						logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次)救单进行中..... 用时:%v(s),在次开启剩余:%v(s) rescuelog: %v,%v,%v,%v,%v,%v  createtime:%v", symbol, group, rescuecount, usetime, remainingtime, *sellorderId, sellprice, *buyorderId, buyprice, btcamount, usdtamount, btime.Format("2006-01-02 15:04:05"))
						message.PrintLog(api.PrintLog, fmt.Sprintf("every_order8:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/(60*5)), logmsg)
					}
					time.Sleep(time.Millisecond * 1000)
				}
			}

		}
	}(&sellorderid, &buyorderid)

	waitGroup.Wait()

	//清除已成交订单
	this.orderPrice.Delete(this.genRotId(group, symbol, self_rescuecount))
	time.Sleep(time.Millisecond * 1500)

	//纠正
	if timeout == true && buyorderid == "FILLED" && sellorderid == "FILLED" {
		timeout = false
	}
	return timeout, nil
}

// 拯救下单
func (this *OnePlatService) rescueLastOrder(rescuecount int, group int, rescue_limittime int64, symbol string, symboldata model.SymbolsData, lastsellorderid *string, lastsellprice float64, lastbuyorderid *string, lastbuyprice float64, lastbtcamount, lastusdtamount float64) (bool, bool) {
	//两笔未均完成
	if *lastsellorderid != "FILLED" && *lastbuyorderid != "FILLED" {
		_, err1 := this.p.CancelOrder(symbol, *lastsellorderid)
		_, err2 := this.p.CancelOrder(symbol, *lastbuyorderid)
		if err1 != nil || err2 != nil {
			logmsg := fmt.Sprintf("\t\t%v:group:%v 撤单失败(%v,%v):sellorderid(%v):%v, buyorderid(%v):%v", symbol, group, lastusdtamount, lastbtcamount, *lastsellorderid, err1, *lastbuyorderid, err2)
			message.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_rescue:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
			message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
		}
		return true, false
	}

	var addProfit = 0.02 //补贴救单任务自成交吃单损失（最多9笔，超出将亏损）
	var maxpricediff float64 = 3
	var maxusdtamount float64 = 100
	//人工救单
	if this.manual_price_enable == true {
		addProfit = 0.05
		maxpricediff = 10
		maxusdtamount = 1000
	}

	//同时只开启一个救单任务
	var curbuyprice float64
	var cursellprice float64
	var curpricediff float64
	var startTime = time.Now().Unix()
	for {
		time.Sleep(time.Second * 1)

		usetime := time.Now().Unix() - startTime
		//超时救单(失败)
		if time.Now().Unix()-startTime > rescue_limittime {
			return false, false
		}

		//已完成，不需要在救单
		if *lastsellorderid == "FILLED" && *lastbuyorderid == "FILLED" {
			logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次):订单已自动成交 FILLED \t(rescuelog1: %v,%v,%v,%v,%v,%v)", symbol, group, rescuecount, *lastsellorderid, lastsellprice, *lastbuyorderid, lastbuyprice, lastbtcamount, lastusdtamount)
			message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
			//todo 可能在救单过程中自行又全成交了
			if rescuecount >= 2 && rescuecount < 1000 { //救单小队只能救单前一小队，订单本身等自动成交
				return false, false
			} else {
				return true, false
			}
		}

		//已完成，不需要在救单
		if *lastsellorderid == "" && *lastbuyorderid == "" {
			logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次):订单已自动成交 orderid已经为空 \t(rescuelog1: %v,%v,%v,%v,%v,%v)", symbol, group, rescuecount, *lastsellorderid, lastsellprice, *lastbuyorderid, lastbuyprice, lastbtcamount, lastusdtamount)
			message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
			//todo 可能在救单过程中自行又全成交了
			if rescuecount >= 2 && rescuecount < 1000 { //救单小队只能救单前一小队，订单本身等自动成交
				return false, false
			} else {
				return true, false
			}
		}

		//已经不存在，说明上笔订单已经完成，不需要在救单
		if _, ok := this.orderPrice.Load(this.genRotId(group, symbol, rescuecount-1)); ok == false {
			logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次):orderPrice不存在，表示已自动完成或已超时 \t(rescuelog1: %v,%v,%v,%v,%v,%v)", symbol, group, rescuecount, *lastsellorderid, lastsellprice, *lastbuyorderid, lastbuyprice, lastbtcamount, lastusdtamount)
			message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
			//todo 可能在救单过程中自行又全成交了
			if rescuecount >= 2 && rescuecount < 1000 { //救单小队只能救单前一小队，订单本身等自动成交
				return false, false
			} else {
				if strings.Index(*lastsellorderid, "TIMEOUT") >= 0 || strings.Index(*lastbuyorderid, "TIMEOUT") >= 0 {
					return false, false
				} else {
					return true, false
				}
			}
		}

		if this.manual_price_enable == true {
			//人工救单
			curbuyprice = this.manual_curbuyprice
			cursellprice = this.manual_cursellprice
		} else {
			//系统救单
			for {
				p1, err := this.p.GetMarketDepth(symbol)
				if err != nil {
					continue
				}
				if p1.Tick.Bids[0][0] == 0 || len(p1.Tick.Bids) < 1 {
					continue
				}
				if p1.Tick.Asks[0][0] == 0 || len(p1.Tick.Asks) < 1 {
					continue
				}
				curbuyprice, _ = decimal.NewFromFloat(p1.Tick.Bids[0][0]).DivRound(decimal.NewFromFloat(1), symboldata.PricePrecision).Float64()
				cursellprice, _ = decimal.NewFromFloat(p1.Tick.Asks[0][0]).DivRound(decimal.NewFromFloat(1), symboldata.PricePrecision).Float64()
				break
			}
		}

		//差价太小
		curpricediff, _ = decimal.NewFromFloat(cursellprice-curbuyprice).DivRound(decimal.NewFromFloat(1), 8).Float64()
		if this.manual_price_enable == false && curpricediff < 0.2 {
			continue
		}

		//出现负值 或 价差太大/用U太多，暂不进入救单队列
		var curusdtamount float64
		if *lastsellorderid == "FILLED" {
			//上次买单没有成交，算卖单损失u，进行补救

			////卖单已成交，说明在上涨
			//if this.manual_price_enable == false {
			//	curbuyprice, _ = decimal.NewFromFloat(curbuyprice).Add(decimal.NewFromFloat(0.2)).Float64()
			//	cursellprice, _ = decimal.NewFromFloat(cursellprice).Add(decimal.NewFromFloat(0.8)).Float64()
			//	curpricediff, _ = decimal.NewFromFloat(cursellprice-curbuyprice).DivRound(decimal.NewFromFloat(1), 8).Float64()
			//}

			//计算亏损:(当前卖一减去上次卖单卖价)*上次卖币数
			lostpricediff, _ := decimal.NewFromFloat(cursellprice-lastsellprice).DivRound(decimal.NewFromInt(1), 8).Float64()
			lossusdtamount, _ := decimal.NewFromFloat(lostpricediff*lastbtcamount+addProfit).DivRound(decimal.NewFromFloat(1), 5).Float64()
			//计算下单量：亏损除以当前买一卖一差价，成交后重新下单并取消上笔交易
			curbtcamount, _ := decimal.NewFromFloat(lossusdtamount).DivRound(decimal.NewFromFloat(curpricediff), symboldata.AmountPrecision).Float64()
			//需要U数量
			curusdtamount, _ = decimal.NewFromFloat(curbtcamount).Mul(decimal.NewFromFloat(curbuyprice)).Float64()

			upPrice := decimal.NewFromFloat(cursellprice-lastbuyprice).DivRound(decimal.NewFromFloat(1), 2)
			upRate, _ := upPrice.DivRound(decimal.NewFromFloat(lastbuyprice), 5).Mul(decimal.NewFromFloat(100)).Float64()
			//todo 未到亏损比例，比如0.3%不触发救单
			if this.manual_price_enable == false && upRate < 0.15 {
				logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次):未到亏损比例，还未加入... 用时:%v(s) \t(rescuelog1: %v,%v,%v,%v,%v,%v 上升:%vu(%.3f%%) 上卖单亏:%v) 现价差:%v 需要:%v(等值U:%v)",
					symbol, group, rescuecount, usetime,
					*lastsellorderid, lastsellprice, *lastbuyorderid, lastbuyprice, lastbtcamount, lastusdtamount,
					upPrice, upRate, lossusdtamount, curpricediff, curbtcamount, curusdtamount)
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue_test:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/(60*5)), logmsg)
				if rescuecount > 1 {
					time.Sleep(time.Second * 10)
				}
				continue
			}

			//上次买单已自动成交
			if this.manual_price_enable == true {
				if orderstatus, _ := this.p.GetOrderStatus(symbol, *lastbuyorderid); orderstatus == true {
					return true, false
				}
			}

			if curpricediff > maxpricediff || curusdtamount > maxusdtamount {
				if this.manual_price_enable == true {
					logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次):价差太大/用U太多，无法人工救单... 用时:%v(s) \t(rescuelog1: %v,%v,%v,%v,%v,%v 上升:%vu(%.3f%%) 上卖单亏:%v) 现价差:%v 需要:%v(等值U:%v)", symbol, group, rescuecount, usetime, *lastsellorderid, lastsellprice, *lastbuyorderid, lastbuyprice, lastbtcamount, lastusdtamount, upPrice, upRate, lossusdtamount, curpricediff, curbtcamount, curusdtamount)
					message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue_test:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()), logmsg)
					return false, false
				}
				logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次):价差太大/用U太多，还未加入... 用时:%v(s) \t(rescuelog1: %v,%v,%v,%v,%v,%v 上升:%vu(%.3f%%) 上卖单亏:%v) 现价差:%v 需要:%v(等值U:%v)",
					symbol, group, rescuecount, usetime,
					*lastsellorderid, lastsellprice, *lastbuyorderid, lastbuyprice, lastbtcamount, lastusdtamount,
					upPrice, upRate, lossusdtamount, curpricediff, curbtcamount, curusdtamount)
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue_test:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/(60*5)), logmsg)
				if rescuecount > 1 {
					time.Sleep(time.Second * 10)
				}
				continue
			}
		}
		if *lastbuyorderid == "FILLED" {
			//上次卖单没有成交，算买单损失u，进行补救

			////买单已成交，说明在下降
			//if this.manual_price_enable == false {
			//	curbuyprice, _ = decimal.NewFromFloat(curbuyprice).Sub(decimal.NewFromFloat(0.8)).Float64()
			//	cursellprice, _ = decimal.NewFromFloat(cursellprice).Sub(decimal.NewFromFloat(0.2)).Float64()
			//	curpricediff, _ = decimal.NewFromFloat(cursellprice-curbuyprice).DivRound(decimal.NewFromFloat(1), 8).Float64()
			//}

			//计算亏损U:(上次买单买价减去当前买一)*上次买币数
			lostpricediff, _ := decimal.NewFromFloat(lastbuyprice-curbuyprice).DivRound(decimal.NewFromInt(1), 8).Float64()
			lossusdtamount, _ := decimal.NewFromFloat(lostpricediff*lastbtcamount+addProfit).DivRound(decimal.NewFromFloat(1), 5).Float64()
			//计算下单量：亏损除以当前买一卖一差价，成交后重新下单并取消上笔交易,
			curbtcamount, _ := decimal.NewFromFloat(lossusdtamount).DivRound(decimal.NewFromFloat(curpricediff), symboldata.AmountPrecision).Float64()
			//需要U数量
			curusdtamount, _ = decimal.NewFromFloat(curbtcamount).Mul(decimal.NewFromFloat(curbuyprice)).Float64()

			downPrice := decimal.NewFromFloat(lastsellprice-curbuyprice).DivRound(decimal.NewFromFloat(1), 2)
			downRate, _ := downPrice.DivRound(decimal.NewFromFloat(curbuyprice), 5).Mul(decimal.NewFromFloat(100)).Float64()
			//todo 未到亏损比例，比如0.3%不触发救单
			if this.manual_price_enable == false && downRate < 0.15 {
				logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次):未到亏损比例，还未加入... 用时:%v(s) \t(rescuelog1: %v,%v,%v,%v,%v,%v 下降:%vu(%.3f%%) 上买单亏:%v) 现价差:%v 需要:%v(等值U:%v)",
					symbol, group, rescuecount, usetime,
					*lastsellorderid, lastsellprice, *lastbuyorderid, lastbuyprice, lastbtcamount, lastusdtamount,
					downPrice, downRate, lossusdtamount, curpricediff, curbtcamount, curusdtamount)
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue_test:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/(60*5)), logmsg)
				if rescuecount > 1 {
					time.Sleep(time.Second * 10)
				}
				continue
			}

			//上次卖单已自动成交
			if this.manual_price_enable == true {
				if orderstatus, _ := this.p.GetOrderStatus(symbol, *lastsellorderid); orderstatus == true {
					return true, false
				}
			}

			if curpricediff > maxpricediff || curusdtamount > maxusdtamount {
				if this.manual_price_enable == true {
					logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次):价差太大/用U太多，无法人工救单... 用时:%v(s) \t(rescuelog1: %v,%v,%v,%v,%v,%v 下降:%vu(%.3f%%) 上买单亏:%v) 现价差:%v 需要:%v(等值U:%v)", symbol, group, rescuecount, usetime, *lastsellorderid, lastsellprice, *lastbuyorderid, lastbuyprice, lastbtcamount, lastusdtamount, downPrice, downRate, lossusdtamount, curpricediff, curbtcamount, curusdtamount)
					message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue_test:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()), logmsg)
					return false, false
				}
				logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次):价差太大/用U太多，还未加入... 用时:%v(s) \t(rescuelog1: %v,%v,%v,%v,%v,%v 下降:%vu(%.3f%%) 上买单亏:%v) 现价差:%v 需要:%v(等值U:%v)",
					symbol, group, rescuecount, usetime,
					*lastsellorderid, lastsellprice, *lastbuyorderid, lastbuyprice, lastbtcamount, lastusdtamount,
					downPrice, downRate, lossusdtamount, curpricediff, curbtcamount, curusdtamount)
				message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue_test:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/(60*5)), logmsg)
				if rescuecount > 1 {
					time.Sleep(time.Second * 10)
				}
				continue
			}
		}

		//救单需要资金不多，说明币价还在周围，仍可能自动成交,不需要救单机器人
		if this.manual_price_enable == false && curusdtamount < 200 {
			logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次): testlog 币价还在周围不用启动救单机器人(现价:%v) rescuelog1: %v,%v,%v,%v,%v,%v", symbol, group, rescuecount, cursellprice, *lastsellorderid, lastsellprice, *lastbuyorderid, lastbuyprice, lastbtcamount, lastusdtamount)
			message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue_test:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/(60*5)), logmsg)
			continue
		}

		//同时只开启一个救单任务
		if rescuecount == 1 && this.workingRescueRots > 0 {
			startTime = time.Now().Unix() //重新计时
			logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次): testlog 还有其它救单机器人 rescuelog1: %v,%v,%v,%v,%v,%v", symbol, group, rescuecount, *lastsellorderid, lastsellprice, *lastbuyorderid, lastbuyprice, lastbtcamount, lastusdtamount)
			message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue_test:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().Unix()/(60*5)), logmsg)
			continue
		} else {
			break
		}
	}

	//统计救单机器人数量
	if rescuecount == 1 {
		this.workingRescueRots += 1
		defer func() {
			this.workingRescueRots -= 1
		}()
	}

	//-----------------------开始救单-----------------------------
	//取消上组未完成订单
	if *lastsellorderid == "FILLED" {
		//上次买单没有成交，算卖单损失u，进行补救

		////卖单已成交，说明在上涨
		//if this.manual_price_enable == false {
		//	curbuyprice, _ = decimal.NewFromFloat(curbuyprice).Add(decimal.NewFromFloat(0.2)).Float64()
		//	cursellprice, _ = decimal.NewFromFloat(cursellprice).Add(decimal.NewFromFloat(0.8)).Float64()
		//	curpricediff, _ = decimal.NewFromFloat(cursellprice-curbuyprice).DivRound(decimal.NewFromFloat(1), 8).Float64()
		//}

		//计算亏损:(当前卖一减去上次卖单卖价)*上次卖币数
		lostpricediff, _ := decimal.NewFromFloat(cursellprice-lastsellprice).DivRound(decimal.NewFromInt(1), 8).Float64()
		lossusdtamount, _ := decimal.NewFromFloat(lostpricediff*lastbtcamount+addProfit).DivRound(decimal.NewFromFloat(1), 5).Float64()
		//计算下单量：亏损除以当前买一卖一差价，成交后重新下单并取消上笔交易
		curbtcamount, _ := decimal.NewFromFloat(lossusdtamount).DivRound(decimal.NewFromFloat(curpricediff), symboldata.AmountPrecision).Float64()
		//需要U数量
		curusdtamount, _ := decimal.NewFromFloat(curbtcamount).Mul(decimal.NewFromFloat(curbuyprice)).Float64()

		logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次):正在拯救上笔低价卖单亏损:%v(u) 上笔卖价:%v 亏损=价差*订单量(%v(u))=%v*%v \t当前价差:%v 下单量:%v(等值U:%v)", symbol, group, rescuecount, lossusdtamount, lastsellprice, lastusdtamount, lostpricediff, lastbtcamount, curpricediff, curbtcamount, curusdtamount)
		message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)

		timeout, err := this.arbitrageOrder(rescuecount, group, symbol, symboldata, curpricediff, cursellprice, curbuyprice, curbtcamount, curusdtamount)
		if err != nil {
			logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次): 拯救下单失败 rescuelog1: %v,%v,%v,%v,%v,%v err:%v", symbol, group, rescuecount, *lastsellorderid, lastsellprice, *lastbuyorderid, lastbuyprice, lastbtcamount, lastusdtamount, err.Error())
			message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue_test:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
			return false, false
		}
		if timeout == false {
			if status, _ := this.p.GetOrderStatus(symbol, *lastbuyorderid); status == false {
				//上次卖出量按当前卖一价买入
				_, err := this.p.OrdersPlace(symbol, fmt.Sprintf("%v", cursellprice), fmt.Sprintf("%v", lastbtcamount), "BUY")
				if err == nil {
					_, err := this.p.CancelOrder(symbol, *lastbuyorderid)
					if err != nil {
						logmsg := fmt.Sprintf("\t\t%v:group:%v 撤单失败(%v,%v):buyorderid(%v):%v", symbol, group, lastusdtamount, lastbtcamount, *lastbuyorderid, err)
						message.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_rescue:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixMicro()), logmsg)
						message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixMicro()), logmsg)
					}
				}
			}
			//todo 可能在救单过程中自行又全成交了
			if rescuecount >= 2 && rescuecount < 1000 { //救单小队只能救单前一小队，订单本身只能等自动成交， 超1000为人工救单
				return false, true
			} else {
				return true, true
			}
		} else {
			return false, true
		}
	}

	if *lastbuyorderid == "FILLED" {
		//上次卖单没有成交，算买单损失u，进行补救

		////买单已成交，说明在下降
		//if this.manual_price_enable == false {
		//	curbuyprice, _ = decimal.NewFromFloat(curbuyprice).Sub(decimal.NewFromFloat(0.8)).Float64()
		//	cursellprice, _ = decimal.NewFromFloat(cursellprice).Sub(decimal.NewFromFloat(0.2)).Float64()
		//	curpricediff, _ = decimal.NewFromFloat(cursellprice-curbuyprice).DivRound(decimal.NewFromFloat(1), 8).Float64()
		//}

		//计算亏损U:(上次买单买价减去当前买一)*上次买币数
		lostpricediff, _ := decimal.NewFromFloat(lastbuyprice-curbuyprice).DivRound(decimal.NewFromInt(1), 8).Float64()
		lossusdtamount, _ := decimal.NewFromFloat(lostpricediff*lastbtcamount+addProfit).DivRound(decimal.NewFromFloat(1), 5).Float64()
		//计算下单量：亏损除以当前买一卖一差价，成交后重新下单并取消上笔交易,
		curbtcamount, _ := decimal.NewFromFloat(lossusdtamount).DivRound(decimal.NewFromFloat(curpricediff), symboldata.AmountPrecision).Float64()
		//需要U数量
		curusdtamount, _ := decimal.NewFromFloat(curbtcamount).Mul(decimal.NewFromFloat(curbuyprice)).Float64()

		logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次):正在拯救上笔高价买单亏损:%v(u) 上笔买价:%v 亏损=价差*订单量(%v(u))=%v*%v \t当前价差:%v,下单量:%v(等值U:%v)", symbol, group, rescuecount, lossusdtamount, lastbuyprice, lastusdtamount, lostpricediff, lastbtcamount, curpricediff, curbtcamount, curusdtamount)
		message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)

		timeout, err := this.arbitrageOrder(rescuecount, group, symbol, symboldata, curpricediff, cursellprice, curbuyprice, curbtcamount, curusdtamount)
		if err != nil {
			logmsg := fmt.Sprintf("\t\t%v:group:%v (%v次): 拯救下单失败 rescuelog1: %v,%v,%v,%v,%v,%v err:%v", symbol, group, rescuecount, *lastsellorderid, lastsellprice, *lastbuyorderid, lastbuyprice, lastbtcamount, lastusdtamount, err.Error())
			message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue_test:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
			return false, false
		}
		if timeout == false {
			if status, _ := this.p.GetOrderStatus(symbol, *lastsellorderid); status == false {
				//上次买入量按当前买一价卖掉
				_, err := this.p.OrdersPlace(symbol, fmt.Sprintf("%v", curbuyprice), fmt.Sprintf("%v", lastbtcamount), "SELL")
				if err == nil {
					_, err := this.p.CancelOrder(symbol, *lastsellorderid)
					if err != nil {
						logmsg := fmt.Sprintf("\t\t%v:group:%v 撤单失败(%v,%v):sellorderid(%v):%v", symbol, group, lastusdtamount, lastbtcamount, *lastsellorderid, err)
						message.SendDingTalkRobit(true, "oneplat", fmt.Sprintf("every_rescue:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixMicro()), logmsg)
						message.PrintLog(api.PrintLog, fmt.Sprintf("every_rescue:%v,%v", group, symbol), fmt.Sprintf("%v", time.Now().UnixMicro()), logmsg)
					} else {
						//logmsg := fmt.Sprintf("(%v次):救单成功 %v %v", rescuetime, *sellorderid, *buyorderid)
						//utils.PrintLog("every_cancel", fmt.Sprintf("%v", time.Now().UnixMicro()), logmsg)
					}
				}
			}
			//todo 可能在救单过程中自行又全成交了
			if rescuecount >= 2 && rescuecount < 1000 { //救单小队只能救单前一小队，订单本身只能等自动成交， 超1000为人工救单
				return false, true
			} else {
				return true, true
			}
		} else {
			return false, true
		}
	}

	return false, false
}

// 是否存在其它救单机器人
func (this *OnePlatService) existRescueRot(self_group int, self_symbol string) bool {
	var exist bool
	this.rescueMap.Range(func(key, value any) bool {
		group, symbol, _ := this.parseRotId(key.(string))
		if group == self_group && symbol == self_symbol {
			//是自已
		} else {
			exist = true
		}
		return true
	})
	return exist
}

func (this *OnePlatService) genRotId(group int, symbol string, rescuecount int) (rotid string) {
	return fmt.Sprintf("%v_%v_%v", group, symbol, rescuecount)
}

func (this *OnePlatService) parseRotId(rotid string) (group int, symbol string, rescuecount int) {
	keys := strings.Split(rotid, "_")
	group, _ = strconv.Atoi(keys[0])
	symbol = keys[1]
	rescuecount, _ = strconv.Atoi(keys[2])
	return
}
