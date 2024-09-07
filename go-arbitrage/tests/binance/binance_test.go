package binance

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"testing"
)

func Test_order(t *testing.T) {
	// mainapi
	apiKey := "mCXfycRaEiffizOajnB1VsVxytyUFnaA1tK4eX8QyuM8G565Weq5s4QXoyhkzwdE"
	secretKey := "wvRdYxo9O4IeBywbDCZgGhflwDwv2ERUbdQHUgoZ8JXTpUDGvFsTnXtzQOHxL9XW"
	client := binance.NewClient(apiKey, secretKey)
	order, err := client.NewGetOrderService().Symbol("DOGEFDUSD").OrderID(767669619).Do(context.Background())
	fmt.Println(order, err)

	////sell
	//order.ExecutedQuantity //doge
	//order.CummulativeQuoteQuantity //fdusd
	//buy

	list, err := client.NewListOrdersService().Symbol("DOGEFDUSD").Limit(50).Do(context.Background())
	for _, order := range list {
		if order.Side == binance.SideTypeBuy && order.Status == binance.OrderStatusTypeFilled {
			fmt.Println(order)
		}
		if order.Side == binance.SideTypeSell && order.Status == binance.OrderStatusTypeFilled {
			fmt.Println(order)
		}
	}

}
