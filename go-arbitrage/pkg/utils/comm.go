package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

func DivRound(amount decimal.Decimal, precision int32) string {
	return amount.Mul(decimal.New(1, precision)).Floor().Div(decimal.New(1, precision)).String()
}

func UseTime(startTime int64) string {
	//运行时间
	runTotalUseTime := time.Now().Unix() - startTime
	runTotalFormatTime, _ := decimal.NewFromFloat(float64(runTotalUseTime)/86400.0).DivRound(decimal.NewFromInt(1), 2).Float64()
	runTotalFormatUnit := "d"
	if runTotalFormatTime < 1 {
		runTotalFormatTime, _ = decimal.NewFromFloat(float64(runTotalUseTime)/3600.0).DivRound(decimal.NewFromInt(1), 2).Float64()
		runTotalFormatUnit = "h"
	}
	return fmt.Sprintf("%v%v", runTotalFormatTime, runTotalFormatUnit)
}
