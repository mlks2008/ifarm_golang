package services

import (
	"fmt"
	"github.com/shopspring/decimal"
	"goarbitrage/pkg/utils"
	"strconv"
	"strings"
	"time"
)

// 每次执行前USDT成本（每次完成套利后写入的值）
func (this *DCAService) RunGetUsdtCost() (decimal.Decimal, decimal.Decimal) {
	cost := utils.ReadFile(fmt.Sprintf("files/dcaservice.%v%v.cost", this.symbolBase, this.symbolQuote))
	cost = strings.Replace(cost, "\n", "", -1)
	if cost == "" {
		v1, v2 := this.InitialUsdtCost()
		this.RunSetUsdtCost(v1.String(), v2.String())
		return this.InitialUsdtCost()
	} else {
		v := strings.Split(cost, ",")
		costQuote, _ := decimal.NewFromString(v[0])
		costostBase, _ := decimal.NewFromString(v[1])
		return costQuote, costostBase
	}
}

// 套利完成后保成USDT成本
func (this *DCAService) RunSetUsdtCost(usdtbalance string, btcbalance string) {
	utils.UpdateFile(fmt.Sprintf("files/dcaservice.%v%v.cost", this.symbolBase, this.symbolQuote), usdtbalance+","+btcbalance)
}

// 首次买入价格
func (this *DCAService) GetFirstBuyPrice() decimal.Decimal {
	val := utils.ReadFile(fmt.Sprintf("files/dcaservice.%v%v.firstbuyprice", this.symbolBase, this.symbolQuote))
	val = strings.Replace(val, "\n", "", -1)
	if val == "" {
		return decimal.Zero
	} else {
		decVal, _ := decimal.NewFromString(val)
		return decVal
	}
}

// 保存首次买入价格
func (this *DCAService) SetFirstBuyPrice(buyprice string) {
	utils.UpdateFile(fmt.Sprintf("files/dcaservice.%v%v.firstbuyprice", this.symbolBase, this.symbolQuote), buyprice)
}

// 买入后最低价格
func (this *DCAService) GetMinPrice() decimal.Decimal {
	val := utils.ReadFile(fmt.Sprintf("files/dcaservice.%v%v.minprice", this.symbolBase, this.symbolQuote))
	val = strings.Replace(val, "\n", "", -1)
	if val == "" {
		return decimal.Zero
	} else {
		decVal, _ := decimal.NewFromString(val)
		return decVal
	}
}

// 保存买入后最低价格
func (this *DCAService) SetMinPrice(minprice string) {
	utils.UpdateFile(fmt.Sprintf("files/dcaservice.%v%v.minprice", this.symbolBase, this.symbolQuote), minprice)
}

// 最近一次买入时间
func (this *DCAService) GetLastBuyTime() int64 {
	val := utils.ReadFile(fmt.Sprintf("files/dcaservice.%v%v.lastbuytime", this.symbolBase, this.symbolQuote))
	val = strings.Replace(val, "\n", "", -1)
	if val == "" {
		return time.Now().Unix() - 10*60
	} else {
		val, _ := strconv.ParseInt(val, 10, 64)
		return val
	}
}

// 保存最近一次买入时间
func (this *DCAService) SetLastBuyTime(lastbuytime int64) {
	utils.UpdateFile(fmt.Sprintf("files/dcaservice.%v%v.lastbuytime", this.symbolBase, this.symbolQuote), fmt.Sprintf("%v", lastbuytime))
}
