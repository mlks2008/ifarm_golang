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
	reportDogeDown()

	c := cron.New(cron.WithSeconds())
	//_, err := c.AddFunc("0 0 */8 * * *", report300U)
	_, err := c.AddFunc("0 0 23 * * *", report300U)
	if err != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}

	_, err = c.AddFunc("0 0 6 * * *", reportDogeDown)
	if err != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}
	_, err = c.AddFunc("59 59 23 * * *", reportDogeDown)
	if err != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}

	c.Start()

	select {}
}
