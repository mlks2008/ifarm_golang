package main

import (
	"components/log/log"
	"fmt"
	"github.com/robfig/cron/v3"
	"math/rand"
	"time"
)

var configFile string
var cmd string

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	log.InitLogger("./", "goreport", true)

	report300U()
	reportDogeDown_oneplat()
	reportDogeDown_mainapi()

	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("0 0 6 * * *", report300U)
	if err != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}

	_, err = c.AddFunc("0 0 6 * * *", reportDogeDown_oneplat)
	if err != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}

	_, err = c.AddFunc("0 0 6 * * *", reportDogeDown_mainapi)
	if err != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}

	c.Start()

	select {}
}
