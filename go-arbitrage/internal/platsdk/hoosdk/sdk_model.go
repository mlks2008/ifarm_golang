/**
 * model.go
 * ============================================================================
 * 简介
 * ============================================================================
 * author: peter.wang
 * createtime: 2020-07-07 23:26
 */

package hoosdk

//---------所有交易对深度-----
type Info struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}
type Depth struct {
	Bids []Info `json:"bids"` // 买盘, [price(成交价), amount(成交量)], 按price降序排列
	Asks []Info `json:"asks"` // 卖盘, [price(成交价), amount(成交量)], 按price升序排列
}

type DepthReturn struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Depth  `json:"data"`
}

type OrdersPlaceInfp struct {
	OrderId string `json:"order_id"`
	TradeNo string `json:"trade_no"`
}
type OrdersPlace struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data OrdersPlaceInfp `json:"data"`
}
