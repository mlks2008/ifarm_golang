package main

import (
	"components/database/redis2"
	"components/log/log"
	"components/message"
	"components/myconfig"
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

func reportFilDown_oneplat() {
	reportDown1("oneplat", "fil")
}

func reportDogeDown_oneplat() {
	reportDown1("oneplat", "doge")
}

func reportDogeDown_mainapi() {
	reportDown1("mainapi", "doge")
}

func reportDown1(robot string, coin string) {
	var getTr = func(lastDate time.Time, date time.Time, count int, totalEarn decimal.Decimal) (string, decimal.Decimal) {
		var redis = redis2.NewRedisCli(myconfig.GConfig.Redis.Host, myconfig.GConfig.Redis.Password, myconfig.GConfig.Redis.DB)

		//当天余额
		key := fmt.Sprintf("%v-%vdown-%v", robot, coin, date.Format("2006-01-02"))
		bal, err := redis.GetDecimal(key)
		if err != nil {
			logmsg := fmt.Sprintf("redis.GetDecimal:%v", err.Error())
			message.SendDingTalkRobit(true, myconfig.GConfig.Project, "redis.GetDecimal", fmt.Sprintf("%v", time.Now().Unix()/(60*60*24)), logmsg)
			return "", totalEarn
		}

		//前一天余额（或者叫上一次余额）
		var lastBal decimal.Decimal
		if bal != decimal.Zero {
			for {
				lastKey := fmt.Sprintf("%v-%vdown-%v", robot, coin, lastDate.Format("2006-01-02"))
				lastBal, err = redis.GetDecimal(lastKey)
				if err != nil {
					logmsg := fmt.Sprintf("redis.GetDecimal:%v", err.Error())
					message.SendDingTalkRobit(true, myconfig.GConfig.Project, "redis.GetDecimal", fmt.Sprintf("%v", time.Now().Unix()/(60*60*24)), logmsg)
					return "", totalEarn
				}
				if lastBal == decimal.Zero {
					if lastDate.Day() == 1 {
						lastBal = bal
						break
					} else {
						lastDate = lastDate.AddDate(0, 0, -1)
					}
				} else {
					break
				}
			}
		}

		//当天没有记录，说明没有收益，使用上一次的
		if bal == decimal.Zero {
			bal = lastBal
		}

		var classVal = "bgcolor1"
		if count%2 != 0 {
			classVal = "bgcolor0"
		}

		earn := bal.Sub(lastBal)
		totalEarn = totalEarn.Add(earn)

		//日期，收益，累计收益，帐户余额
		var tr = fmt.Sprintf(`<tr class=%v>
		       	<td class='padding'>%v</td> 
		       	<td class='padding'>%v</td>
				<td class='padding'>%v</td>
		       	<td class='padding'>%v</td>
		       </tr>`, classVal, date.Format("01-02"), earn.String(), totalEarn.String(), bal.String())
		return tr, earn
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
	var table = fmt.Sprintf(`<div style='text-align: left'><h3>%v每日收益报表 %v</h3><h4>日收益：当天余额-前日余额</h4></div>
	       <table>
	       <thead>
	       <tr class='bgcolor0'>
	           	<th class='padding'>日期</th>
	           	<th class='padding'>日收益 $</th>
	           	<th class='padding'>月收益 $</th>
	           	<th class='padding'>帐户余额 $</th>
	       </tr>
	       </thead>
	       <tbody>%v</tbody></table>`, strings.ToUpper(coin), nowDate, trs)

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
	var subject = fmt.Sprintf(`%v - Farm %v Daily Report %v`, robot, strings.ToUpper(coin), nowDate)
	var err = new(EmailHelper).SendEmail(subject, emailHtml)
	if err != nil {
		log.Logger.Error("SendEmail", err)
		logmsg := fmt.Sprintf("SendEmail:%v", err.Error())
		message.SendDingTalkRobit(true, myconfig.GConfig.Project, "SendEmail", fmt.Sprintf("%v", time.Now().Unix()/(60*60*24)), logmsg)
	}
	log.Logger.Debug("report ok")
}
