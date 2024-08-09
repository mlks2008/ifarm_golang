package model

import "github.com/shopspring/decimal"

// ---------所有交易对-----
type SymbolsData struct {
	BaseCurrency    string  `json:"baseCurrency"` // 基础币种
	BasePrecision   int32   `json:"baseCurrency"`
	QuoteCurrency   string  `json:"quoteCurrency"` // 计价币种/报价币种
	QuotePrecision  int32   `json:"quotePrecision"`
	PricePrecision  int32   `json:"pricePrecision"`  //交易对精度
	AmountPrecision int32   `json:"amountPrecision"` //下单数精度
	MinQuoteOrder   float64 `json:"minQuoteOrder"`   ///最小下单量(报价币种)
	MaxNumOrders    float64 `json:"maxNumOrders"`    //最大订单数
	Symbol          string  `json:"symbol"`          // 交易对
	State           bool
}

type SymbolsReturn struct {
	Code int           `json:"code"`
	Data []SymbolsData `json:"data"`
}

// ---------所有交易对深度-----
type MarketDepth struct {
	Bids [][]float64 `json:"bids"` // 买盘, [price(成交价), amount(成交量)], 按price降序排列
	Asks [][]float64 `json:"asks"` // 卖盘, [price(成交价), amount(成交量)], 按price升序排列
}

type MarketDepthReturn struct {
	Code int         `json:"code"` // 请求状态, ok或者error
	Msg  string      `json:"msg"`
	Tick MarketDepth `json:"tick"` // Depth数据
}

type SchedulingIndex struct {
	Step1Index int `json:"step_1_index"`
	Step2Index int `json:"step_2_index"`
	Step3Index int `json:"step_3_index"`
}

type Balance struct {
	Asset   string          `json:"asset"`
	Total   decimal.Decimal `json:"total"`
	Free    decimal.Decimal `json:"free"`
	Locked  decimal.Decimal `json:"locked"`
	AddTime int64
}

type Order struct {
	Symbol           string `json:"symbol"`
	OrderID          int64  `json:"orderId"`
	Price            string `json:"price"`
	OrigQuantity     string `json:"origQty"`
	ExecutedQuantity string `json:"executedQty"`
	Side             string `json:"side"`
	Time             int64  `json:"time"`
}
