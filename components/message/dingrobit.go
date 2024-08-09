package message

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Log struct {
	Enable  bool
	Robot   string
	Typ     string
	TxId    string
	Content string
}

var chLog = make(chan Log, 10000)

func init() {
	var cashLog = make(map[string]bool)

	go func() {
		for {

			log := <-chLog

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
			}

			var secret string

			var accesstoken string

			var data string

			if strings.ToLower(log.Robot) == "all" {
				//全机器人
				secret = "SEC1ee94f53b28f4452a7da8d6abee1068d173538e4259ca933a713518ed8b0c81f"
				accesstoken = "454d64566f9bddf9051351051af8c9dd42e4a1858854e9620d460b1d0e725fc6"
				data = fmt.Sprintf(`{"msgtype": "text","text": {"content": "%v"}}`, fmt.Sprintf("%v\nmd:%v \t%v", log.Content, log.Robot, time.Now().Format("01-02 15:04")))
			} else if strings.ToLower(log.Robot) == "token" {
				//牛机器人
				secret = "SEC42a41f4480883d98b9c41e402a58a578599bff98ab66f96e465e0ed90c68f8e0"
				accesstoken = "ee52b626cc4fb9d7d81dc69e8759349c50e0f10546cd4d9a08d13281e8d5ce9b"
				data = fmt.Sprintf(`{"msgtype": "text","text": {"content": "%v"}}`, fmt.Sprintf("%v\nmd:%v \t%v", log.Content, log.Robot, time.Now().Format("01-02 15:04")))
			} else if strings.ToLower(log.Robot) == "node" {
				//w机器人
				secret = "SEC73d9e7dbe6886e8a571864c171ddad5c73c8338d35c24d0ff0093f4c95b97510"
				accesstoken = "08fefe23f6dc8bdf3486f379c42c4b1be9c3a988360ea0592e70b9b0c35a2c71"
				data = fmt.Sprintf(`{"msgtype": "text","text": {"content": "%v"}}`, fmt.Sprintf("%v\nmd:%v \t%v", log.Content, log.Robot, time.Now().Format("01-02 15:04")))
			} else if strings.ToLower(log.Robot) == "oneplat" {
				//套币机器人
				secret = "SEC05dd86028265682c00b7d8e7abf1b0a73e71277414fe10c670cce83483e6c597"
				accesstoken = "45ed11014c6e2ca3bc171166f4307213b682b0960280433c5a9f7063863c7603"
				data = fmt.Sprintf(`{"msgtype": "text","text": {"content": "%v"}}`, fmt.Sprintf("%v\nmd:%v \t%v", log.Content, log.Robot, time.Now().Format("01-02 15:04")))
			} else if strings.ToLower(log.Robot) == "alert" {
				//观察机器人
				secret = "SECdbbec643fc237f3d6adad04abef87b947748d598be3f9385daf072c7b98f74bf"
				accesstoken = "bde3255fc76589180d1283f23fadb53e6276714ce450dabd1cc2ef4cc8bd11a9"
				data = fmt.Sprintf(`{"msgtype": "text","text": {"content": "%v"}}`, fmt.Sprintf("%v\nmd:%v \t%v", log.Content, log.Robot, time.Now().Format("01-02 15:04")))
			}

			var timestamp = time.Now().UnixNano() / 1e6

			stringToSign := fmt.Sprintf("%v\n%v", timestamp, secret)

			sign := signData(stringToSign, secret)

			var apiurl = fmt.Sprintf(
				"https://oapi.dingtalk.com/robot/send?access_token=%v&timestamp=%v&sign=%v",
				accesstoken,
				timestamp,
				sign)

			_, err := HttpPostRequest(apiurl, []byte(data))
			if err != nil {
				fmt.Println("dingrobit err", err.Error())
			} else {
				//fmt.Println("dingrobit ", err)
			}
		}
	}()
}

func SendDingTalkRobit(dingtalk bool, robot, typ, txid, content string) {
	log := Log{}
	log.Enable = dingtalk
	log.Robot = robot
	log.Typ = typ
	log.TxId = txid
	log.Content = content
	chLog <- log
	////filelog
	//PrintLog(typ, txid, content)
}

func signData(message string, secret string) string {
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(message))

	hmData := h.Sum(nil)

	bData := base64.StdEncoding.EncodeToString(hmData)

	urlData := url.QueryEscape(bData)

	return urlData
}
