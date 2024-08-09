package api

import (
	"flag"
)

var Proxy = flag.Bool("proxy", false, "[可选]")

var (
	OnePlatStop        bool    = false
	PrintLog           bool    = true //默认开启日志
	CheckBalance       bool    = true //默认开启下单余额检测
	OnePlatUsdtAmount1 float64 = 20
	OnePlatUsdtAmount2 float64 = 20
	OnePlatUsdtAmount3 float64 = 20
)

func init() {
	////机器人数
	//for i := 0; i < 1; i++ {
	//	OnePlatSymbols = append(OnePlatSymbols, "BTC/USDT")
	//}
}

func GetPrintInfo() interface{} {
	data := struct {
		OnePlatStop        bool    `json:"OnePlatStop"`
		OnePlatUsdtAmount1 float64 `json:"OnePlatUsdtAmount1"`
		OnePlatUsdtAmount2 float64 `json:"OnePlatUsdtAmount2"`
		OnePlatUsdtAmount3 float64 `json:"OnePlatUsdtAmount3"`
	}{
		OnePlatStop:        OnePlatStop,
		OnePlatUsdtAmount1: OnePlatUsdtAmount1,
		OnePlatUsdtAmount2: OnePlatUsdtAmount2,
		OnePlatUsdtAmount3: OnePlatUsdtAmount3,
	}
	return data
}
