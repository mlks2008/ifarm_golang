package main

import (
	"components/database/redis2"
	"components/log/log"
	"components/message"
	"components/myconfig"
	"flag"
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"os"
	"time"

	"github.com/adshao/go-binance/v2"
)

var (
	qtyPrecision       int32 = 2
	pricePrecision     int32 = 3
	symbol                   = "FILUSDT"
	firstQty                 = 70.0  //首笔挂单量
	firstPriceIncrease       = 1.004 //首笔价格增长比例

	baseQty      = 39.0 //Safety order size
	baseIncrease = 2.0  //Safety order大小倍数

	priceIncrease = 1.006 //每笔订单间隔比例(from init price)
	priceFactor   = 1.4   //Safety order间隔倍数
	profitFactor  = 0.004 //Target profit

	minAllowPrice    = 0.1 //最小允许价格（小于此值不在执行）
	maxSellOrders    = 8   //最大订单数
	actionSellOrders = 3   //活跃订单数
	// mainapi
	//apiKey           = "mCXfycRaEiffizOajnB1VsVxytyUFnaA1tK4eX8QyuM8G565Weq5s4QXoyhkzwdE"
	//secretKey        = "wvRdYxo9O4IeBywbDCZgGhflwDwv2ERUbdQHUgoZ8JXTpUDGvFsTnXtzQOHxL9XW"
	// oneplat
	apiKey          = "3JiMItY7JeQoxWNAlylhsxCI38hysP5OZUgypWewm3PhKUaVx9pMv3dTUlyT5sbS"
	secretKey       = "iPP0IRusNqUhKtVyl0gSteRnTEpMUXttUWSekH2MeqljcLkfzwyJ6J8nmUyUOxhn"
	initSellQty     = make([]decimal.Decimal, maxSellOrders+1)
	initSellPrice   = make([]decimal.Decimal, maxSellOrders+1)
	initRadios      = make([]decimal.Decimal, maxSellOrders+1)
	baseDoubleIndex = 0 //前几手可以补单，之后不补单，否则很快占用大量资金
)

var (
	gitHash                 string
	Filed_BuyOrderId        string = "buyorderid"
	Filed_EachRoundTime            = "eachroundtime" //每轮值不同,redisKey的一部分：用于到redis中检测该价格是否已经挂过卖单
	Filed_PlaceSellLastTime string = "placeSellLastTime"
)

var client *binance.Client

var stopByBalance bool //CRV:0时暂停服务，其它值恢复
var stop bool          //暂停买卖挂单

var buyOrderId int64        //当前买单orderId
var returnProfitFil int64   //回撤数
var placeBuyLastTime int64  //最近一次挂买单时间，定时检测，超过没有发生更新买单，重新计算下单
var buySuccLastTime int64   //最近一次买单成功执行时间，挂买单时，如果发现很久没有成交，需要进行回撤处理
var placeSellLastTime int64 //最近一次挂卖单时间，定时检测，如果发现很久时没有挂卖了，这时可能在下跌中，重新挂卖单列表
var startTime = time.Now().Unix()

var (
	robot string
)

func init() {
	flag.StringVar(&robot, "robot", "", "eg: -robot oneplat/mainapi")
	flag.Float64Var(&minAllowPrice, "minAllowPrice", minAllowPrice, "")
	flag.IntVar(&maxSellOrders, "maxSellOrders", maxSellOrders, "")
	flag.Parse()

	initSellQty = make([]decimal.Decimal, maxSellOrders+1)
	initSellPrice = make([]decimal.Decimal, maxSellOrders+1)
	initRadios = make([]decimal.Decimal, maxSellOrders+1)
}

func main() {
	if robot == "oneplat" {
		// oneplat
		apiKey = "3JiMItY7JeQoxWNAlylhsxCI38hysP5OZUgypWewm3PhKUaVx9pMv3dTUlyT5sbS"
		secretKey = "iPP0IRusNqUhKtVyl0gSteRnTEpMUXttUWSekH2MeqljcLkfzwyJ6J8nmUyUOxhn"
	} else if robot == "mainapi" {
		// mainapi
		apiKey = "mCXfycRaEiffizOajnB1VsVxytyUFnaA1tK4eX8QyuM8G565Weq5s4QXoyhkzwdE"
		secretKey = "wvRdYxo9O4IeBywbDCZgGhflwDwv2ERUbdQHUgoZ8JXTpUDGvFsTnXtzQOHxL9XW"
	} else {
		panic("apikey not exist")
	}

	log.InitLogger("./", "testFil", true)

	client = binance.NewClient(apiKey, secretKey)

	err := initSellOrders(false)
	if err != nil {
		panic(err)
	}
	if os.Getenv("local") == "true" {
		return
	}

	//加载买入buyorderid
	buyOrderId = RunGetInt64(robot, symbol, Filed_BuyOrderId)
	//加载最近卖出时间
	placeSellLastTime = RunGetInt64(robot, symbol, Filed_PlaceSellLastTime)
	log.Logger.Debugf("Load buyOrderId: %v, placeSellLastTime: %v", buyOrderId, placeSellLastTime)

	initialUSDT, initialFIL, _ := RunGetFilCost(robot, symbol)
	log.Logger.Debugf("Initial balances: %s FIL, %s FDUSD", initialFIL.String(), initialUSDT.String())
	RunGetFilCost(robot, symbol+"-INIT")

	currentFIL, currentUSDT, stopBalance, err := getBalances()
	if err != nil {
		panic(err)
	}
	checkStopByBalance(currentUSDT.String(), currentFIL.String(), stopBalance)

	log.Logger.Debugf("Robot:%v, minAllowPrice:%v, maxSellOrders:%v, gitHash:%v, stopByBalance:%v", robot, minAllowPrice, maxSellOrders, gitHash, stopByBalance)

	go checkFinish()

	for {
		if checkFee() {
			//每轮开始前检测是否已经暂停
			currentFIL, currentUSDT, stopBalance, err := getBalances()
			if err != nil {
				log.Logger.Errorf("getBalances", err)
				time.Sleep(time.Second * 20)
				continue
			}
			checkStopByBalance(currentUSDT.String(), currentFIL.String(), stopBalance)

			haveNewSell, openSells, openBuys := placeSells()
			time.Sleep(time.Second * 20)
			placeBuy(haveNewSell, openSells, openBuys)
			time.Sleep(time.Second * 20)
		}
	}
}

func checkFee() bool {
	//fees, err := client.NewTradeFeeService().Symbol(symbol).Do(context.Background())
	//if err != nil {
	//	log.Logger.Errorf("NewTradeFeeService %v", err)
	//	time.Sleep(time.Minute)
	//	return false
	//}
	//makerfee, err := decimal.NewFromString(fees[0].MakerCommission) //挂单
	//if makerfee.Cmp(decimal.Zero) > 0 {
	//	log.Logger.Errorf("makerfee:%v >0 ", fees[0].MakerCommission)
	//	message.SendDingTalkRobit(true, robot, "fil2_every_fee_"+symbol, fmt.Sprintf("%v", time.Now().Unix()/3600), "makerfee:%v >0")
	//	time.Sleep(time.Minute)
	//	return false
	//}
	return true
}

func checkFinish() {
	var profitTimes int
	for {
		time.Sleep(time.Second * 5)
		openOrders, err := openOrders()
		if err != nil {
			stop = false
			log.Logger.Error(err)
			continue
		}

		if buyOrderId > 0 {
			orderstatus, qty, err := getOrderStatus(symbol, buyOrderId)
			if err != nil {
				log.Logger.Errorf("[checkFinish] getOrderStatus %v %v", buyOrderId, err)
				continue
			}
			if orderstatus == true {
				stop = true
				time.Sleep(time.Second)
				//重新最新余额
				currentFIL, currentUSDT, _, err := getBalances()
				if err != nil {
					continue
				}

				//最近一次买成功时间
				buySuccLastTime = time.Now().Unix()

				//套利通知
				{
					//计算本次套利
					runUSDT, runFIL, _ := RunGetFilCost(robot, symbol)
					filDelta, _ := currentFIL.Sub(runFIL).Float64()
					usdtDelta, _ := currentUSDT.Sub(runUSDT).Float64()
					//计算累计套利
					initUSDT, initFIL, initTime := RunGetFilCost(robot, symbol+"-INIT")
					totalFilDelta, _ := currentFIL.Sub(initFIL).Float64()
					totalUsdtDelta, _ := currentUSDT.Sub(initUSDT).Float64()

					//说明又有卖单成交了，这次套利还要继续(要扣除回撤部分)
					if filDelta < float64(0-returnProfitFil) {
						stop = false
						msg := fmt.Sprintf("发生了买单已成交，但关闭前又有卖单成交，继续交易(filDelta：%v)...", filDelta)
						log.Logger.Error(msg)
						message.SendDingTalkRobit(true, robot, "fil2_every_continue_"+symbol, fmt.Sprintf("%v", time.Now().Unix()/3600), msg)
						continue
					}

					profitTimes++
					price, _ := getCurrentPrice()

					logmsg := fmt.Sprintf("E... %v 交易量:%v 套利:%vfil(%vusdt) 总套利:%vfil(%vusdt) \n套利次数:%v 余额:%vfil(%vusdt) 当前价格:%v",
						symbol, qty, filDelta, usdtDelta, totalFilDelta, totalUsdtDelta,
						profitTimes, currentFIL.String(), currentUSDT.String(), price)
					log.Logger.Debugf("[checkFinish] profit: %v", logmsg)
					message.SendDingTalkRobit(true, robot, "fil2_every_profit_"+symbol, fmt.Sprintf("%v", time.Now().Unix()), logmsg)

					//保存套利后最新余额,重置初始投入值
					filBalanceSaveFile(initTime, currentUSDT.String(), currentFIL.String())

					//余额保存到redis,后面做邮件报表使用
					filBalanceSaveRedis(currentFIL.String())
				}

				time.Sleep(time.Second * 5)

				//重新铺单
				cancelOrders(binance.SideTypeSell, symbol, openOrders)
				initSellOrders(true)
				placeSellLastTime = time.Now().Unix()
				RunSetInt64(robot, symbol, Filed_PlaceSellLastTime, placeSellLastTime)
				buyOrderId = 0
				RunSetInt64(robot, symbol, Filed_BuyOrderId, 0)
				stop = false
			}
		} else {
			if stopByBalance == true {
				continue
			}
			//本轮长时间没成交
			if buyOrderId == 0 && placeSellLastTime > 0 && time.Now().Unix()-placeSellLastTime > 10*60 {
				curPrice, err := getCurrentPrice()
				if err != nil {
					log.Logger.Error("[placeSells] Error getCurrentPrice:", err)
					continue
				}
				initPrice, err := RunGetInitPrice(robot, symbol)
				if err != nil {
					log.Logger.Error(err)
					continue
				}
				//当前价格小于当时铺单价格，说明是下跌导致的未成交需要重新铺单(如果当前价格比当时铺单价格还高仍未成交，说明价格太稳定不需要重新铺单）
				if curPrice < initPrice {
					stop = true
					cancelOrders(binance.SideTypeSell, symbol, openOrders)
					initSellOrders(true)
					placeSellLastTime = time.Now().Unix()
					RunSetInt64(robot, symbol, Filed_PlaceSellLastTime, placeSellLastTime)
					stop = false
				}
			}
		}
	}
}

func initSellOrders(init bool) error {
	if init {
		RunSetInitPrice(robot, symbol, 0)
		RunSetInt64(robot, symbol, Filed_EachRoundTime, time.Now().Unix())
	}

	//每次重启服务后，仍然根据上次的价格创建initSellPrice数组
	price, err := RunGetInitPrice(robot, symbol)
	if err != nil {
		log.Logger.Error(err)
		return err
	}
	if price == 0 {
		price, err = getCurrentPrice()
		if err != nil {
			log.Logger.Error(err)
			return err
		}
		RunSetInitPrice(robot, symbol, price)
	}

	//每轮的开始时间
	eachRoundTime := RunGetInt64(robot, symbol, Filed_EachRoundTime)
	if eachRoundTime == 0 {
		RunSetInt64(robot, symbol, Filed_EachRoundTime, time.Now().Unix())
	}

	log.Logger.Debugf("[initSellOrders] Current price: %f", price)

	initSellQty[0] = decimal.NewFromFloat(firstQty).Truncate(qtyPrecision)
	initRadios[0] = decimal.NewFromFloat(firstPriceIncrease - 1)
	initSellPrice[0] = decimal.NewFromFloat(price * math.Pow(firstPriceIncrease, 1)).Truncate(pricePrecision)
	for i := 0; i < maxSellOrders; i++ {
		qty := decimal.NewFromFloat(baseQty * math.Pow(baseIncrease, float64(i))).Truncate(qtyPrecision)
		initSellQty[i+1] = qty

		if i == 0 {
			initRadios[i+1] = decimal.NewFromFloat(priceIncrease - 1)
			initSellPrice[i+1] = decimal.NewFromFloat(price * math.Pow(priceIncrease, 1)).Truncate(pricePrecision)
		} else {
			initRadios[i+1] = (decimal.NewFromFloat(priceIncrease - 1)).Add(initRadios[i].Mul(decimal.NewFromFloat(priceFactor)))
			initSellPrice[i+1] = initRadios[i+1].Add(decimal.NewFromInt(1)).Mul(decimal.NewFromFloat(price)).Truncate(pricePrecision)
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

// 否：如果有成交卖单，定时会把成交卖单在补上 --> 容易占用大量资金，一直上涨亏损过大
func placeSells() (bool, int, int) {
	if stop == true || stopByBalance == true {
		return false, -1, -1
	}

	//没轮开始前检测，已经在进行中不需要检测
	if buyOrderId == 0 {
		price, err := RunGetInitPrice(robot, symbol)
		if err != nil {
			log.Logger.Error(err)
			return false, -1, -1
		}
		//当价格小于此值时自动暂停
		if price < minAllowPrice {
			initSellOrders(true)
			message.SendDingTalkRobit(true, robot, "fil2_every_autostop_"+symbol, fmt.Sprintf("%v", time.Now().Unix()/(60*60*12)), fmt.Sprintf("当前价格%v小于%v，不挂单。", price, minAllowPrice))
			return false, -1, -1
		}
	}

	openOrders, err := openOrders()
	if err != nil {
		log.Logger.Error("[placeSells] Error openOrders:", err)
		return false, -1, -1
	}

	//统计卖挂单数
	var openSells int
	var openBuys int
	for _, order := range openOrders {
		if order.Side == binance.SideTypeSell {
			openSells++
		}
		if order.Side == binance.SideTypeBuy {
			openBuys++
		}
	}

	var haveNewSell bool
	//卖挂单数少于活跃订单数
	if openSells < actionSellOrders {
		var redis = redis2.NewRedisCli(myconfig.GConfig.Redis.Host, myconfig.GConfig.Redis.Password, myconfig.GConfig.Redis.DB)

		curPrice, err := getCurrentPrice()
		if err != nil {
			log.Logger.Error("[placeSells] Error getCurrentPrice:", err)
			return false, openSells, openBuys
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
						key := fmt.Sprintf("sameSellPrice#%v#%v#%v", symbol, RunGetInt64(robot, symbol, Filed_EachRoundTime), sellPrice.String())
						val, err := redis.GetString(key)
						if err != nil {
							log.Logger.Error(err)
							return false, openSells, openBuys
						}
						//不存在进行补单和前几手可以重复挂单
						if val == "" || i < baseDoubleIndex {
							openSells++ //局部变量
							haveNewSell = true
							placeSellLastTime = time.Now().Unix()
							RunSetInt64(robot, symbol, Filed_PlaceSellLastTime, placeSellLastTime)
							if _, err := placeOrder("SELL", initSellQty[i].String(), sellPrice.String()); err == nil {
								redis.SetEX(key, "1", 7*24*3600*time.Second)
							}
						}
					}
				}
			}
		}

		//卖单已全部成交，买单还在进行中
		if openSells == 0 {
			message.SendDingTalkRobit(true, robot, "fil2_every_allsell_"+symbol, fmt.Sprintf("%v", time.Now().Unix()/(3600*8)), "卖单已全部成交，买单还在进行中")
		}
	}

	//超时没有成交了
	var timeout = (time.Now().Unix() - placeSellLastTime) / (3600 * 8)
	if placeSellLastTime > 0 && timeout > 1 {
		message.SendDingTalkRobit(true, robot, "fil2_every_sell_"+symbol, fmt.Sprintf("%v", time.Now().Unix()/(3600*8)), fmt.Sprintf("超过%v小时没有新卖单", timeout*8))
	}
	return haveNewSell, openSells, openBuys
}

func placeBuy(haveNewSell bool, openSells, openBuys int) {
	if stop == true || stopByBalance == true {
		return
	}
	//没有产生新卖单情况
	if haveNewSell == false {
		//placeSells()发生了异常 或 价格过小自动暂停了
		if openSells == -1 && openBuys == -1 {
			return
		}
		//当前没有买单:可能发生了取消买单成功但重新下买单失败，所以这里有个超过n分钟允许重新下单
		if openBuys == 0 && openSells >= 0 && time.Now().Unix()-placeBuyLastTime < 5*60 {
			return
		}
		//todo 当前有买单:可能已经达到最大挂单数，不会在挂新订单，此时需要定时重挂一次
		if openBuys > 0 && openSells >= 0 && time.Now().Unix()-placeBuyLastTime < 5*60 {
			return
		}
	}

	openOrders, err := openOrders()
	if err != nil {
		log.Logger.Error("[placeBuy] Error openOrders:", err)
		return
	}

	//如果买单已经是部分成交，不在进行重挂
	for _, order := range openOrders {
		if order.OrderID == buyOrderId && order.Status == binance.OrderStatusTypePartiallyFilled {
			return
		}
	}

	runInitialUSDT, runInitialFIL, _ := RunGetFilCost(robot, symbol)
	log.Logger.Debugf("[placeBuy] Initial balances: %s FIL, %s FDUSD", runInitialFIL.String(), runInitialUSDT.String())

	currentFIL, currentUSDT, _, err := getBalances()
	if err != nil {
		log.Logger.Error("[placeBuy] Error getBalances:", err)
		return
	}
	log.Logger.Debugf("[placeBuy] current balances: %s FIL, %s FDUSD", currentFIL.String(), currentUSDT.String())

	//计算卖掉的fil和获得的usdt
	var calcDelta = func(filDeltaOld float64, usdtBalance decimal.Decimal, initialUSDT decimal.Decimal) (float64, float64) {
		////还没有卖出或没有挂单
		//if filDeltaOld >= 0 {
		//	return 0, 0
		//}

		//当前挂单中最小的卖单价
		var minSellPrice = decimal.Zero
		//var partiallyFilled bool
		for _, order := range openOrders {
			var price, _ = decimal.NewFromString(order.Price)
			if order.Side == binance.SideTypeSell && (minSellPrice == decimal.Zero || price.Cmp(minSellPrice) < 0) {
				minSellPrice = price
			}
			//if order.Side == binance.SideTypeSell && order.Status == binance.OrderStatusTypePartiallyFilled {
			//	partiallyFilled = true
			//}
		}

		var totalUSDT decimal.Decimal
		var totalFil decimal.Decimal
		//已成交
		for i, sellPrice := range initSellPrice {
			if buyOrderId > 0 && minSellPrice.Cmp(decimal.Zero) == 0 { //全卖掉了
				totalUSDT = totalUSDT.Add(sellPrice.Mul(initSellQty[i]))
				totalFil = totalFil.Add(initSellQty[i])
			} else if minSellPrice.Cmp(decimal.Zero) > 0 && sellPrice.Cmp(minSellPrice) < 0 {
				totalUSDT = totalUSDT.Add(sellPrice.Mul(initSellQty[i]))
				totalFil = totalFil.Add(initSellQty[i])
			}
		}
		//部分成交
		for _, order := range openOrders {
			if order.Side == binance.SideTypeSell && order.Status == binance.OrderStatusTypePartiallyFilled {
				usdt1, _ := decimal.NewFromString(order.CummulativeQuoteQuantity)
				fil1, _ := decimal.NewFromString(order.ExecutedQuantity)
				totalUSDT = totalUSDT.Add(usdt1)
				totalFil = totalFil.Add(fil1)
			}
		}
		log.Logger.Debugf("[placeBuy] calcSell: %s FIL, %s FDUSD, minSellPrice: %v", totalFil.String(), totalUSDT.String(), minSellPrice)

		//if partiallyFilled == true {
		//	var curUsdt, _ = usdtBalance.Float64()
		//	var initUsdt, _ = initialUSDT.Float64()
		//	return curUsdt - initUsdt, true
		//} else {
		//两者取小
		if totalUSDT.Cmp(usdtBalance) > 0 {
			tFil, _ := totalFil.Float64()
			tUsdt, _ := usdtBalance.Float64()
			return -tFil, tUsdt
		} else {
			tFil, _ := totalFil.Float64()
			tUsdt, _ := totalUSDT.Float64()
			////todo 调整参数后需要先用余额计算方式
			//tUsdt, _ = usdtBalance.Float64()
			return -tFil, tUsdt
		}
	}

	filDelta, _ := currentFIL.Sub(runInitialFIL).Float64()
	//usdtDelta, _ := currentUSDT.Sub(runInitialUSDT).Float64()
	filDeltaNew, usdtDelta := calcDelta(filDelta, currentUSDT, runInitialUSDT)
	if filDeltaNew != filDelta {
		message.SendDingTalkRobit(true, robot, "fil2_every_checkfil_"+symbol,
			fmt.Sprintf("%v", time.Now().Unix()/(8*60*60)),
			fmt.Sprintf("两种计算结束不一致(filDelta: %v,filDeltaNew: %v)", filDelta, filDeltaNew))
	}
	log.Logger.Debugf("[placeBuy] filDelta:%v, filDeltaNew: %v, usdtDelta: %v buyOrderId: %v", filDelta, filDeltaNew, usdtDelta, buyOrderId)

	//fil为负，表示已有卖单成交，开始挂买单
	if filDelta < 0 {
		if usdtDelta <= 0 {
			logmsg := "异常:套利还未执行完，U的余额增量居然小于等于0"
			message.SendDingTalkRobit(true, robot, "fil2_every_profit_"+symbol, fmt.Sprintf("%v", time.Now().Unix()/(10*60)), logmsg)
			return
		}

		var tmpReturnProfitFil int64
		//可能一直上涨没有大的回调，这时需要把之前的收益拿出来，减少本次的买回量(fil),确保可以成交（收益回撤了）
		_, firstInitialFIL, _ := RunGetFilCost(robot, symbol+"-INIT")
		totalProfitFil, _ := runInitialFIL.Sub(firstInitialFIL).Float64()
		if buySuccLastTime > 0 && time.Now().Unix()-buySuccLastTime > 24*3600 {
			tmpReturnProfitFil = int64(totalProfitFil / 3)
			filDelta = filDelta + float64(tmpReturnProfitFil)
		} else if buySuccLastTime > 0 && time.Now().Unix()-buySuccLastTime > 20*3600 {
			tmpReturnProfitFil = int64(totalProfitFil / 5)
			filDelta = filDelta + float64(tmpReturnProfitFil)
		} else if buySuccLastTime > 0 && time.Now().Unix()-buySuccLastTime > 12*3600 {
			tmpReturnProfitFil = int64(totalProfitFil / 10)
			filDelta = filDelta + float64(tmpReturnProfitFil)
		}

		//要买回的币
		filToBuyBack := decimal.NewFromFloat((0-filDelta)*(1+profitFactor)).DivRound(decimal.NewFromInt(1), qtyPrecision)
		//购买价
		newBuyPrice := decimal.NewFromFloat(usdtDelta).Div(filToBuyBack).Truncate(pricePrecision)

		curPrice, err := getCurrentPrice()
		if err != nil {
			log.Logger.Error("[placeBuy] Error fetching current price:", err)
			return
		}
		//新购买价高于当前价
		if newBuyPrice.Cmp(decimal.NewFromFloat(curPrice)) > 0 {
			newBuyPrice = decimal.NewFromFloat(curPrice).Mul(decimal.NewFromFloat(1 - 0.0005)).Truncate(pricePrecision)
		}

		//取消买单
		if buyOrderId > 0 {
			cancelOrder(buyOrderId)
			time.Sleep(time.Millisecond * 500)
		}

		//计算新买价和数量
		var needUsdt = filToBuyBack.Mul(newBuyPrice)
		var needDownRatio = decimal.NewFromFloat(curPrice).Sub(newBuyPrice).Div(decimal.NewFromFloat(curPrice)).Truncate(4).Mul(decimal.NewFromFloat(100)).String() + "%"
		log.Logger.Debugf("[placeBuy] filToBuyBack: %v, newBuyPrice: %v curPrice: %v(needDownRatio:%v) (needUsdt: %v, usdtBalance: %v)", filToBuyBack, newBuyPrice, curPrice, needDownRatio, needUsdt, currentUSDT)

		orderId, err := placeOrder("BUY", filToBuyBack.String(), newBuyPrice.String())
		if err != nil {
			log.Logger.Error("[placeBuy] Error placeOrder:", err)
			message.SendDingTalkRobit(true, robot, "fil2_every_buy_"+symbol, fmt.Sprintf("%v", time.Now().Unix()/(4*60*60)), err.Error())
		} else {
			buyOrderId = orderId
			returnProfitFil = tmpReturnProfitFil
			RunSetInt64(robot, symbol, Filed_BuyOrderId, buyOrderId)
			placeBuyLastTime = time.Now().Unix()
		}
	}
}

// stopBalance: 0暂停服务，其它值恢复
func checkStopByBalance(currentUSDT, currentFIL string, stopBalance decimal.Decimal) {
	if stopByBalance == false && stopBalance.Cmp(decimal.NewFromFloat(0)) == 0 {
		//确认是否结束
		if buyOrderId > 0 {
			message.SendDingTalkRobit(true, robot, "fil2_every_stop1_"+symbol, fmt.Sprintf("%v", time.Now().Unix()/(3600*8)), "待本轮结束后暂停...")
			orderstatus, _, err := getOrderStatus(symbol, buyOrderId)
			if err != nil {
				log.Logger.Errorf("[checkStopByBalance] getOrderStatus %v %v", buyOrderId, err)
				return
			}
			//还未结束
			if orderstatus == false {
				return
			}
		}

		//已经结束，可以暂停
		stopByBalance = true
		for {
			//取消可能存在的未成交卖挂单
			openOrders, err := openOrders()
			if err != nil {
				log.Logger.Error(err)
				time.Sleep(time.Second * 5)
				continue
			}
			cancelOrders(binance.SideTypeSell, symbol, openOrders)
			break
		}
		message.SendDingTalkRobit(true, robot, "fil2_every_stop2_"+symbol, fmt.Sprintf("%v", time.Now().Unix()), "已暂停")
	}

	if stopByBalance == true && stopBalance.Cmp(decimal.NewFromFloat(0)) != 0 {
		RunSetFilCost(robot, symbol, currentUSDT, currentFIL)
		for {
			//取消可能存在的未成交卖挂单
			openOrders, err := openOrders()
			if err != nil {
				log.Logger.Error(err)
				time.Sleep(time.Second * 5)
				continue
			}
			cancelOrders(binance.SideTypeSell, symbol, openOrders)
			break
		}
		initSellOrders(true)
		placeSellLastTime = time.Now().Unix()
		RunSetInt64(robot, symbol, Filed_PlaceSellLastTime, placeSellLastTime)

		//重置回来后重新计算，避免执行回撤
		buySuccLastTime = 0

		stopByBalance = false

		message.SendDingTalkRobit(true, robot, "fil2_every_start_"+symbol, fmt.Sprintf("%v", time.Now().Unix()), "已恢复")
	}
}

func filBalanceSaveFile(initTime int64, currentUSDT, currentFIL string) {
	//保存当前余额
	log.Logger.Debugf("[checkFinish] Initial balances: %s FIL, %s FDUSD", currentFIL, currentUSDT)
	RunSetFilCost(robot, symbol, currentUSDT, currentFIL)

	//每24小时结算一次：重置初始投入值，在回撤计算时最多回撤24小时收益
	if (time.Now().Unix()-initTime)/(24*3600) >= 1 {
		RunSetFilCost(robot, symbol+"-INIT", currentUSDT, currentFIL)
	}
}

// 每天余额保存到redis,后面做邮件报表使用
func filBalanceSaveRedis(currentFIL string) {
	var redis = redis2.NewRedisCli("localhost:6379", "", 0)
	var key = fmt.Sprintf("%v-fildown-%v", robot, time.Now().Format("2006-01-02"))
	redis.SetEX(key, currentFIL, 60*24*3600*time.Second)
}
