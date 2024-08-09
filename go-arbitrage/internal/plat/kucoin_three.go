package plat

import (
	"fmt"
	"github.com/Kucoin/kucoin-go-sdk"
	"goarbitrage/internal/model"
)

//支持的交易区
func (this *KuCoin) GetScope() (map[string]string, error) {
	return map[string]string{"BTC": "", "ETH": "", "USDT": "", "EOS": ""}, nil
}

//所有交易对
func (this *KuCoin) GetSymbols() (*model.SymbolsReturn, error) {
	rsp, err := this.client.Symbols("")
	if err != nil {
		return nil, err
	}

	l := kucoin.SymbolsModel{}
	if err := rsp.ReadData(&l); err != nil {
		fmt.Println(err)
	}

	symbolsReturn := &model.SymbolsReturn{}
	for _, s := range l {
		symbolsReturn.Data = append(symbolsReturn.Data, model.SymbolsData{
			BaseCurrency:  s.BaseCurrency,
			QuoteCurrency: s.QuoteCurrency,
			Symbol:        s.Symbol,
			State:         s.EnableTrading,
		})
	}

	return symbolsReturn, nil
}
