package message

import (
	"fmt"
	"time"
)

var chLogLog = make(chan Log, 10000)

func init() {
	var cashLog = make(map[string]bool)

	go func() {

		for {

			log := <-chLogLog

			if log.Enable == false {
				continue
			}

			key := fmt.Sprintf("%v_%v", log.Typ, log.TxId)

			if _, ok := cashLog[key]; ok {
				continue
			} else {
				//重置
				if len(cashLog) >= 1000 {
					cashLog = make(map[string]bool)
				}
				cashLog[key] = true
				fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " ", log.Content)
			}

		}
	}()
}

func PrintLog(printlog bool, typ, txid, content string) {
	log := Log{}
	log.Enable = printlog
	log.Typ = typ
	log.TxId = txid
	log.Content = content
	chLogLog <- log
}
