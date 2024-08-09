package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

var body []byte

// 使用示例
func main() {

	dev()
}

func dev() {
	//body = nil
	//request("GET", "http://127.0.0.1:8000/v1/assets/switch", body)

	body = nil
	request("GET", "http://127.0.0.1:8000/v1/users/0ce6c237-5538-46dc-9e32-0cda4e5bf16d/assets", body)

	//body = []byte(`{"asset_address":"","asset_id":"","chain":"","page":0,"page_size":10,"type":""}`)
	//request("POST", "http://127.0.0.1:8000/v1/users/0ce6c237-5538-46dc-9e32-0cda4e5bf16d/transactions", body)

	//body = []byte(`{"code": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1aWQiOiIxNDk4MDUyMTkxOTkzOTI1IiwiaXNHdWVzdCI6ZmFsc2UsImV4cCI6MTcwMzgyMTQ0MH0.n6vdGGVkbxfjYG58cPI-WDG69bEAy2_9HIjoNttebrQ"}`)
	//request("POST", "http://127.0.0.1:8000/v1/auth/verify-authorization-code", body)

	body = []byte(`{"user_id":"0ce6c237-5538-46dc-9e32-0cda4e5bf16d","chain":"","asset_id":"gp","asset_address":""}`)
	request("POST", "http://127.0.0.1:8000/v1/assets/withdraw-rules", body)
}

func request(method, url string, body []byte) {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	cli := http.Client{Timeout: 5 * time.Second}
	resp, err := cli.Do(req)
	if err != nil {
		fmt.Println("req err:", err)
		return
	}
	body, _ = io.ReadAll(resp.Body)

	fmt.Println("data", string(body))
}
