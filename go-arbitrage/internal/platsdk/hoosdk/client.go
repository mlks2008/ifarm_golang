package hoosdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"goarbitrage/pkg/utils"
	"io"
	"net/url"
	"time"
)

type (
	Client struct {
		Node      string
		ApiKey    string
		ApiSecret string
	}
)

// NewClient creates a new Client struct, with a single node URI.
func NewClient(nodeURI, apiKey, apiSecret string) *Client {
	c := &Client{
		Node:      nodeURI,
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
	}

	return c
}

// 查询系统支持的所有交易及精度
func (this Client) GetSymbols() ([]byte, error) {
	params := this.getSignAndParams()

	method := fmt.Sprintf("%s/open/v1/tickers/market", this.Node)

	return this.executeRequest(method, params)
}

// 获取交易深度信息
func (this Client) GetMarketDepth(symbol string) ([]byte, error) {
	params := this.getSignAndParams()
	params.Add("symbol", symbol)

	method := fmt.Sprintf("%s/open/v1/depth", this.Node)

	return this.executeRequest(method, params)
}

// 下单
func (this Client) OrdersPlace(symbol string, price string, quantity string, side string) ([]byte, error) {
	params := this.getSignAndParams()
	params.Add("symbol", symbol)   //交易对
	params.Add("price", symbol)    //价格
	params.Add("quantity", symbol) //数量
	if side == "SELL" {
		params.Add("side", "-1") //方向,1买，-1卖
	} else {
		params.Add("side", "1") //方向,1买，-1卖
	}
	method := fmt.Sprintf("%s/open/v1/orders/place", this.Node)

	return this.executePostRequest(method, params)
}
func (this Client) getSignAndParams() url.Values {
	params := url.Values{}
	params.Add("client_id", this.ApiKey)
	params.Add("nonce", utils.RandString(16))
	params.Add("ts", fmt.Sprintf("%v", time.Now().Unix()))

	message := fmt.Sprintf("client_id=%v&nonce=%v&ts=%v", params.Get("client_id"), params.Get("nonce"), params.Get("ts"))
	h := hmac.New(sha256.New, []byte(this.ApiSecret))
	io.WriteString(h, message)

	params.Add("sign", fmt.Sprintf("%x", h.Sum(nil)))
	return params
}
