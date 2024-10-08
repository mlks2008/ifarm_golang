package main

import (
	"components/database/redis2"
	"components/log/log"
	"components/message"
	"components/myconfig"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/shopspring/decimal"
	"time"
)

func report300U() {
	var binanceCli = NewCexBinanceCli(myconfig.GConfig.Binance.ApiKey, myconfig.GConfig.Binance.SecretKey)

	//当月转出历史
	var earnOuts, err1 = binanceCli.CurMonthTransferHistory()
	if err1 != nil {
		log.Logger.Error("CurMonthTransferHistory", err1)
		logmsg := fmt.Sprintf("CurMonthTransferHistory:%v", err1.Error())
		message.SendDingTalkRobit(true, myconfig.GConfig.Project, "CurMonthTransferHistory", fmt.Sprintf("%v", time.Now().Unix()/(60*60*24)), logmsg)
	}

	var getEarnOut = func(date time.Time) decimal.Decimal {
		var amount = decimal.Zero
		for _, out := range earnOuts {
			if out.Status == "SUCCESS" && out.Type == binance.SubAccountTransferTypeTransferOut && (out.Asset == "FDUSD" || out.Asset == "USDT" || out.Asset == "USDC") {
				if time.Unix(out.Time/1000, 0).Format("2006-01-02") == date.Format("2006-01-02") {
					var t, _ = decimal.NewFromString(out.Qty)
					amount = amount.Add(t)
				}
			}
		}
		return amount
	}

	var getTr = func(lastDate time.Time, date time.Time, count int, totalEarn decimal.Decimal) (string, decimal.Decimal) {
		var redis = redis2.NewRedisCli(myconfig.GConfig.Redis.Host, myconfig.GConfig.Redis.Password, myconfig.GConfig.Redis.DB)

		//前一天余额
		lastKey := fmt.Sprintf("binance-%v", lastDate.Format("2006-01-02"))
		lastBal, err := redis.GetDecimal(lastKey)
		if err != nil {
			logmsg := fmt.Sprintf("redis.GetDecimal:%v", err.Error())
			message.SendDingTalkRobit(true, myconfig.GConfig.Project, "redis.GetDecimal", fmt.Sprintf("%v", time.Now().Unix()/(60*60*24)), logmsg)
			return "", totalEarn
		}

		//当天余额
		key := fmt.Sprintf("binance-%v", date.Format("2006-01-02"))
		var bal decimal.Decimal
		if date.Format("2006-01-02") == time.Now().Format("2006-01-02") {
			var bals, err = binanceCli.SpotBalances()
			if err != nil {
				logmsg := fmt.Sprintf("Binance.SpotBalances:%v", err.Error())
				message.SendDingTalkRobit(true, myconfig.GConfig.Project, "Binance.SpotBalances", fmt.Sprintf("%v", time.Now().Unix()/(60*60*24)), logmsg)
				return "", totalEarn
			}
			for _, b := range bals.Balances {
				if b.Asset == "FDUSD" || b.Asset == "USDT" || b.Asset == "USDC" {
					var t, _ = decimal.NewFromString(b.Free)
					bal = bal.Add(t)
				}
			}
			redis.SetEX(key, bal.String(), 60*24*3600*time.Second)
		} else {
			bal, err = redis.GetDecimal(key)
			if err != nil {
				logmsg := fmt.Sprintf("redis.GetDecimal:%v", err.Error())
				message.SendDingTalkRobit(true, myconfig.GConfig.Project, "redis.GetDecimal", fmt.Sprintf("%v", time.Now().Unix()/(60*60*24)), logmsg)
				return "", totalEarn
			}
		}

		var classVal = "bgcolor1"
		if count%2 != 0 {
			classVal = "bgcolor0"
		}

		earn := bal.Sub(lastBal)
		earnOut := getEarnOut(date)
		totalEarn = totalEarn.Add(earn).Add(earnOut)

		//日期，收益，累计收益，帐户余额
		var tr = fmt.Sprintf(`<tr class=%v>
		       	<td class='padding'>%v</td> 
		       	<td class='padding'>%v</td>
				<td class='padding'>%v</td>
		       	<td class='padding'>%v</td>
		       	<td class='padding'>%v</td>
		       </tr>`, classVal, date.Format("01-02"), earn.String(), earnOut.String(), totalEarn.String(), bal.String())
		return tr, earn.Add(earnOut)
	}

	// 计算每天收益
	var count int
	var totalEarn = decimal.Zero
	var startDate = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	var curDate = startDate.AddDate(0, 0, time.Now().Day()-1)

	var trList = make([]string, 0)
	for date := startDate; date.Before(curDate) || date.Equal(curDate); date = date.AddDate(0, 0, 1) {
		var lastDate = date.AddDate(0, 0, -1)
		var tr, earn = getTr(lastDate, date, count, totalEarn)
		trList = append(trList, tr)
		totalEarn = totalEarn.Add(earn)
		count++
	}

	// 形成table
	var trs string
	for i := len(trList) - 1; i >= 0; i-- {
		trs += trList[i]
	}
	location, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Unix(time.Now().Unix(), 0).In(location)
	var nowDate = t.Format("2006-01-02 15:04:05")
	var table = fmt.Sprintf(`<div style='text-align: left'><h3>手动每日收益报表 %v</h3><h4>日收益：当天余额-前日余额，月收益：日收益+收益转出</h4></div>
	       <table>
	       <thead>
	       <tr class='bgcolor0'>
	           	<th class='padding'>日期</th>
	           	<th class='padding'>日收益 $</th>
				<th class='padding'>收益转出 $</th>
	           	<th class='padding'>月收益 $</th>
	           	<th class='padding'>帐户余额 $</th>
	       </tr>
	       </thead>
	       <tbody>%v</tbody></table>`, nowDate, trs)

	//发送邮件
	var emailHtml = fmt.Sprintf(`<style>
	   table {table-layout: auto; text-align: center; font-size: 14px;border: 1px #d3d3d3 solid;}
	   td {white-space: nowrap;overflow: hidden;text-overflow: ellipsis;}
	   .padding {padding: 6px 6px;}
	   .bgcolor0{background-color: #f2f2f2;}
	   .bgcolor1{background-color: #ffffff;}
	   .reason{text-align: left;}
	   </style>
		%v`, table)
	var subject = fmt.Sprintf(`Farm Manual Daily Report %v`, nowDate)
	var err = new(EmailHelper).SendEmail(subject, emailHtml)
	if err != nil {
		log.Logger.Error("SendEmail", err)
		logmsg := fmt.Sprintf("SendEmail:%v", err.Error())
		message.SendDingTalkRobit(true, myconfig.GConfig.Project, "SendEmail", fmt.Sprintf("%v", time.Now().Unix()/(60*60*24)), logmsg)
	}
	log.Logger.Debug("report300U ok")
}
