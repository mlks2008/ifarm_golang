/**
 * dcaService.go
 * ============================================================================
 *1、在大跌之后，后面大概率一波反弹行情，这种情况调整首买单额，比如bnbusdt从3500上调到5000，arbusdt从700上调到1400。
 *2、日常平稳行情时，下跌买入间距调小，提升资金利用率，涉及三档
 *   arbusdt:高风险（间距小，控制在5%内回撤，适合平稳升和上涨行情）、中风险（间距适中，控制在10%内回撤，适合平稳降和下降行情），低风险（间距大，控制在15%以上回撤，适合下降行情）
 *   bnbusdt:高风险（间距小，控制在3%内回撤，适合平稳升和上涨行情）、中风险（间距适中，控制在6%内回撤，适合平稳降和下降行情），低风险（间距大，控制在10%以上回撤，适合下降行情）
 * ============================================================================
 * author: peter.wang
 * createtime: 2020-07-07 22:35
 */

package services

import (
	"components/message"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"goarbitrage/api"
	"goarbitrage/internal/model"
	"goarbitrage/internal/plat"
	"goarbitrage/pkg/utils"
	"strings"
	"sync"
	"time"
)

type DCAService struct {
	p              plat.Plat
	dayStartTime   int64
	dayTotalCount  int64
	dayTotalProfit decimal.Decimal
	totalProfit    decimal.Decimal
	startTime      int64
	totalCount     int64
	symbol         string //交易对,比如BTC/USDT
	symbolBase     string //BTC
	symbolQuote    string //USDT
	//初始投入
	initInput map[string]model.Balance
	//启动时帐户余额
	runBalance map[string]model.Balance
	//交易对
	symbols map[string]model.SymbolsData
	//挂单列表
	openorders sync.Map
}

func NewDCAService(platcode string, symbol string) *DCAService {
	symbolBase := strings.Split(symbol, "/")[0]
	symbolQuote := strings.Split(symbol, "/")[1]

	ser := &DCAService{
		p:            plat.Get(platcode),
		dayStartTime: time.Now().Unix(),
		startTime:    time.Now().Unix(),
		symbol:       symbol,
		symbolBase:   symbolBase,
		symbolQuote:  symbolQuote,
	}
	ser.initInput = ser.p.GetInitialInput(symbol)
	ser.runBalance, _ = ser.p.GetAccountBalance(symbol, true)
	ser.symbols = make(map[string]model.SymbolsData)
	return ser
}

// 买单下跌超0.4%补一单，收益超0.2%即卖
func (this *DCAService) Start() {
	//启动时依上次完成订单余额为启始值
	var run_b_quote, run_b_base = this.RunGetUsdtCost()
	var run_bBalance = make(map[string]model.Balance)
	run_bBalance[this.symbolQuote] = model.Balance{Total: run_b_quote}
	run_bBalance[this.symbolBase] = model.Balance{Total: run_b_base}
	//启动时间
	var run_startTime = time.Now().Unix()
	//每次交易开始时间(计算每次交易用时)
	var every_trading_startTime = time.Now().Unix()
	//上次买单时间(计算买单时间间隔，超时降低buypoint)
	var last_buy_time = this.GetLastBuyTime()

	var trailing = false               //进入回撤时刻
	var trailingMaxRate = decimal.Zero //进入回撤时刻最高涨幅

	go this.checkOpenOrders()

	for {
		//当前余额
		var curBalance = this.GetBalance()

		//币USDT成本(每次执行前USDT成本 - 当前USDT余额)
		var costQuote, _ = this.RunGetUsdtCost()
		costQuote = costQuote.Sub(curBalance[this.symbolQuote].Total)

		//币市值(买价*币数量)
		buyprice, err := this.GetBuyPrice()
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		var coinMarket = buyprice.Mul(curBalance[this.symbolBase].Total)

		//计算下单次数
		var baseBalance = curBalance[this.symbolBase].Total
		var buyTimes, firstBuyAmount, firstBuyUsdt = this.GetBuyTimes(buyprice, baseBalance)

		//获取下跌多少在买入
		var buyPoint2, profitTargetPoint2 = this.GetBuyPoint(buyTimes, last_buy_time)

		//首单涨跌比
		var firstBuyPrice = this.GetFirstBuyPrice()
		var firstCurRate decimal.Decimal
		if firstBuyPrice.Cmp(decimal.Zero) > 0 {
			firstCurRate = buyprice.Sub(firstBuyPrice).DivRound(firstBuyPrice, 4).Mul(decimal.NewFromInt(100))
		}

		//收益率
		var earningRate decimal.Decimal
		if costQuote.Cmp(decimal.Zero) > 0 && curBalance[this.symbolBase].Total.Cmp(decimal.Zero) > 0 {
			earningRate = coinMarket.DivRound(costQuote, 4).Sub(decimal.NewFromFloat(1)).Mul(decimal.NewFromInt(100))

			this.log(buyprice, earningRate, costQuote, coinMarket, last_buy_time, buyTimes, firstBuyPrice, firstCurRate, curBalance)
		}

		if costQuote.Cmp(decimal.Zero) == 0 || costQuote.Cmp(firstBuyUsdt.Div(decimal.NewFromInt(2))) < 0 || firstCurRate.Cmp(decimal.NewFromFloat(-buyPoint2)) <= 0 {
			//触发买
			//没有币开始买 或者首次成交不足一半 或者 下跌0.5%买一次

			//10分钟内最多买一次,避免集中买入
			if time.Now().Unix()-last_buy_time < 5*60 {
				time.Sleep(time.Second * 5)
				continue
			}

			////arbusdt 临时添加
			//if buyprice.Cmp(decimal.NewFromFloat(1.03)) > 0 {
			//	time.Sleep(time.Second * 5)
			//	continue
			//}

			var first bool
			var buyAmount decimal.Decimal
			if costQuote.Cmp(firstBuyUsdt.Div(decimal.NewFromInt(2))) < 0 { //不到首次够买额的一半，说明之前没有成交完
				first = true
				buyAmount = firstBuyAmount
			} else {
				buyAmount = curBalance[this.symbolBase].Total.Mul(decimal.NewFromFloat(1))
				//余额buyAmount*sellprice > usdt余额时，所有usdt买入
				if buyAmount.Mul(buyprice).Cmp(curBalance[this.symbolQuote].Total) > 0 {
					buyAmount = curBalance[this.symbolQuote].Total.Div(buyprice).Floor()
				}
			}

			orderAmount, orderid, err := this.Buy(buyAmount, buyprice)
			if err != nil {
				time.Sleep(time.Second * 5)
				continue
			}

			if first {
				//更新首次买入价格
				this.SetFirstBuyPrice(buyprice.String())
				//更新minprice
				this.SetMinPrice(buyprice.String())
			}

			//检测是否成交，超时取消
			var check_time = time.Now().Unix()
			for {
				orderstatus, err := this.p.GetOrderStatus(this.symbol, orderid)
				if err != nil {
					time.Sleep(time.Second * 3)
					continue
				}
				if orderstatus == true {
					//记录购买时间
					last_buy_time = time.Now().Unix()
					//保存记录购买时间
					this.SetLastBuyTime(last_buy_time)

					if costQuote.Cmp(decimal.Zero) > 0 && curBalance[this.symbolBase].Total.Cmp(decimal.Zero) > 0 {
						var curBalance = this.GetBalance()
						var baseBalance = curBalance[this.symbolBase].Total
						var buyTimes, _, _ = this.GetBuyTimes(buyprice, baseBalance)
						var buyPoint2, profitTargetPoint2 = this.GetBuyPoint(buyTimes, last_buy_time)

						var coinMarket = buyprice.Mul(curBalance[this.symbolBase].Total)
						var costQuote, _ = this.RunGetUsdtCost()
						costQuote = costQuote.Sub(curBalance[this.symbolQuote].Total)
						var earningRate = coinMarket.DivRound(costQuote, 4).Sub(decimal.NewFromFloat(1)).Mul(decimal.NewFromInt(100))

						logmsg := fmt.Sprintf("买入%v: %v %v %v %v %v %v %v\n",
							this.symbolBase,
							fmt.Sprintf("币成本:%v", costQuote.DivRound(decimal.NewFromInt(1), 4)),
							fmt.Sprintf("首单涨跌:%v%%", firstCurRate.String()),
							fmt.Sprintf("首买价:%v ", firstBuyPrice),
							fmt.Sprintf("现价:%v", buyprice.DivRound(decimal.NewFromInt(1), 4).String()),
							fmt.Sprintf("目标价:%v", costQuote.Mul(decimal.NewFromFloat(1+profitTargetPoint2/100)).DivRound(curBalance[this.symbolBase].Total, 4).String()),

							fmt.Sprintf("已买次数:%v,下次买入点:-%.2f%%(%v)", buyTimes, buyPoint2, firstBuyPrice.Mul(decimal.NewFromFloat(1-buyPoint2/100))),
							fmt.Sprintf("收益率:%v%%", earningRate.String()))
						message.PrintLog(api.PrintLog, fmt.Sprintf("已买入_%v", this.symbol), fmt.Sprintf("%v", time.Now().UnixNano()), logmsg)
						message.SendDingTalkRobit(true, "oneplat", "every_profit_"+this.symbol, fmt.Sprintf("%v", time.Now().Unix()), logmsg)
					}
					//首笔买单不挂网格卖单
					if first == false {
						//网格思路：立即挂卖单 todo
						for {
							var sellorderid string
							var err error
							if this.symbol == "BNB/FDUSD" {
								sellorderid, err = this.Sell(orderAmount, buyprice.Mul(decimal.NewFromFloat(1+(profitTargetPoint2+3)/100)))
							} else {
								sellorderid, err = this.Sell(orderAmount, buyprice.Mul(decimal.NewFromFloat(1+(profitTargetPoint2+5)/100)))
							}
							if err == nil {
								//加入列表，实时监测成交提醒
								this.openorders.Store(sellorderid, orderAmount.String())
								break
							} else {
								logmsg := fmt.Sprintf("cellsell %v", err.Error())
								message.SendDingTalkRobit(true, "alert", "cellsell_"+this.symbol, fmt.Sprintf("%v", time.Now().Unix()/3600), logmsg)
							}
							time.Sleep(time.Second * 3)
						}
					}
					break
				} else if time.Now().Unix()-check_time > 5*60 {
					//超时取消买单
					this.p.CancelOrder(this.symbol, orderid)
					break
				}
				time.Sleep(time.Second * 3)
			}
		} else if earningRate.Cmp(decimal.NewFromFloat(profitTargetPoint2)) >= 0 {
			//触发卖

			//卖出流程添加回撤时刻，避免还在上升时过早卖出
			if trailing == false {
				//进入回撤时刻
				trailing = true
				trailingMaxRate = earningRate
				time.Sleep(time.Second * 2)
				continue
			} else {
				if earningRate.Cmp(trailingMaxRate) >= 0 {
					//还在上涨阶段
					trailingMaxRate = earningRate
					time.Sleep(time.Second * 2)
					continue
				} else {
					//开始回落卖掉
					trailing = false
					trailingMaxRate = decimal.Zero
				}
			}

			//网格思路：取消所有挂单 todo
			this.p.CancelOpenOrders(this.symbol)

			//卖
			orderid, err := this.SellAll(curBalance[this.symbolBase].Total, buyprice)
			if err != nil {
				time.Sleep(time.Second * 2)
				continue
			}

			//检测是否成交，超时取消
			var sell_time = time.Now().Unix()
			for {
				time.Sleep(time.Second * 5)
				orderstatus, err := this.p.GetOrderStatus(this.symbol, orderid)
				if err != nil {
					continue
				}
				if orderstatus == true {
					var newBalance = this.GetBalance()
					if newBalance[this.symbolQuote].Total.Cmp(curBalance[this.symbolQuote].Total) != 0 {
						runPreCostQuote, _ := this.RunGetUsdtCost()
						//套利完成,更新usdt成本
						this.RunSetUsdtCost(newBalance[this.symbolQuote].Total.String(), newBalance[this.symbolBase].Total.String())
						//更新首次买入价格
						this.SetFirstBuyPrice("0")
						//更新minprice
						this.SetMinPrice("0")

						//机器人通知
						this.pushNotice(run_startTime, buyprice,
							every_trading_startTime,
							newBalance[this.symbolQuote].Total.Sub(runPreCostQuote),
							run_bBalance,
							newBalance, coinMarket)

						//等5分钟
						time.Sleep(time.Second * 5)

						every_trading_startTime = time.Now().Unix()
						break
					}
				} else if time.Now().Unix()-sell_time > 5*60 {
					//超时取消卖单
					this.p.CancelOrder(this.symbol, orderid)
					//logmsg := fmt.Sprintf("触发盈利但未成交：%v", false)
					//utils.SendDingTalkRobit(true, "alert", "sellorder_"+this.symbol, fmt.Sprintf("%v", time.Now().Unix()), logmsg)
					break
				}
			}
		}

		time.Sleep(time.Second * 30)
	}
}

func (this *DCAService) log(buyprice decimal.Decimal, earningRate decimal.Decimal, costQuote, coinMarket decimal.Decimal, last_buy_time int64, buyTimes int, firstBuyPrice, firstCurRate decimal.Decimal, curBalance map[string]model.Balance) {
	var buyPoint2, profitTargetPoint2 = this.GetBuyPoint(buyTimes, last_buy_time)

	//首单涨跌比
	var firstMaxRate decimal.Decimal
	if firstBuyPrice.Cmp(decimal.Zero) > 0 {
		//买入后出现的最低价
		var minprice = this.GetMinPrice()
		if minprice.Cmp(decimal.Zero) == 0 || buyprice.Cmp(minprice) < 0 {
			minprice = buyprice
			this.SetMinPrice(buyprice.String())
		}
		//买入后首单最大下跌比
		firstMaxRate = minprice.Sub(firstBuyPrice).DivRound(firstBuyPrice, 4).Mul(decimal.NewFromInt(100))
	}

	//距上次买入已过多久
	pastTime := time.Now().Unix() - last_buy_time
	pastTimeFormatTime, _ := decimal.NewFromFloat(float64(pastTime)/3600.0).DivRound(decimal.NewFromInt(1), 1).Float64()
	pastTimeFormatUnit := "h"
	if pastTimeFormatTime < 1 {
		pastTimeFormatTime, _ = decimal.NewFromFloat(float64(pastTime)/60.0).DivRound(decimal.NewFromInt(1), 0).Float64()
		pastTimeFormatUnit = "m"
	}

	fmt.Println(
		fmt.Sprintf("币成本:%v", costQuote.DivRound(decimal.NewFromInt(1), 4)),
		fmt.Sprintf("币市值:%v", coinMarket.DivRound(decimal.NewFromInt(1), 4)),
		fmt.Sprintf("首单涨跌:%v%%(最大跌幅:%v%%)", firstCurRate.String(), firstMaxRate),
		fmt.Sprintf("首买价:%v ", firstBuyPrice),
		fmt.Sprintf("目标价:%v", costQuote.Mul(decimal.NewFromFloat(1+profitTargetPoint2/100)).DivRound(curBalance[this.symbolBase].Total, 4).String()),
		fmt.Sprintf("现价:%v", buyprice.DivRound(decimal.NewFromInt(1), 4).String()),

		fmt.Sprintf("\n距上次买入:%v,已买次数:%v,新买入点:-%.2f%%", fmt.Sprintf("%v%v", pastTimeFormatTime, pastTimeFormatUnit), buyTimes, buyPoint2),
		fmt.Sprintf("收益率:%v%%", earningRate.String()),

		fmt.Sprintf("\t%v", time.Now().Format("2006-01-02 15:04")))
}

// 计算下单次数
func (this *DCAService) GetBuyTimes(buyprice, baseBalance decimal.Decimal) (int, decimal.Decimal, decimal.Decimal) {
	var firstBuyAmount decimal.Decimal
	var firstBuyUsdt decimal.Decimal
	//if this.symbol == "BNB/FDUSD" {
	//	//首单3500usdt
	//	firstBuyAmount = decimal.NewFromFloat(100*30).DivRound(buyprice, 0) //100u对应币数
	//	firstBuyUsdt = decimal.NewFromFloat(100 * 30)
	//} else {
	//	//首单700usdt
	//	firstBuyAmount = decimal.NewFromFloat(100*5).DivRound(buyprice, 0) //100u对应币数
	//	firstBuyUsdt = decimal.NewFromFloat(100 * 5)
	//}

	firstBuyAmountKey := fmt.Sprintf("%v.firstbuyamount", this.symbol)
	if !viper.IsSet(firstBuyAmountKey) {
		panic(fmt.Sprintf("%v not exist", firstBuyAmountKey))
	}
	firstBuyAmount = decimal.NewFromFloat(viper.GetFloat64(firstBuyAmountKey)).DivRound(buyprice, 0) //100u对应币数
	firstBuyUsdt = decimal.NewFromFloat(viper.GetFloat64(firstBuyAmountKey))
	if firstBuyUsdt.Cmp(decimal.Zero) <= 0 {
		panic(fmt.Sprintf("%v 值无效", firstBuyAmountKey))
	}

	//计算下单次数
	var buytimes int = 1
	for {
		if baseBalance.Cmp(firstBuyAmount.Mul(decimal.NewFromFloat(1.03))) <= 0 {
			break
		} else {
			buytimes++
			baseBalance = baseBalance.Div(decimal.NewFromInt(1 + 1))
		}
	}
	return buytimes, firstBuyAmount, firstBuyUsdt
}

// 根据下单次数获取下次买入的损失比率
func (this *DCAService) GetBuyPoint(buytimes int, last_buy_time int64) (buyPoint float64, profitTargetPoint float64) {
	defer func() {
		if this.symbol == "BNB/FDUSD" {

		} else {
			//超时未买单，降低一半的buypoint值
			if time.Now().Unix()-last_buy_time > 90*60 && buyPoint >= 5 {
				buyPoint = buyPoint / 4 * 3
			}
		}
	}()

	buyPointKey := fmt.Sprintf("%v.buytimes%v.buyPoint", this.symbol, buytimes)
	if !viper.IsSet(buyPointKey) {
		buyPointKey = fmt.Sprintf("%v.default.buyPoint", this.symbol)
	}
	profitTargetPointKey := fmt.Sprintf("%v.buytimes%v.profitTargetPoint", this.symbol, buytimes)
	if !viper.IsSet(profitTargetPointKey) {
		profitTargetPointKey = fmt.Sprintf("%v.default.profitTargetPoint", this.symbol)
	}

	decBuyPoint, err := decimal.NewFromString(viper.GetString(buyPointKey))
	if err != nil {
		panic(fmt.Sprintf("%v value invalid. %v", buyPointKey, err))
	}
	buyPoint, _ = decBuyPoint.Float64()

	decProfitTargetPoint, err := decimal.NewFromString(viper.GetString(profitTargetPointKey))
	if err != nil {
		panic(fmt.Sprintf("%v value invalid. %v", profitTargetPointKey, err))
	}
	profitTargetPoint, _ = decProfitTargetPoint.Float64()

	if buyPoint <= 0 {
		panic(fmt.Sprintf("%v value <= 0.", buyPointKey))
	}
	if profitTargetPoint <= 0 {
		panic(fmt.Sprintf("%v value <= 0.", decProfitTargetPoint))
	}
	return
	//if this.symbol == "BNB/FDUSD" {
	//	switch {
	//	case buytimes == 1:
	//		//下跌买入点，上涨卖出点
	//		return 0.5, 0.6
	//	case buytimes == 2:
	//		return 1.0, 0.6
	//	case buytimes == 3:
	//		return 1.6, 0.6
	//	case buytimes == 4:
	//		return 3.6, 0.6
	//	case buytimes == 5:
	//		return 7.0, 1.2
	//	case buytimes == 6:
	//		return 11.11, 1.6
	//	case buytimes == 7:
	//		return 15.0, 1.6
	//	case buytimes == 8:
	//		return 20.0, 1.0
	//	case buytimes == 9:
	//		return 25.0, 1.0
	//	case buytimes == 10:
	//		return 30.0, 1.0
	//	default:
	//		return 30.0, 1.0
	//	}
	//
	//} else {
	//	//中风险
	//	switch {
	//	case buytimes == 1:
	//		//下跌买入点，上涨卖出点
	//		return 0.8, 1.5
	//	case buytimes == 2:
	//		return 1.6, 1.0
	//	case buytimes == 3:
	//		return 3.2, 1.0
	//	case buytimes == 4:
	//		return 7.4, 1.2
	//	case buytimes == 5:
	//		return 4, 1.5
	//	case buytimes == 6:
	//		return 20.0, 1.5
	//	case buytimes == 7:
	//		return 25.0, 1.5
	//	case buytimes == 8:
	//		return 30.0, 1.5
	//	case buytimes == 9:
	//		return 35.0, 1.5
	//
	//	default:
	//		return 50.0, 1.5
	//	}
	//}
}

// 初始投入USDT成本
func (this *DCAService) InitialUsdtCost() (decimal.Decimal, decimal.Decimal) {
	return this.initInput[this.symbolQuote].Total, this.initInput[this.symbolBase].Total
}

func (this *DCAService) Buy(buyAmount decimal.Decimal, buyPrice decimal.Decimal) (decimal.Decimal, string, error) {
	symboldata := this.getSymbol(this.symbol)
	buyAmount = buyAmount.Mul(decimal.New(1, symboldata.AmountPrecision)).Floor().Div(decimal.New(1, symboldata.AmountPrecision))
	orderid, err := this.p.OrdersPlace(this.symbol, buyPrice.String(), buyAmount.String(), "BUY")
	return buyAmount, orderid, err
}

func (this *DCAService) Sell(haveAmount, minSellPrice decimal.Decimal) (string, error) {
	sellPrice, _ := this.GetSellPrice()
	if sellPrice.Cmp(minSellPrice) < 0 {
		sellPrice = minSellPrice
	}

	//精度处理
	symboldata := this.getSymbol(this.symbol)
	haveAmount = haveAmount.Mul(decimal.New(1, symboldata.AmountPrecision)).Floor().Div(decimal.New(1, symboldata.AmountPrecision))
	sellPrice = sellPrice.DivRound(decimal.NewFromInt(1), symboldata.PricePrecision)

	//挂卖单
	orderid, err := this.p.OrdersPlace(this.symbol, sellPrice.String(), haveAmount.String(), "SELL")
	return orderid, err
}

func (this *DCAService) SellAll(haveAmount, minSellPrice decimal.Decimal) (string, error) {
	sellPrice, _ := this.GetSellPrice()
	if sellPrice.Cmp(minSellPrice) < 0 {
		sellPrice = minSellPrice
	}

	//精度处理
	//fmt.Println(haveAmount, sellPrice)
	symboldata := this.getSymbol(this.symbol)
	haveAmount = haveAmount.Mul(decimal.New(1, symboldata.AmountPrecision)).Floor().Div(decimal.New(1, symboldata.AmountPrecision))
	sellPrice = sellPrice.DivRound(decimal.NewFromInt(1), symboldata.PricePrecision)
	//fmt.Println(haveAmount, sellPrice)

	//先取消现在的卖单
	openorders, _ := this.p.OpenOrders(this.symbol)
	for i := 0; i < len(openorders); i++ {
		open := openorders[i]
		if open.Side == "SELL" {
			_, err := this.p.CancelOrder(this.symbol, fmt.Sprintf("%v", open.OrderID))
			if err != nil {
				i--
				time.Sleep(time.Second * 2)
			}
		}
	}
	if len(openorders) > 0 {
		time.Sleep(time.Second * 2)
	}

	//重新挂卖单
	orderid, err := this.p.OrdersPlace(this.symbol, sellPrice.String(), haveAmount.String(), "SELL")
	return orderid, err
}

func (this *DCAService) profitFunc(sellPrice decimal.Decimal, iBalance map[string]model.Balance, eBalance map[string]model.Balance, base, quote string) (decimal.Decimal, decimal.Decimal) {
	//多出币换算为等值U
	changeBase := eBalance[base].Total.Sub(iBalance[base].Total)
	changeBaseToQuote := changeBase.Mul(sellPrice)
	//当前最新U余额
	nowQuoteBalance := eBalance[quote].Total.Add(changeBaseToQuote)
	//收益
	profit := nowQuoteBalance.Sub(iBalance[quote].Total).DivRound(decimal.NewFromFloat(1), 5)
	return profit, changeBase
}

func (this *DCAService) pushNotice(run_startTime int64, sellPrice decimal.Decimal, trading_startTime int64, profit decimal.Decimal, run_bBalance map[string]model.Balance, bot_eBalance map[string]model.Balance, coinMarket decimal.Decimal) {

	this.totalCount++

	//当天
	if time.Now().Format("2006-01-02") != time.Unix(this.dayStartTime, 0).Format("2006-01-02") {
		this.dayStartTime = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).Unix()
		this.dayTotalCount = 1
		this.dayTotalProfit = profit
	} else {
		this.dayTotalCount += 1
		this.dayTotalProfit = this.dayTotalProfit.Add(profit)
	}
	//总收益
	this.totalProfit = this.totalProfit.Add(profit)

	//运行累计收益
	runProfit, runChangeBase := this.profitFunc(sellPrice, run_bBalance, bot_eBalance, this.symbolBase, this.symbolQuote)
	//帐户累计收益
	accountProfit, accountChangeBase := this.profitFunc(sellPrice, this.initInput, bot_eBalance, this.symbolBase, this.symbolQuote)

	//运行时间
	runTotalUseTime := time.Now().Unix() - run_startTime
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
	estimateDayProfit := this.dayTotalProfit.DivRound(decimal.NewFromInt(this.dayTotalCount), 4).Mul(decimal.NewFromInt(estimateDayCount)).String()

	var group string = "1"
	logmsg := fmt.Sprintf("group#%v   %v/%v(%v,%v) 交易额:%vu   套利:%vu(今日:%vu) 运行累计:%vu(chanB:%v) 帐户累计:%vu(chanB:%v)   已运行:%v {curprice}   预估Day:%vu/%vt usetime:%vs avgtime:%vs count:%v\n",
		group, this.symbolBase, this.symbolQuote, bot_eBalance[this.symbolBase].Total, utils.DivRound(bot_eBalance[this.symbolQuote].Total, 2),
		coinMarket, utils.DivRound(profit, 2), utils.DivRound(this.dayTotalProfit, 2),
		runProfit, runChangeBase, accountProfit.String(), accountChangeBase.String(),
		fmt.Sprintf("%v%v", runTotalFormatTime, runTotalFormatUnit),
		estimateDayProfit, estimateDayCount, time.Now().Unix()-trading_startTime, avgUseTime, this.totalCount)
	message.PrintLog(api.PrintLog, fmt.Sprintf("every_profit_%v,%v", group, this.symbol), fmt.Sprintf("%v", time.Now().UnixNano()), strings.Replace(logmsg, "{curprice}", "", 1))

	logmsg = strings.Replace(logmsg, "{curprice}", fmt.Sprintf("币价:%v", sellPrice), 1)
	logmsg = strings.Replace(logmsg, "套利", "\n套利", 1)
	logmsg = strings.Replace(logmsg, "已运行", "\n\n已运行", 1)
	logmsg = strings.Replace(logmsg, "预估Day", "\n预估Day", 1)
	message.SendDingTalkRobit(true, "oneplat", "every_profit_"+this.symbol, fmt.Sprintf("%v", time.Now().Unix()), logmsg)
}

func (this *DCAService) checkOpenOrders() {
	////所有卖单
	//openorders, _ := this.p.OpenOrders(this.symbol)
	//for i := 0; i < len(openorders); i++ {
	//	open := openorders[i]
	//	if open.Side == "SELL" {
	//		this.openorders.Store(fmt.Sprintf("%v", open.OrderID), open.OrigQuantity)
	//	}
	//}

	var today = time.Now().Format("2006-01-02")

	for {
		//var suss int
		//var amounts decimal.Decimal
		//this.openorders.Range(func(key, value any) bool {
		//	orderstatus, _ := this.p.GetOrderStatus(this.symbol, key.(string))
		//	if orderstatus == true {
		//		suss++
		//		tval, _ := decimal.NewFromString(value.(string))
		//		amounts = amounts.Add(tval)
		//		this.openorders.Delete(key)
		//	}
		//	time.Sleep(time.Second * 5)
		//	return true
		//})
		//if suss > 0 {
		//	//buyprice, err := this.GetBuyPrice()
		//	//if err == nil {
		//	//	logmsg := fmt.Sprintf("网格交易成交数:%v，成交额:%v", suss, amounts.Mul(buyprice).String())
		//	//	utils.SendDingTalkRobit(true, "alert", "openorders_"+this.symbol, fmt.Sprintf("%v", time.Now().Unix()), logmsg)
		//	//} else {
		//	//	logmsg := fmt.Sprintf("网格交易成交数:%v，成交量:%v", suss, amounts.String())
		//	//	utils.SendDingTalkRobit(true, "alert", "openorders_"+this.symbol, fmt.Sprintf("%v", time.Now().Unix()), logmsg)
		//	//}
		//}

		//日报
		if time.Now().Format("2006-01-02") != today {
			today = time.Now().Format("2006-01-02")

			//当前余额
			var curBalance = this.GetBalance()
			//币市值(买价*币数量)
			buyprice, err := this.GetBuyPrice()
			if err == nil {
				var coinMarket = buyprice.Mul(curBalance[this.symbolBase].Total)
				var balance = curBalance[this.symbolQuote].Total.Add(coinMarket).String()
				var logmsg = fmt.Sprintf("余额市值：%v  --------------------<日报>--------------------", balance)
				message.SendDingTalkRobit(true, "alert", "balance_"+this.symbol, fmt.Sprintf("%v", time.Now().Unix()), logmsg)
			}
		}

		time.Sleep(time.Second * 10)
	}
}

// 交易对信息
func (this *DCAService) getSymbol(symbol string) model.SymbolsData {
	if val, ok := this.symbols[symbol]; ok {
		return val
	} else {
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
			this.symbols[symbol] = fdSymbol
			return fdSymbol
		} else {
			panic(fmt.Sprintf("not find symbol %v", symbol))
		}
	}
}

func (this *DCAService) GetBalance() map[string]model.Balance {
	for {
		balance, err := this.p.GetAccountBalance(this.symbol, true)
		if err != nil {
			time.Sleep(time.Second * 5)
		} else {
			return balance
		}
	}
}

func (this *DCAService) GetBuyPrice() (decimal.Decimal, error) {
	p1, err := this.p.GetMarketDepth(this.symbol)
	if err != nil {
		return decimal.Zero, err
	}
	if this.symbol == "-BNB/FDUSD-" {
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

func (this *DCAService) GetSellPrice() (decimal.Decimal, error) {
	p1, err := this.p.GetMarketDepth(this.symbol)
	if err != nil {
		return decimal.Zero, err
	}

	if this.symbol == "-BNB/FDUSD-" {
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
