package message

import (
	"testing"
	"time"
)

func Test_SendDingTalkRobit(t *testing.T) {
	//SendDingTalkRobit(true, "all", "all", "txid", "all")
	//SendDingTalkRobit(true, "token", "token", "txid", "token")
	//SendDingTalkRobit(true, "node", "node", "txid", "node")
	SendDingTalkRobit(true, "oneplat", "oneplat", "txid", "oneplat")
	SendDingTalkRobit(true, "alert", "alert", "txid", "alert")
	time.Sleep(time.Second * 10)
}
