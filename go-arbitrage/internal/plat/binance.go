/**
 * binance.go
 * ============================================================================
 * https://github.com/adshao/go-binance
 * ============================================================================
 * author: peter.wang
 */

package plat

import (
	"components/message"
	"context"
	"errors"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/common"
	"github.com/shopspring/decimal"
	"goarbitrage/api"
	"goarbitrage/internal/model"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Binance struct {
	Base
	client      map[string]*binance.Client
	initial     map[string]map[string]model.Balance
	balances    sync.Map //创建/取消/完成订单时更新余额（降底请求量，权重太大)
	bookTickers sync.Map //实时最优挂单
	cache       sync.Map //缓存数据
}

func NewBinance() *Binance {
	obj := &Binance{}

	obj.client = make(map[string]*binance.Client)
	obj.initial = make(map[string]map[string]model.Balance)
	//obj.client["BTCUSDT"].Debug = true

	//mainapi
	obj.client["BTCFDUSD"] = binance.NewClient("mCXfycRaEiffizOajnB1VsVxytyUFnaA1tK4eX8QyuM8G565Weq5s4QXoyhkzwdE", "wvRdYxo9O4IeBywbDCZgGhflwDwv2ERUbdQHUgoZ8JXTpUDGvFsTnXtzQOHxL9XW")
	obj.initial["BTCFDUSD"] = map[string]model.Balance{
		"BTC":   model.Balance{Asset: "BTC", Total: decimal.NewFromFloat(0), Free: decimal.NewFromFloat(0)},
		"FDUSD": model.Balance{Asset: "FDUSD", Total: decimal.NewFromFloat(30 * 10000), Free: decimal.NewFromFloat(30 * 10000)},
	}
	////oneplat_virtual@***.com
	//obj.client["BTCFDUSD"] = binance.NewClient("3JiMItY7JeQoxWNAlylhsxCI38hysP5OZUgypWewm3PhKUaVx9pMv3dTUlyT5sbS", "iPP0IRusNqUhKtVyl0gSteRnTEpMUXttUWSekH2MeqljcLkfzwyJ6J8nmUyUOxhn")
	//obj.initial["BTCFDUSD"] = map[string]model.Balance{
	//	"BTC":   model.Balance{Asset: "BTC", Total: decimal.NewFromFloat(0), Free: decimal.NewFromFloat(0)},
	//	"FDUSD": model.Balance{Asset: "FDUSD", Total: decimal.NewFromFloat(10 * 10000), Free: decimal.NewFromFloat(10 * 10000)},
	//}

	//btcbusd_virtual@***.com
	obj.client["MEMEUSDT"] = binance.NewClient("JUZWoFOfVM7abdnHYHgSevlGIfi6XyrZCWZc8YhKHDCor5g1An55CUCvrhS6aGx9", "t7SPjnakqH6F0bWe2dzh1TAREHMuzY0OS4sk7YV0Brfzyl1noExAV6Bnwd95F7si")
	obj.initial["MEMEUSDT"] = map[string]model.Balance{
		"MEME": model.Balance{Asset: "MEME", Total: decimal.NewFromFloat(0), Free: decimal.NewFromFloat(0)},
		"USDT": model.Balance{Asset: "USDT", Total: decimal.NewFromFloat(0), Free: decimal.NewFromFloat(0)},
	}

	//btcusdt_virtual@***.com
	obj.client["ARBUSDT"] = binance.NewClient("O6VRBuk0qYI3Q5J0t1sqsHXCjbO3HyDr1ef9KqtGABqneBbLzrCJuIEzHgkjtpap", "bHAr64boBpbVHpu2ZbpIIiFAlP40VuhDjeqxPFXSG4GewlWWVvB2mZk8xjoQFJgz")
	obj.initial["ARBUSDT"] = map[string]model.Balance{
		"ARB":  model.Balance{Asset: "ARB", Total: decimal.NewFromFloat(0), Free: decimal.NewFromFloat(0)},
		"USDT": model.Balance{Asset: "USDT", Total: decimal.NewFromFloat(20000), Free: decimal.NewFromFloat(20000)},
	}
	obj.client["ARBFDUSD"] = binance.NewClient("O6VRBuk0qYI3Q5J0t1sqsHXCjbO3HyDr1ef9KqtGABqneBbLzrCJuIEzHgkjtpap", "bHAr64boBpbVHpu2ZbpIIiFAlP40VuhDjeqxPFXSG4GewlWWVvB2mZk8xjoQFJgz")
	obj.initial["ARBFDUSD"] = map[string]model.Balance{
		"ARB":   model.Balance{Asset: "ARB", Total: decimal.NewFromFloat(0), Free: decimal.NewFromFloat(0)},
		"FDUSD": model.Balance{Asset: "FDUSD", Total: decimal.NewFromFloat(20000), Free: decimal.NewFromFloat(20000)},
	}

	//busdusdt_virtual@***.com
	obj.client["BNBFDUSD"] = binance.NewClient("4O7Pa8bbwTBcsKMEWb9qlwoXUU56pJEhOMbI4nTyQKDTxE1GhAp0tW7HR9dpFZbs", "v1zTAZ0y5GfnKSf5QZ9u9NX4WaOsIPoHwnh5C2KinGgBBQyT8YXffuyGWxORI559")
	obj.initial["BNBFDUSD"] = map[string]model.Balance{
		"BNB":   model.Balance{Asset: "BNB", Total: decimal.NewFromFloat(0), Free: decimal.NewFromFloat(0), AddTime: time.Date(2023, 11, 8, 0, 0, 0, 0, time.Local).Unix()},
		"FDUSD": model.Balance{Asset: "FDUSD", Total: decimal.NewFromFloat(100000), Free: decimal.NewFromFloat(100000), AddTime: time.Date(2023, 11, 9, 0, 0, 0, 0, time.Local).Unix()},
	}
	obj.client["FDUSDUSDT"] = binance.NewClient("4O7Pa8bbwTBcsKMEWb9qlwoXUU56pJEhOMbI4nTyQKDTxE1GhAp0tW7HR9dpFZbs", "v1zTAZ0y5GfnKSf5QZ9u9NX4WaOsIPoHwnh5C2KinGgBBQyT8YXffuyGWxORI559")
	obj.initial["FDUSDUSDT"] = map[string]model.Balance{
		"FDUSD": model.Balance{Asset: "FDUSD", Total: decimal.NewFromFloat(0), Free: decimal.NewFromFloat(0), AddTime: time.Date(2023, 11, 9, 0, 0, 0, 0, time.Local).Unix()},
		"USDT":  model.Balance{Asset: "USDT", Total: decimal.NewFromFloat(0), Free: decimal.NewFromFloat(0), AddTime: time.Date(2023, 11, 9, 0, 0, 0, 0, time.Local).Unix()},
	}

	var symbols = make([]string, 0)
	for key, _ := range obj.client {
		symbols = append(symbols, key)
	}

	//实时接收最佳挂单
	go func(oobj *Binance) {
		var doneC, stopC chan struct{}
		var wserr error
		var handler = func(event *binance.WsBookTickerEvent) {
			var marketDepthReturn1 = model.MarketDepthReturn{}
			price1, _ := strconv.ParseFloat(event.BestBidPrice, 10)
			num1, _ := strconv.ParseFloat(event.BestBidQty, 10)
			marketDepthReturn1.Tick.Bids = make([][]float64, 1)
			marketDepthReturn1.Tick.Bids[0] = make([]float64, 2)
			marketDepthReturn1.Tick.Bids[0][0] = price1 //价格
			marketDepthReturn1.Tick.Bids[0][1] = num1   //数量

			price2, _ := strconv.ParseFloat(event.BestAskPrice, 10)
			num2, _ := strconv.ParseFloat(event.BestAskQty, 10)
			marketDepthReturn1.Tick.Asks = make([][]float64, 1)
			marketDepthReturn1.Tick.Asks[0] = make([]float64, 2)
			marketDepthReturn1.Tick.Asks[0][0] = price2 //价格
			marketDepthReturn1.Tick.Asks[0][1] = num2   //数量
			oobj.bookTickers.Store(event.Symbol, marketDepthReturn1)
			//fmt.Println(event.Symbol, marketDepthReturn1)
		}
		var errHandler = func(err error) {
			if err != nil {
				wserr = err
				oobj.bookTickers = sync.Map{}

				logmsg := fmt.Sprintf("errHandler:%v", err.Error())
				message.PrintLog(api.PrintLog, "Binance.getAccountBalance", fmt.Sprintf("%v", time.Now().Unix()/(60*10)), logmsg)
				message.SendDingTalkRobit(true, "oneplat", "Binance.WsBookTicker", fmt.Sprintf("%v", time.Now().Unix()/(60*60)), logmsg)
			}
		}
		doneC, stopC, wserr = binance.WsCombinedBookTickerServe(symbols, handler, errHandler)
		for {
			if wserr != nil {
				if stopC != nil {
					stopC <- struct{}{}
				}
				if doneC != nil {
					<-doneC
				}

				doneC, stopC, wserr = binance.WsCombinedBookTickerServe(symbols, handler, errHandler)
				if wserr != nil {
					oobj.bookTickers = sync.Map{}

					logmsg := fmt.Sprintf("WsServe:%v", wserr.Error())
					message.PrintLog(api.PrintLog, "Binance.getAccountBalance", fmt.Sprintf("%v", time.Now().Unix()/(60*10)), logmsg)
					message.SendDingTalkRobit(true, "oneplat", "Binance.WsBookTicker", fmt.Sprintf("%v", time.Now().Unix()/(60*60)), logmsg)
				}
			}
			time.Sleep(time.Second)
		}
	}(obj)

	//n秒同步一次余额
	go func(oobj *Binance) {
		for {
			for symbol, _ := range oobj.client {
				symbol = oobj.FormatSymbol(symbol)
				_, _, err := oobj.updateAccountBalance(symbol)
				if err != nil {
					break
				}
				time.Sleep(time.Second * 20)
			}
			time.Sleep(time.Second * 60 * 5)
		}
	}(obj)

	return obj
}

// 灵活测试子帐号
func (this *Binance) SetKey(symbol string, apiKey, secretKey string, initial map[string]model.Balance) {
	symbol = this.FormatSymbol(symbol)
	this.client[symbol] = binance.NewClient(apiKey, secretKey)
	this.initial[symbol] = initial
}

func (this *Binance) GetInitialInput(symbol string) map[string]model.Balance {
	symbol = this.FormatSymbol(symbol)

	if initial, ok := this.initial[symbol]; ok {
		return initial
	} else {
		panic(fmt.Sprintf("%v not exist initial input", symbol))
	}
}

func (this *Binance) PlatCode() string {
	return P_Binance
}

// 币安交易对格式:ETHUSDT
func (this *Binance) FormatSymbol(symbol string) string {
	return strings.Replace(symbol, "/", "", 1)
}

func (this *Binance) GetMarketDepth(symbol string) (model.MarketDepthReturn, error) {
	symbol = this.FormatSymbol(symbol)

	if val, ok := this.bookTickers.Load(symbol); ok {
		value, ok := val.(model.MarketDepthReturn)
		if ok {
			return value, nil
		} else {
			return model.MarketDepthReturn{}, errors.New("nil")
		}
	} else {
		this.wait(this.PlatCode(), "GetMarketDepth", 2, symbol)

		resp, err := this.client[symbol].NewListBookTickersService().Symbol(symbol).Do(context.Background())
		if err != nil {
			return model.MarketDepthReturn{}, err
		}
		var marketDepthReturn1 = model.MarketDepthReturn{}
		for _, ticker := range resp {
			price1, _ := strconv.ParseFloat(ticker.BidPrice, 10)
			num1, _ := strconv.ParseFloat(ticker.BidQuantity, 10)
			marketDepthReturn1.Tick.Bids = make([][]float64, 1)
			marketDepthReturn1.Tick.Bids[0] = make([]float64, 2)
			marketDepthReturn1.Tick.Bids[0][0] = price1 //价格
			marketDepthReturn1.Tick.Bids[0][1] = num1   //数量

			price2, _ := strconv.ParseFloat(ticker.AskPrice, 10)
			num2, _ := strconv.ParseFloat(ticker.AskQuantity, 10)
			marketDepthReturn1.Tick.Asks = make([][]float64, 1)
			marketDepthReturn1.Tick.Asks[0] = make([]float64, 2)
			marketDepthReturn1.Tick.Asks[0][0] = price2 //价格
			marketDepthReturn1.Tick.Asks[0][1] = num2   //数量
		}
		return marketDepthReturn1, nil
	}
}

func (this *Binance) AveragePrice(symbol string) (decimal.Decimal, error) {
	symbol = this.FormatSymbol(symbol)

	this.wait(this.PlatCode(), "AveragePrice", 2, symbol)

	avgprice, err := this.client[symbol].NewAveragePriceService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return decimal.Zero, err
	}
	decPrice, err := decimal.NewFromString(avgprice.Price)
	return decPrice, err
}

func (this *Binance) GetAccountBalance(symbol string, realtime bool) (map[string]model.Balance, error) {
	symbol = this.FormatSymbol(symbol)

	if val, ok := this.balances.Load(symbol); ok == true && realtime == false {
		values := val.([]model.Balance)
		var res = make(map[string]model.Balance)
		for _, val := range values {
			res[val.Asset] = val
		}
		return res, nil
	} else {
		res, _, err := this.updateAccountBalance(symbol)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func (this *Binance) updateAccountBalance(symbol string) (map[string]model.Balance, []model.Balance, error) {
	this.wait2(this.PlatCode()+"2", "GetAccountBalance", 1, symbol)

	account, err := this.client[symbol].NewGetAccountService().Do(context.Background())
	if err != nil {
		this.balances.Delete(symbol)
		logmsg := fmt.Sprintf("updateAccountBalance:%v", err.Error())
		message.PrintLog(api.PrintLog, "Binance.getAccountBalance", fmt.Sprintf("%v", time.Now().Unix()/10), logmsg)
		message.SendDingTalkRobit(true, "oneplat", "Binance.getAccountBalance", fmt.Sprintf("%v", time.Now().Unix()/(60*60*24)), logmsg)
		return nil, nil, err
	}

	var res = make(map[string]model.Balance)
	var values = make([]model.Balance, 0)
	for _, balance := range account.Balances {
		free, _ := decimal.NewFromString(balance.Free)
		locked, _ := decimal.NewFromString(balance.Locked)
		total := free.Add(locked)
		if total.Cmp(decimal.NewFromFloat(0)) > 0 {
			res[balance.Asset] = model.Balance{Asset: balance.Asset, Total: total, Free: free, Locked: locked}
			values = append(values, model.Balance{Asset: balance.Asset, Total: total, Free: free, Locked: locked})
		}
	}

	this.balances.Store(symbol, values)
	return res, values, nil
}

func (this *Binance) OrdersPlace(symbol string, price string, quantity string, side string) (string, error) {
	symbol = this.FormatSymbol(symbol)

	var sidetype binance.SideType
	if side == "BUY" {
		sidetype = binance.SideTypeBuy
	} else {
		sidetype = binance.SideTypeSell
	}

	order, err := this.client[symbol].NewCreateOrderService().Symbol(symbol).
		Side(sidetype).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(price).Do(context.Background())
	if err != nil {
		return "", err
	}

	//下单成功，更新余额
	this.updateAccountBalance(symbol)

	return fmt.Sprintf("%v", order.OrderID), nil
}

func (this *Binance) OrdersPlace2(clientsymbol string, symbol string, price string, quantity string, side string) (string, error) {
	clientsymbol = this.FormatSymbol(clientsymbol)
	symbol = this.FormatSymbol(symbol)

	var sidetype binance.SideType
	if side == "BUY" {
		sidetype = binance.SideTypeBuy
	} else {
		sidetype = binance.SideTypeSell
	}

	order, err := this.client[clientsymbol].NewCreateOrderService().Symbol(symbol).
		Side(sidetype).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(price).Do(context.Background())
	if err != nil {
		return "", err
	}

	//下单成功，更新余额
	this.updateAccountBalance(clientsymbol)

	return fmt.Sprintf("%v", order.OrderID), nil
}

func (this *Binance) CancelOrder(symbol string, orderid string) (bool, error) {
	symbol = this.FormatSymbol(symbol)

	this.wait(this.PlatCode(), "CancelOrder", 2, symbol)

	orderId, _ := strconv.ParseInt(orderid, 10, 64)

	_, err := this.client[symbol].NewCancelOrderService().Symbol(symbol).OrderID(orderId).Do(context.Background())
	if err != nil {
		commerr, ok := err.(*common.APIError)
		//已经成交或已经取消
		if ok && commerr.Code == -2011 {
			//取消订单，更新余额
			this.updateAccountBalance(symbol)
			return true, nil
		}
		return false, err
	}

	//取消订单，更新余额
	this.updateAccountBalance(symbol)

	return true, nil
}

func (this *Binance) CancelOrder2(clientsymbol string, symbol string, orderid string) (bool, error) {
	clientsymbol = this.FormatSymbol(clientsymbol)
	symbol = this.FormatSymbol(symbol)

	this.wait(this.PlatCode(), "CancelOrder", 2, symbol)

	orderId, _ := strconv.ParseInt(orderid, 10, 64)

	_, err := this.client[clientsymbol].NewCancelOrderService().Symbol(symbol).OrderID(orderId).Do(context.Background())
	if err != nil {
		commerr, ok := err.(*common.APIError)
		//已经成交或已经取消
		if ok && commerr.Code == -2011 {
			//取消订单，更新余额
			this.updateAccountBalance(clientsymbol)
			return true, nil
		}
		return false, err
	}

	//取消订单，更新余额
	this.updateAccountBalance(clientsymbol)

	return true, nil
}

func (this *Binance) GetOrderStatus(symbol string, orderid string) (bool, error) {
	symbol = this.FormatSymbol(symbol)

	//订单已经完成或超时
	if orderid == "" || strings.Index(orderid, "FILLED") >= 0 {
		fmt.Println("---------------------- why", fmt.Sprintf(".%v.", orderid))
		return true, nil
	}

	//去除超时未完交易TIMEOUT前缀
	orderid = strings.Replace(orderid, "TIMEOUT", "", 1)

	this.wait(this.PlatCode(), "GetOrders", 2, symbol)

	orderId, _ := strconv.ParseInt(orderid, 10, 64)

	order, err := this.client[symbol].NewGetOrderService().Symbol(symbol).OrderID(orderId).Do(context.Background())
	if err != nil {
		//不存在
		commerr, ok := err.(*common.APIError)
		if ok && commerr.Code == -2013 {
			panic(fmt.Sprintf("%v %v %v", orderId, orderid, err))
		}
		return false, err
	}

	if order.Status == binance.OrderStatusTypePartiallyFilled || order.Status == binance.OrderStatusTypeFilled {
		//订单执行更新余额
		this.updateAccountBalance(symbol)
	}

	if order.OrigQuantity == order.ExecutedQuantity {
		return true, nil
	}

	if order.Status == binance.OrderStatusTypeCanceled {
		return true, nil
	}
	return false, nil
}

func (this *Binance) GetOrderStatus2(clientsymbol string, symbol string, orderid string) (bool, error) {
	clientsymbol = this.FormatSymbol(clientsymbol)
	symbol = this.FormatSymbol(symbol)

	//订单已经完成或超时
	if orderid == "" || strings.Index(orderid, "FILLED") >= 0 {
		fmt.Println("---------------------- why", fmt.Sprintf(".%v.", orderid))
		return true, nil
	}

	//去除超时未完交易TIMEOUT前缀
	orderid = strings.Replace(orderid, "TIMEOUT", "", 1)

	this.wait(this.PlatCode(), "GetOrders", 2, symbol)

	orderId, _ := strconv.ParseInt(orderid, 10, 64)

	order, err := this.client[clientsymbol].NewGetOrderService().Symbol(symbol).OrderID(orderId).Do(context.Background())
	if err != nil {
		//不存在
		commerr, ok := err.(*common.APIError)
		if ok && commerr.Code == -2013 {
			panic(fmt.Sprintf("%v %v %v", orderId, orderid, err))
		}
		return false, err
	}

	if order.Status == binance.OrderStatusTypePartiallyFilled || order.Status == binance.OrderStatusTypeFilled {
		//订单执行更新余额
		this.updateAccountBalance(clientsymbol)
	}

	if order.OrigQuantity == order.ExecutedQuantity {
		return true, nil
	}

	if order.Status == binance.OrderStatusTypeCanceled {
		return true, nil
	}
	return false, nil
}

func (this *Binance) OpenOrders(symbol string) ([]*model.Order, error) {
	symbol = this.FormatSymbol(symbol)

	this.wait(this.PlatCode(), "OpenOrdersCount", 2, symbol)

	orders, err := this.client[symbol].NewListOpenOrdersService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return nil, err
	}

	var res = make([]*model.Order, len(orders))
	for i, order := range orders {
		res[i] = &model.Order{
			Symbol:           order.Symbol,
			OrderID:          order.OrderID,
			Price:            order.Price,
			OrigQuantity:     order.OrigQuantity,
			ExecutedQuantity: order.ExecutedQuantity,
			Side:             fmt.Sprintf("%v", order.Side),
			Time:             order.Time,
		}
	}

	return res, nil
}

func (this *Binance) OpenOrders2(clientsymbol string, symbol string) ([]*model.Order, error) {
	clientsymbol = this.FormatSymbol(clientsymbol)
	symbol = this.FormatSymbol(symbol)

	this.wait(this.PlatCode(), "OpenOrdersCount", 2, symbol)

	orders, err := this.client[clientsymbol].NewListOpenOrdersService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return nil, err
	}

	var res = make([]*model.Order, len(orders))
	for i, order := range orders {
		res[i] = &model.Order{
			Symbol:           order.Symbol,
			OrderID:          order.OrderID,
			Price:            order.Price,
			OrigQuantity:     order.OrigQuantity,
			ExecutedQuantity: order.ExecutedQuantity,
			Side:             fmt.Sprintf("%v", order.Side),
			Time:             order.Time,
		}
	}

	return res, nil
}

func (this *Binance) CancelOpenOrders(symbol string) error {
	symbol = this.FormatSymbol(symbol)

	_, err := this.client[symbol].NewCancelOpenOrdersService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (this *Binance) GetTradeFee(symbol string) (makerfee decimal.Decimal, takerfee decimal.Decimal, err error) {
	type cacheFee struct {
		Time     int64
		MakerFee decimal.Decimal
		TakerFee decimal.Decimal
	}

	symbol = this.FormatSymbol(symbol)

	key := fmt.Sprintf("GetTradeFee_%v", symbol)

	fee, ok := this.cache.Load(key)

	//有效期内
	if ok == true && time.Now().Unix()-fee.(cacheFee).Time < 60*2 {
		makerfee = fee.(cacheFee).MakerFee
		takerfee = fee.(cacheFee).TakerFee
		return
	}

	var fees []*binance.TradeFeeDetails
	fees, err = this.client[symbol].NewTradeFeeService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return
	}
	makerfee, err = decimal.NewFromString(fees[0].MakerCommission) //挂单
	if err != nil {
		return
	}
	takerfee, err = decimal.NewFromString(fees[0].TakerCommission) //吃单
	if err != nil {
		return
	}
	this.cache.Store(key, cacheFee{Time: time.Now().Unix(), MakerFee: makerfee, TakerFee: takerfee})
	return
}
