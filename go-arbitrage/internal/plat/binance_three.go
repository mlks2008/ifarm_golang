package plat

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"goarbitrage/internal/model"
	"strconv"
	"strings"
)

//支持的交易区
func (this *Binance) GetScope() (map[string]string, error) {
	return map[string]string{"BTC": "", "ETH": "", "USDT": "", "EOS": ""}, nil
}

//所有交易对
func (this *Binance) GetSymbols() (*model.SymbolsReturn, error) {
	this.wait(this.PlatCode(), "GetSymbols", 1, "symbol")

	var client *binance.Client
	for _, c := range this.client {
		client = c
	}

	resp, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	symbolsReturn := &model.SymbolsReturn{}
	for _, s := range resp.Symbols {
		//if s.Symbol == "CFXUSDT" {
		//	fmt.Println("")
		//}
		//交易价格精度
		var pricePrecision int32
		var amountPrecision int32
		var minQuoteOrder float64
		var maxNumOrders float64
		for _, filter := range s.Filters {
			if filter["filterType"] == "PRICE_FILTER" {
				pricePrecision = int32(strings.Index(fmt.Sprintf("%v", filter["tickSize"]), "1") - 1)
				if pricePrecision < 0 {
					pricePrecision = 0
				}
			}
			if filter["filterType"] == "LOT_SIZE" {
				amountPrecision = int32(strings.Index(fmt.Sprintf("%v", filter["stepSize"]), "1") - 1)
				if amountPrecision < 0 {
					amountPrecision = 0
				}
			}
			if filter["filterType"] == "MIN_NOTIONAL" {
				minQuoteOrder, _ = strconv.ParseFloat(filter["minNotional"].(string), 64)
			}
			if filter["filterType"] == "MAX_NUM_ORDERS" {
				maxNumOrders, _ = filter["maxNumOrders"].(float64)
			}

		}

		var state bool
		if s.Status == "TRADING" {
			state = true
		}

		symbolsReturn.Data = append(symbolsReturn.Data, model.SymbolsData{
			BaseCurrency:    s.BaseAsset,
			BasePrecision:   int32(s.BaseAssetPrecision), //币种精度
			QuoteCurrency:   s.QuoteAsset,
			QuotePrecision:  int32(s.QuoteAssetPrecision), //报价币精度
			PricePrecision:  pricePrecision,               //交易对精度
			AmountPrecision: amountPrecision,              //下单数精度
			MinQuoteOrder:   minQuoteOrder,                //最小下单量
			MaxNumOrders:    maxNumOrders,                 //最大下单量
			Symbol:          s.Symbol,
			State:           state,
		})
	}

	return symbolsReturn, nil
}
