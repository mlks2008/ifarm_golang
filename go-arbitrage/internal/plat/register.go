package plat

import (
	"fmt"
	"github.com/shopspring/decimal"
	"goarbitrage/internal/model"
)

var (
	regPlat map[string]Plat
)

func Register() {
	regPlat = make(map[string]Plat)
	register(NewKuCoin())
	register(NewBinance())
}

func register(p Plat) {
	if _, ok := regPlat[p.PlatCode()]; ok {
		panic(fmt.Sprintf("register plat [%v] duplicate", p.PlatCode()))
	}
	regPlat[p.PlatCode()] = p
}

type Plat interface {
	PlatCode() string
	FormatSymbol(symbol string) string
	GetSymbols() (*model.SymbolsReturn, error)
	GetScope() (map[string]string, error) //支持的交易区
	GetMarketDepth(symbol string) (model.MarketDepthReturn, error)
	AveragePrice(symbol string) (decimal.Decimal, error)
	GetAccountBalance(symbol string, realtime bool) (map[string]model.Balance, error)
	SetKey(symbol string, apiKey, secretKey string, initial map[string]model.Balance)
	GetInitialInput(symbol string) map[string]model.Balance //初始投入
	OrdersPlace(symbol string, price string, quantity string, side string) (string, error)
	OrdersPlace2(clientsymbol string, symbol string, price string, quantity string, side string) (string, error)
	CancelOrder(symbol string, orderid string) (bool, error)
	GetOrderStatus(symbol string, orderid string) (bool, error)
	GetOrderStatus2(clientsymbol string, symbol string, orderid string) (bool, error)
	OpenOrders(symbol string) ([]*model.Order, error)
	OpenOrders2(clientsymbol string, symbol string) ([]*model.Order, error)
	CancelOpenOrders(symbol string) error
	GetTradeFee(symbol string) (makerfee decimal.Decimal, takerfee decimal.Decimal, err error)
}

func Get(platcode string) Plat {
	if p, ok := regPlat[platcode]; ok {
		return p
	} else {
		panic(fmt.Sprintf("plat [%v] not register", platcode))
	}
}
