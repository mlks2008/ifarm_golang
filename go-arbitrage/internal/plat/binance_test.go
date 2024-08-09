package plat

import (
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"goarbitrage/internal/model"
	"goarbitrage/pkg/utils"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func init() {
	viper.SetConfigFile("../../configs/hotbit_dev.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func Test_Binance_OneTicker(t *testing.T) {
	plat := NewBinance()

	fmt.Println(plat.OpenOrders("BTCUSDT"))
	//fmt.Println(plat.GetTradeFee("BTCUSDT"))
	//fmt.Println(plat.CancelOrder("BTCUSDT", "1"))

	////币种信息
	//resp, err := plat.client.NewExchangeInfoService().Do(context.Background())
	//symbolsReturn := &model.SymbolsReturn{}
	//for _, ss := range resp.Symbols {
	//	if ss.Symbol == "BTCUSDT" {
	//		fmt.Println(ss.Filters)
	//	}
	//	symbolsReturn.Data = append(symbolsReturn.Data, model.SymbolsData{
	//		BaseCurrency:  ss.BaseAsset,
	//		QuoteCurrency: ss.QuoteAsset,
	//		Symbol:        ss.Symbol,
	//	})
	//}

	//prices, err := plat.client.NewListPricesService().Do(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//for _, p := range prices {
	//	fmt.Println(p)
	//}

	////深度
	//res, err := plat.client.NewDepthService().Symbol("BTCUSDT").Do(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(res.Bids[0], res.Asks[0])

	////最优买卖挂单
	//bookTicker, err := plat.client.NewListBookTickersService().Symbol("BTCUSDT").Do(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//for _, ticker := range bookTicker {
	//	fmt.Println(ticker)
	//}

	time.Sleep(time.Second)
}

// 快速上涨跟一单
func Test_Binance_BookTicker(t *testing.T) {
	plat := NewBinance()

	var list, _ = plat.GetSymbols()
	var symbols = make([]string, 0)
	var amountPrecisions = make(map[string]int32)
	var pricePrecisions = make(map[string]int32)
	var quoteCurrency = "USDT"
	for _, symbolinfo := range list.Data {
		if symbolinfo.QuoteCurrency == quoteCurrency {
			symbols = append(symbols, symbolinfo.Symbol)
			amountPrecisions[symbolinfo.Symbol] = symbolinfo.AmountPrecision
			pricePrecisions[symbolinfo.Symbol] = symbolinfo.PricePrecision
		}
	}
	//symbols = []string{"BTCUSDT", "FILUSDT"}

	var buysell = func(symbol string, bprice float64, factor float64) {
		symbolBase := strings.TrimRight(symbol, quoteCurrency)
		symbolQuote := quoteCurrency

		buyprice := utils.DivRound(decimal.NewFromFloat(bprice).Mul(decimal.NewFromFloat(factor)), pricePrecisions[symbol])
		sellprice := utils.DivRound(decimal.NewFromFloat(bprice).Mul(decimal.NewFromFloat(factor+0.01)), pricePrecisions[symbol])

		var usdtamount int64 = 20
		var decbuyprice, _ = decimal.NewFromString(buyprice)
		var btcamount = decimal.NewFromInt(usdtamount).DivRound(decbuyprice, amountPrecisions[symbol]).String()

		buyorderid, err := plat.OrdersPlace2("BUSD/USDT", symbol, buyprice, btcamount, "BUY")
		if err != nil {
			fmt.Println(err)
			return
		}
		btime := time.Now().Unix()
		for {
			status, _ := plat.GetOrderStatus2("BUSD/USDT", symbol, buyorderid)
			if status == true {
				time.Sleep(time.Millisecond * 1000)
				balance, _ := plat.GetAccountBalance("BUSD/USDT", true)
				if balance[symbolBase].Free.Cmp(decimal.NewFromFloat(0)) > 0 {
					btcamount = utils.DivRound(balance[symbolBase].Free, amountPrecisions[symbol])
					orderid, err := plat.OrdersPlace2("BUSD/USDT", symbol, sellprice, btcamount, "SELL")
					if err != nil {
						panic(err)
					}
					for {
						status, _ := plat.GetOrderStatus2("BUSD/USDT", symbol, orderid)
						if status == true {
							time.Sleep(time.Millisecond * 500)
							endbalance, _ := plat.GetAccountBalance("BUSD/USDT", true)
							fmt.Println(symbol, "买卖完成", balance[symbolQuote].Free.Add(decimal.NewFromInt(usdtamount)), endbalance[symbolQuote].Free)
							return
						}
					}
				} else {
					time.Sleep(time.Second)
				}
			} else {
				//超时取消买单
				if time.Now().Unix()-btime > 1*60 {
					res, err := plat.CancelOrder2("BUSD/USDT", symbol, buyorderid)
					fmt.Println("取消买单", res, err)
				} else {
					time.Sleep(time.Millisecond * 500)
				}
			}
		}
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

			var min float64 = 0.1
			lastkeymin := fmt.Sprintf("%v_%v", event.Symbol, time.Now().Unix()/int64(min*60)-1)
			nowkeymin := fmt.Sprintf("%v_%v", event.Symbol, time.Now().Unix()/int64(min*60))

			if _, ok := oobj.bookTickers.Load(nowkeymin); ok == false {
				oobj.bookTickers.Store(nowkeymin, marketDepthReturn1)
			}

			if v, ok := oobj.bookTickers.Load(lastkeymin); ok {
				lastMarketDepthReturn := v.(model.MarketDepthReturn)

				lastBuyPrice := lastMarketDepthReturn.Tick.Bids[0][0]
				nowBuyPrice := marketDepthReturn1.Tick.Bids[0][0]
				diff := nowBuyPrice - lastBuyPrice
				if diff > 0 && diff/lastBuyPrice > 0.015 {
					//大幅上涨
					utils.PrintLog(true, "req_test"+event.Symbol, fmt.Sprintf("%v", time.Now().Unix()/int64(min*60)), fmt.Sprintf("%v %v %v %v", event.Symbol, nowBuyPrice,
						"大幅上涨",
						utils.DivRound(decimal.NewFromFloat(diff/lastBuyPrice*100), 2)+"%"))

					//buysell(event.Symbol, nowBuyPrice, 1.00)
				}
				if diff < 0 && diff/lastBuyPrice < -0.015 {
					//大幅下跌
					utils.PrintLog(true, "req_test"+event.Symbol, fmt.Sprintf("%v", time.Now().Unix()/int64(min*60)), fmt.Sprintf("%v %v %v %v", event.Symbol, nowBuyPrice,
						"大幅下跌",
						utils.DivRound(decimal.NewFromFloat(diff/lastBuyPrice*100), 2)+"%"))

					buysell(event.Symbol, nowBuyPrice, 0.98)
				}
			}
		}
		var errHandler = func(err error) {
			if err != nil {
				wserr = err
				oobj.bookTickers = sync.Map{}

				//utils.PrintLog(api.PrintLog, "Binance.getAccountBalance", fmt.Sprintf("%v", time.Now().Unix()/(60*10)), err.Error())
				//utils.SendDingTalkRobit(true, "oneplat", "Binance.WsBookTicker", fmt.Sprintf("%v", time.Now().Unix()/(60*10)), err.Error())
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
					//utils.PrintLog(api.PrintLog, "Binance.getAccountBalance", fmt.Sprintf("%v", time.Now().Unix()/(60*10)), wserr.Error())
					//utils.SendDingTalkRobit(true, "oneplat", "Binance.WsBookTicker", fmt.Sprintf("%v", time.Now().Unix()/(60*10)), wserr.Error())
				}
			}
			time.Sleep(time.Second)
		}
	}(plat)

	select {}
}
