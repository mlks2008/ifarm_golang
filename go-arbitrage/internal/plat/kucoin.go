/**
 * kucoin.go
 * ============================================================================
 * https://github.com/Kucoin/kucoin-go-sdk
 * ============================================================================
 * author: peter.wang
 */

package plat

import (
	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"goarbitrage/internal/model"
	"strconv"
	"strings"
)

type KuCoin struct {
	Base
	client *kucoin.ApiService
}

func NewKuCoin() *KuCoin {
	obj := &KuCoin{}
	obj.client = kucoin.NewApiService(
		kucoin.ApiBaseURIOption(viper.GetString("plats.kucoin.apiaddr")),
		kucoin.ApiKeyOption(viper.GetString("plats.kucoin.access_key")),
		kucoin.ApiSecretOption(viper.GetString("plats.kucoin.secret_key")),
		kucoin.ApiPassPhraseOption(viper.GetString("plats.kucoin.passphrase")),
	)
	return obj
}

// 灵活测试子帐号
func (this *KuCoin) SetKey(symbol string, apiKey, secretKey string, initial map[string]model.Balance) {

}

func (this *KuCoin) GetInitialInput(symbol string) map[string]model.Balance {
	return nil
}

func (this *KuCoin) PlatCode() string {
	return P_KuCoin
}

// KuCoin交易对格式:ETH-USDT
func (this *KuCoin) FormatSymbol(symbol string) string {
	return strings.Replace(symbol, "/", "-", 1)
}

func (this *KuCoin) GetMarketDepth(symbol string) (model.MarketDepthReturn, error) {
	symbol = this.FormatSymbol(symbol)

	rsp, err := this.client.AggregatedPartOrderBook(symbol, 20)
	if err != nil {
		return model.MarketDepthReturn{}, err
	}
	c := &kucoin.PartOrderBookModel{}
	if err := rsp.ReadData(c); err != nil {

		return model.MarketDepthReturn{}, err

	} else {
		var marketDepthReturn = model.MarketDepthReturn{}

		if len(c.Bids) > 0 {
			price, _ := strconv.ParseFloat(c.Bids[0][0], 10)
			num, _ := strconv.ParseFloat(c.Bids[0][1], 10)

			marketDepthReturn.Tick.Bids = make([][]float64, 1)
			marketDepthReturn.Tick.Bids[0] = make([]float64, 2)
			marketDepthReturn.Tick.Bids[0][0] = price //价格
			marketDepthReturn.Tick.Bids[0][1] = num   //数量
		}

		if len(c.Asks) > 0 {
			price, _ := strconv.ParseFloat(c.Asks[0][0], 10)
			num, _ := strconv.ParseFloat(c.Asks[0][1], 10)

			marketDepthReturn.Tick.Asks = make([][]float64, 1)
			marketDepthReturn.Tick.Asks[0] = make([]float64, 2)
			marketDepthReturn.Tick.Asks[0][0] = price //价格
			marketDepthReturn.Tick.Asks[0][1] = num   //数量
		}

		return marketDepthReturn, nil
	}
}

func (this *KuCoin) AveragePrice(symbol string) (decimal.Decimal, error) {
	return decimal.Zero, nil
}

func (this *KuCoin) GetAccountBalance(symbol string, realtime bool) (map[string]model.Balance, error) {
	return nil, nil
}

func (this *KuCoin) OrdersPlace(symbol string, price string, quantity string, side string) (string, error) {
	return "", nil
}

func (this *KuCoin) OrdersPlace2(clientsymbol string, symbol string, price string, quantity string, side string) (string, error) {
	return "", nil
}

func (this *KuCoin) CancelOrder(symbol string, orderid string) (bool, error) {
	return false, nil
}

func (this *KuCoin) GetOrderStatus(symbol string, orderid string) (bool, error) {
	return false, nil
}

func (this *KuCoin) GetOrderStatus2(clientsymbol string, symbol string, orderid string) (bool, error) {
	return false, nil
}

func (this *KuCoin) OpenOrders(symbol string) ([]*model.Order, error) {
	return nil, nil
}

func (this *KuCoin) OpenOrders2(clientsymbol string, symbol string) ([]*model.Order, error) {
	return nil, nil
}

func (this *KuCoin) CancelOpenOrders(symbol string) error {
	return nil
}

func (this *KuCoin) GetTradeFee(symbol string) (makerfee decimal.Decimal, takerfee decimal.Decimal, err error) {
	return
}
