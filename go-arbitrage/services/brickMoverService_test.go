/**
 * brickMoverService.go
 * ============================================================================
 * 2、多交易所搬砖套利
 * ============================================================================
 * author: peter.wang
 */

package services

import (
	"github.com/spf13/viper"
	_ "goarbitrage/api/router"
	"goarbitrage/internal/plat"
	"testing"
)

func init() {
	viper.SetConfigFile("../configs/hotbit_dev.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	plat.Register()
}

//交易所间搬砖
func Test_TWOPlat_Start(t *testing.T) {
	//plat1 := plat.Get(plat.P_Binance)
	//plat2 := plat.Get(plat.P_KuCoin)
	//
	//for {
	//	symbols, _ := plat2.GetSymbols()
	//	symbols = &model.SymbolsReturn{
	//		Data: []model.SymbolsData{
	//			model.SymbolsData{Symbol: "TRX-USDT"},
	//			//model.SymbolsData{Symbol: "BTC-USDT"},
	//		},
	//	}
	//	for _, tmp := range symbols.Data {
	//		symbol := strings.Replace(tmp.Symbol, "-", "/", 1)
	//
	//		//如果plat1的买一 > plat2的卖一，计算价差高于两边手续费比例，则plat1卖，plat2买
	//		//如果plat2的买一 > plat1的卖一，计算价差高于两边手续费比例，则plat2卖，plat1买
	//		for i := 0; i < 5; i++ {
	//			p1, err := plat1.GetMarketDepth(symbol)
	//			if err != nil {
	//				//fmt.Println(fmt.Sprintf("p1 err %v", err))
	//				break
	//			}
	//			if p1.Tick.Bids[0][0] == 0 {
	//				break
	//			}
	//
	//			p2, err := plat2.GetMarketDepth(symbol)
	//			if err != nil {
	//				fmt.Println(fmt.Sprintf("p2 err %v", err))
	//				break
	//			}
	//			if p2.Tick.Bids[0][0] == 0 {
	//				break
	//			}
	//
	//			if p1.Tick.Bids[0][0] > p2.Tick.Asks[0][0] {
	//				rate := (p1.Tick.Bids[0][0] - p2.Tick.Asks[0][0]) / p1.Tick.Bids[0][0] * 1000
	//				if rate > 0 {
	//					fmt.Println(fmt.Sprintf("%v p1卖价:%v,p2买价:%v, 价差(千分比):%v", symbol, p1.Tick.Bids[0][0], p2.Tick.Asks[0][0], rate))
	//					fmt.Println("")
	//					break
	//				}
	//			} else if p2.Tick.Bids[0][0] > p1.Tick.Asks[0][0] {
	//				rate := (p2.Tick.Bids[0][0] - p1.Tick.Asks[0][0]) / p2.Tick.Bids[0][0] * 1000
	//				if rate > 0 {
	//					fmt.Println(fmt.Sprintf("%v p2卖价:%v,p1买价:%v, 价差(千分比):%v", symbol, p2.Tick.Bids[0][0], p1.Tick.Asks[0][0], rate))
	//					fmt.Println("")
	//					break
	//				}
	//			}
	//		}
	//	}
	//	time.Sleep(time.Second * 1)
	//}
}
