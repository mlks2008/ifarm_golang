package plat

import (
	"encoding/json"
	"github.com/pkg/errors"
	"goarbitrage/internal/model"
	"goarbitrage/internal/platsdk/hoosdk"
	"log"
	"strconv"
	"strings"
)

//支持的交易区
func (this *Hoo) GetScope() (map[string]string, error) {
	return map[string]string{"BTC": "", "ETH": "", "USDT": ""}, nil
}

//所有交易对
func (this *Hoo) GetSymbols() (*model.SymbolsReturn, error) {
	symbolsReturn := &model.SymbolsReturn{}

	body, err := this.client.GetSymbols()
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal([]byte(body), &symbolsReturn)

	for i, _ := range symbolsReturn.Data {
		symbol := strings.Split(symbolsReturn.Data[i].Symbol, "-")
		symbolsReturn.Data[i].BaseCurrency = symbol[0]
		symbolsReturn.Data[i].QuoteCurrency = symbol[1]
		symbolsReturn.Data[i].State = true
	}

	return symbolsReturn, nil
}

//深度
func (this *Hoo) GetMarketDepth(symbol string) (*model.MarketDepthReturn, error) {
	body, err := this.client.GetMarketDepth(symbol)
	if err != nil {
		log.Println(err)
	}
	res := &hoosdk.DepthReturn{}
	json.Unmarshal([]byte(body), &res)
	if res.Code != 0 {
		return nil, errors.New(res.Msg)
	}

	if len(res.Data.Asks) > 0 && len(res.Data.Bids) > 0 {
		marketDepthReturn := &model.MarketDepthReturn{}
		marketDepthReturn.Tick.Asks = make([][]float64, 1)
		marketDepthReturn.Tick.Asks[0] = make([]float64, 2)
		marketDepthReturn.Tick.Asks[0][0], _ = strconv.ParseFloat(res.Data.Asks[0].Price, 10)
		marketDepthReturn.Tick.Asks[0][1], _ = strconv.ParseFloat(res.Data.Asks[0].Quantity, 10)
		marketDepthReturn.Tick.Bids = make([][]float64, 1)
		marketDepthReturn.Tick.Bids[0] = make([]float64, 2)
		marketDepthReturn.Tick.Bids[0][0], _ = strconv.ParseFloat(res.Data.Bids[0].Price, 10)
		marketDepthReturn.Tick.Bids[0][1], _ = strconv.ParseFloat(res.Data.Bids[0].Quantity, 10)
		return marketDepthReturn, nil
	} else {
		return nil, errors.New("empty")
	}
}
