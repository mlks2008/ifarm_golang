package services

type OrderInfo struct {
	CreateTime  int64
	SellOrderId *string
	SellPrice   float64
	BuyOrderId  *string
	BuyPrice    float64
}
