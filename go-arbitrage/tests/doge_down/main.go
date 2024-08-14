package main

import (
	"components/database/redis2"
	"components/log/log"
	"components/message"
	"components/myconfig"
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"time"

	"github.com/adshao/go-binance/v2"
)

var (
	symbol             = "DOGEFDUSD"
	firstQty           = 9000.0 //首笔挂单量
	firstPriceIncrease = 1.004  //首笔价格增长比例

	baseQty      = 4500.0 //Safety order size
	baseIncrease = 1.814  //Safety order大小倍数

	priceIncrease = 1.0064 //每笔订单间隔比例(from init price)
	priceFactor   = 1.404  //Safety order间隔倍数
	profitFactor  = 0.0035 //Target profit

	maxSellOrders    = 10 //最大订单数
	actionSellOrders = 2  //活跃订单数
	apiKey           = "mCXfycRaEiffizOajnB1VsVxytyUFnaA1tK4eX8QyuM8G565Weq5s4QXoyhkzwdE"
	secretKey        = "wvRdYxo9O4IeBywbDCZgGhflwDwv2ERUbdQHUgoZ8JXTpUDGvFsTnXtzQOHxL9XW"
	initSellQty      = make([]decimal.Decimal, maxSellOrders+1)
	initSellPrice    = make([]decimal.Decimal, maxSellOrders+1)
	initRadios       = make([]decimal.Decimal, maxSellOrders+1)
	baseDoubleIndex  = 0 //前几手可以补单，之后不补单，否则很快占用大量资金
)

var client *binance.Client
var buyOrderId int64      //当前买单orderId
var buySuccLastTime int64 //最近一次买单成功执行时间，挂买单时，如果发现很久没有下，需要进行回撤处理
var stop bool             //暂停买卖挂单
var sellOrders int        //已挂卖单数

var placeBuyLastTime int64  //最近一次挂买单时间，定时检测，超过没有发生更新买单，重新计算下单
var placeSellLastTime int64 //最近一次挂卖单时间，定时检测，如果发现很久时没有挂卖了，这时可能在下跌中，重新挂卖单列表
var startTime = time.Now().Unix()

func main() {
	log.InitLogger("./", "testDoge", true)

	client = binance.NewClient(apiKey, secretKey)

	err := initSellOrders(false)
	if err != nil {
		panic(err)
	}
	//return

	initialUSDT, initialDOGE, _ := RunGetDogeCost(symbol)
	log.Logger.Debugf("Initial balances: %s DOGE, %s FDUSD", initialDOGE.String(), initialUSDT.String())
	//初始投入值
	RunGetDogeCost(symbol + "-INIT")

	//重启时取消所有挂单
	openOrders, err := openOrders()
	if err != nil {
		panic(err)
	}
	cancelOrders("all", openOrders)
	//重置每轮时间，重新挂单
	RunSetEachRoundTime(symbol, time.Now().Unix())

	go checkFinish()

	for {
		if checkFee() {
			haveNewSell, openSells, openBuys, minSellPrice := placeSells()
			time.Sleep(time.Second * 10)
			placeBuy(haveNewSell, openSells, openBuys, minSellPrice)
			time.Sleep(time.Second * 60)
		}
	}
}

func checkFee() bool {
	fees, err := client.NewTradeFeeService().Symbol(symbol).Do(context.Background())
	if err != nil {
		log.Logger.Errorf("NewTradeFeeService %v", err)
		time.Sleep(time.Minute)
		return false
	}
	makerfee, err := decimal.NewFromString(fees[0].MakerCommission) //挂单
	if makerfee.Cmp(decimal.Zero) > 0 {
		log.Logger.Errorf("makerfee:%v >0 ", fees[0].MakerCommission)
		message.SendDingTalkRobit(true, "oneplat", "doge2_every_fee_"+symbol, fmt.Sprintf("%v", time.Now().Unix()/3600), "makerfee:%v >0")
		time.Sleep(time.Minute)
		return false
	}
	return true
}

func initSellOrders(init bool) error {
	if init {
		sellOrders = 0
		RunSetInitPrice(symbol, 0)
		RunSetEachRoundTime(symbol, time.Now().Unix())
	}

	//每次重启服务后，仍然根据上次的价格创建initSellPrice数组
	price, err := RunGetInitPrice(symbol)
	if err != nil {
		log.Logger.Error(err)
		return err
	}
	RunSetInitPrice(symbol, price)

	//每轮的开始时间
	eachRoundTime := RunGetEachRoundTime(symbol)
	RunSetEachRoundTime(symbol, eachRoundTime)

	log.Logger.Debugf("[initSellOrders] Current price: %f", price)

	initSellQty[0] = decimal.NewFromFloat(firstQty).Truncate(0)
	initRadios[0] = decimal.NewFromFloat(firstPriceIncrease - 1)
	initSellPrice[0] = decimal.NewFromFloat(price * math.Pow(firstPriceIncrease, 1)).Truncate(5)
	for i := 0; i < maxSellOrders; i++ {
		qty := decimal.NewFromFloat(baseQty * math.Pow(baseIncrease, float64(i))).Truncate(0)
		initSellQty[i+1] = qty

		if i == 0 {
			initRadios[i+1] = decimal.NewFromFloat(priceIncrease - 1)
			initSellPrice[i+1] = decimal.NewFromFloat(price * math.Pow(priceIncrease, 1)).Truncate(5)
		} else {
			initRadios[i+1] = (decimal.NewFromFloat(priceIncrease - 1)).Add(initRadios[i].Mul(decimal.NewFromFloat(priceFactor)))
			initSellPrice[i+1] = initRadios[i+1].Add(decimal.NewFromInt(1)).Mul(decimal.NewFromFloat(price)).Truncate(5)
		}
	}

	var total = decimal.Zero
	var totalU = decimal.Zero
	for i, qty := range initSellQty {
		sellPrice := initSellPrice[i]
		total = total.Add(qty)
		totalU = totalU.Add(qty.Mul(sellPrice))
		profitPrice := totalU.Div(total.Mul(decimal.NewFromFloat(1 + profitFactor))).Truncate(5)
		if i < 50 || i > len(initSellQty)-20 {
			fmt.Println(fmt.Sprintf("%v: ", i),
				fmt.Sprintf("大小:%v 累计:%v", qty.String(), total.String()),
				fmt.Sprintf("卖价:%v(涨幅:%v)", sellPrice.String(), sellPrice.Sub(decimal.NewFromFloat(price)).DivRound(decimal.NewFromFloat(price), 5).Mul(decimal.NewFromFloat(100)).String()+"%"),
				fmt.Sprintf("止盈价:%v(下跌:%v) 累计U:%v", profitPrice.String(), sellPrice.Sub(profitPrice).DivRound(sellPrice, 5).Mul(decimal.NewFromFloat(100)).String()+"%", totalU.String()))
		}
	}
	return nil
}

// 挂卖
// 否：如果有成交卖单，定时会把成交卖单在补上 --> 容易占用大量资金，一直上涨亏损过大
func placeSells() (bool, int, int, decimal.Decimal) {
	if stop == true {
		return false, -1, -1, decimal.Zero
	}

	openOrders, err := openOrders()
	if err != nil {
		log.Logger.Error("[placeSells] Error openOrders:", err)
		return false, -1, -1, decimal.Zero
	}
	//统计卖挂单数
	var openSells int
	var openBuys int
	var minSellPrice decimal.Decimal
	for _, order := range openOrders {
		var price, _ = decimal.NewFromString(order.Price)
		if order.Side == binance.SideTypeSell {
			openSells++
		}
		if order.Side == binance.SideTypeBuy {
			openBuys++
		}
		//当前挂单中最小的卖单价
		if order.Side == binance.SideTypeSell && (minSellPrice == decimal.Zero || price.Cmp(minSellPrice) < 0) {
			minSellPrice = price
		}
	}

	var haveNewSell bool
	//卖挂单数少于活跃订单数
	if openSells < actionSellOrders {
		var redis = redis2.NewRedisCli(myconfig.GConfig.Redis.Host, myconfig.GConfig.Redis.Password, myconfig.GConfig.Redis.DB)

		curPrice, err := getCurrentPrice()
		if err != nil {
			log.Logger.Error("[placeSells] Error getCurrentPrice:", err)
			return false, openSells, openBuys, minSellPrice
		}

		for i, sellPrice := range initSellPrice {
			if openSells < actionSellOrders {
				if sellPrice.Cmp(decimal.NewFromFloat(curPrice)) > 0 {
					//是否存在相同的价格
					var sameprice bool
					for _, order := range openOrders {
						orderPrice, _ := decimal.NewFromString(order.Price)
						if order.Side == binance.SideTypeSell && sellPrice.Cmp(orderPrice) == 0 {
							sameprice = true
							break
						}
					}
					if sameprice == false {
						//是否存在相同的价格
						key := fmt.Sprintf("sameSellPrice#%v#%v#%v", symbol, RunGetEachRoundTime(symbol), sellPrice.String())
						val, err := redis.GetString(key)
						if err != nil {
							log.Logger.Error(err)
							return false, openSells, openBuys, minSellPrice
						}
						//不存在进行补单和前几手可以重复挂单
						if val == "" || i < baseDoubleIndex {
							openSells++  //局部变量
							sellOrders++ //全局变量
							haveNewSell = true
							placeSellLastTime = time.Now().Unix()
							if _, err := placeOrder("SELL", initSellQty[i].String(), sellPrice.String()); err == nil {
								redis.SetEX(key, "1", 7*24*3600*time.Second)
								//当前挂单中最小的卖单价
								if sellPrice.Cmp(minSellPrice) < 0 {
									minSellPrice = sellPrice
								}
							}
						}
					}
				}
			}
		}
	}

	//超时没有成交了
	var timeout = (time.Now().Unix() - placeSellLastTime) / (3600 * 2)
	if timeout >= 1 {
		message.SendDingTalkRobit(true, "oneplat", "doge2_every_sell_"+symbol, fmt.Sprintf("%v", time.Now().Unix()/(3600*2)), fmt.Sprintf("超过%v小时没有新卖单", timeout*2))
	}
	return haveNewSell, openSells, openBuys, minSellPrice
}

func placeBuy(haveNewSell bool, openSells, openBuys int, minSellPrice decimal.Decimal) {
	if stop == true {
		return
	}
	//没有产生新卖单情况
	if haveNewSell == false {
		//placeSells()发生了异常
		if openSells == -1 && openBuys == -1 {
			return
		}
		//当前有买单:不需要重新挂买单(安全点加个超时吧）
		if openSells > 0 && openBuys > 0 && time.Now().Unix()-placeBuyLastTime < 60*60 {
			return
		}
		//当前没有买单:可能发生了取消买单成功但重新下买单失败，所以这里有个超过n分钟允许重新下单
		if (openSells > 0 && openBuys == 0) && time.Now().Unix()-placeBuyLastTime < 5*60 {
			return
		}
	}

	runInitialUSDT, runInitialDOGE, _ := RunGetDogeCost(symbol)
	log.Logger.Debugf("[placeBuy] Initial balances: %s DOGE, %s FDUSD", runInitialDOGE.String(), runInitialUSDT.String())

	currentDOGE, currentUSDT, err := getBalances()
	if err != nil {
		log.Logger.Error("[placeBuy] Error getBalances:", err)
		return
	}
	log.Logger.Debugf("[placeBuy] current balances: %s DOGE, %s FDUSD", currentDOGE.String(), currentUSDT.String())

	//计算本轮所有卖单获得U与当前U余额比较，取小值作为本次的购买本金
	var calcUsdtDelta = func(usdtBalance decimal.Decimal, minSellPrice decimal.Decimal) (float64, bool) {
		var total decimal.Decimal
		for i, sellPrice := range initSellPrice {
			if sellPrice.Cmp(minSellPrice) < 0 {
				total = total.Add(sellPrice.Mul(initSellQty[i]))
			}
		}
		if total.Cmp(usdtBalance) > 0 {
			return usdtBalance.Float64()
		} else {
			return total.Float64()
		}
	}

	dogeDelta, _ := currentDOGE.Sub(runInitialDOGE).Float64()
	//usdtDelta, _ := currentUSDT.Sub(runInitialUSDT).Float64()
	usdtDelta, _ := calcUsdtDelta(currentUSDT, minSellPrice)
	log.Logger.Debugf("[placeBuy] dogeDelta: %v, usdtDelta: %v buyOrderId: %v", dogeDelta, usdtDelta, buyOrderId)
	//doge为负，表示已有卖单成交，开始挂买单
	if dogeDelta < 0 {
		if usdtDelta <= 0 {
			logmsg := "异常:套利还未执行完，U的余额增量居然小于等于0"
			message.SendDingTalkRobit(true, "oneplat", "doge2_every_profit_"+symbol, fmt.Sprintf("%v", time.Now().Unix()/5*60), logmsg)
		}

		//可能一直上涨没有大的回调，这时需要把之前的收益拿出来，减少本次的买回量(doge),确保可以成交（收益回撤了）
		_, firstInitialDOGE, _ := RunGetDogeCost(symbol + "-INIT")
		totalProfitDoge, _ := runInitialDOGE.Sub(firstInitialDOGE).Float64()
		if buySuccLastTime > 0 && time.Now().Unix()-buySuccLastTime > 24*3600 {
			dogeDelta = dogeDelta + totalProfitDoge/3
		} else if buySuccLastTime > 0 && time.Now().Unix()-buySuccLastTime > 12*3600 {
			dogeDelta = dogeDelta + totalProfitDoge/4
		} else if buySuccLastTime > 0 && time.Now().Unix()-buySuccLastTime > 6*3600 {
			dogeDelta = dogeDelta + totalProfitDoge/5
		}

		//要买回的币
		dogeToBuyBack := decimal.NewFromFloat((0-dogeDelta)*(1+profitFactor)).DivRound(decimal.NewFromInt(1), 0)
		//购买价
		newBuyPrice := decimal.NewFromFloat(usdtDelta).Div(dogeToBuyBack).Truncate(5)

		curPrice, err := getCurrentPrice()
		if err != nil {
			log.Logger.Error("[placeBuy] Error fetching current price:", err)
			return
		}
		//新购买价高于当前价
		if newBuyPrice.Cmp(decimal.NewFromFloat(curPrice)) > 0 {
			newBuyPrice = decimal.NewFromFloat(curPrice).Mul(decimal.NewFromFloat(1 - 0.0005)).Truncate(5)
		}

		//取消买单
		if buyOrderId > 0 {
			cancelOrder(buyOrderId)
			time.Sleep(time.Millisecond * 500)
		}

		//计算新买价和数量
		var needUsdt = dogeToBuyBack.Mul(newBuyPrice)
		var needDownRatio = decimal.NewFromFloat(curPrice).Sub(newBuyPrice).Div(decimal.NewFromFloat(curPrice)).Truncate(4).Mul(decimal.NewFromFloat(100)).String() + "%"
		log.Logger.Debugf("[placeBuy] dogeToBuyBack: %v, newBuyPrice: %v curPrice: %v(needDownRatio:%v) (needUsdt: %v, usdtBalance: %v)", dogeToBuyBack, newBuyPrice, curPrice, needDownRatio, needUsdt, currentUSDT)

		orderId, err := placeOrder("BUY", dogeToBuyBack.String(), newBuyPrice.String())
		if err != nil {
			log.Logger.Error("[placeBuy] Error placeOrder:", err)
			message.SendDingTalkRobit(true, "oneplat", "doge2_every_buy_"+symbol, fmt.Sprintf("%v", time.Now().Unix()/10*60), err.Error())
		} else {
			buyOrderId = orderId
			placeBuyLastTime = time.Now().Unix()
		}
	}
}

func checkFinish() {
	var profitTimes int
	for {
		time.Sleep(time.Second * 10)
		openOrders, err := openOrders()
		if err != nil {
			stop = false
			log.Logger.Error(err)
			continue
		}

		if buyOrderId > 0 {
			orderstatus, qty, err := getOrderStatus(symbol, buyOrderId)
			if err != nil {
				log.Logger.Errorf("[checkFinish] getOrderStatus", err)
				continue
			}
			if orderstatus == true {
				//stop后有可能placeSells或placeBuy还在执行，这里先sleep会
				stop = true
				time.Sleep(time.Second * 2)

				currentDOGE, currentUSDT, err := getBalances()
				if err != nil {
					stop = false
					log.Logger.Errorf("getBalances", err)
					continue
				}

				var openSells int
				for _, order := range openOrders {
					if order.Side == binance.SideTypeSell {
						openSells++
					}
				}

				//最近一次买成功时间
				buySuccLastTime = time.Now().Unix()

				//套利通知
				{
					profitTimes++

					_, runDOGE, _ := RunGetDogeCost(symbol)
					initUSDT, initDOGE, initTime := RunGetDogeCost(symbol + "-INIT")

					dogeDelta, _ := currentDOGE.Sub(runDOGE).Float64()
					//说明又有卖单成交了，这次套利还要继续
					if dogeDelta < 0 {
						stop = false
						log.Logger.Errorf("发生了买单已成交，但关闭前又有卖单成交，继续交易...")
						message.SendDingTalkRobit(true, "oneplat", "doge2_every_continue_"+symbol, fmt.Sprintf("%v", time.Now().Unix()/3600), "发生了买单已成交，但关闭前又有卖单成交，继续交易...")
						continue
					}

					totalDogeDelta, _ := currentDOGE.Sub(initDOGE).Float64()
					totalUsdtDelta, _ := currentUSDT.Sub(initUSDT).Float64()

					price, _ := getCurrentPrice()

					logmsg := fmt.Sprintf("E... %v 交易量:%v 卖单成交数:%v 套利:%vdoge 总套利:%vdoge(chanU:%v) \n套利次数:%v 余额:%vdoge(chanU:%v) 当前价格:%v",
						symbol, qty, sellOrders-openSells, dogeDelta, totalDogeDelta, totalUsdtDelta,
						profitTimes, currentDOGE.String(), currentUSDT.String(), price)
					log.Logger.Debugf("[checkFinish] profit: %v", logmsg)
					message.SendDingTalkRobit(true, "oneplat", "doge2_every_profit_"+symbol, fmt.Sprintf("%v", time.Now().Unix()), logmsg)

					//保存当前余额,重置初始投入值
					dogeBalanceSaveFile(initTime, currentUSDT.String(), currentDOGE.String())

					//每天余额保存到redis,后面做邮件报表使用
					dogeBalanceSaveRedis(currentDOGE.String())
				}

				time.Sleep(time.Second * 5)

				//重新铺单
				cancelOrders(binance.SideTypeSell, openOrders)
				initSellOrders(true)
				buyOrderId = 0
				stop = false

				continue
			} else {
				////运行中卖单长时间没有成交
				//if buyOrderId > 0 && placeSellLastTime > 0 && time.Now().Unix()-placeSellLastTime > 30*60 {
				//	stop = true
				//	cancelOrders(binance.SideTypeSell, openOrders)
				//	initSellOrders(false)
				//	placeSellLastTime = time.Now().Unix()
				//	stop = false
				//}
			}
		} else {
			//本轮长时间没成交
			if buyOrderId == 0 && placeSellLastTime > 0 && time.Now().Unix()-placeSellLastTime > 10*60 {
				curPrice, err := getCurrentPrice()
				if err != nil {
					log.Logger.Error("[placeSells] Error getCurrentPrice:", err)
					continue
				}
				initPrice, err := RunGetInitPrice(symbol)
				if err != nil {
					log.Logger.Error(err)
					continue
				}
				//当前价格小于当时铺单价格，说明是下跌导致的未成交需要重新铺单(如果当前价格比当时铺单价格还高仍未成交，说明价格太稳定不需要重新铺单）
				if curPrice < initPrice {
					stop = true
					cancelOrders(binance.SideTypeSell, openOrders)
					initSellOrders(true)
					placeSellLastTime = time.Now().Unix()
					stop = false
				}
			}
		}
	}
}

func dogeBalanceSaveFile(initTime int64, currentUSDT, currentDOGE string) {
	////todo usdt始终按0算
	//currentUSDT = "0"

	//保存当前余额
	log.Logger.Debugf("[checkFinish] Initial balances: %s DOGE, %s FDUSD", currentDOGE, currentUSDT)
	RunSetDogeCost(symbol, currentUSDT, currentDOGE)

	//每24小时结算一次：重置初始投入值，在回撤计算时最多回撤24小时收益
	if (time.Now().Unix()-initTime)/(24*3600) >= 1 {
		RunSetDogeCost(symbol+"-INIT", currentUSDT, currentDOGE)
	}
}

// 每天余额保存到redis,后面做邮件报表使用
func dogeBalanceSaveRedis(currentDOGE string) {
	var redis = redis2.NewRedisCli("localhost:6379", "", 0)
	var key = fmt.Sprintf("dogedown-%v", time.Now().Format("2006-01-02"))
	redis.SetEX(key, currentDOGE, 60*24*3600*time.Second)
}
