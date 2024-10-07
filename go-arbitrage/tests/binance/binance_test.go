package binance

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"testing"
)

func getClient(robot string) *binance.Client {
	var apiKey, secretKey string
	if robot == "oneplat" {
		// oneplat
		apiKey = "3JiMItY7JeQoxWNAlylhsxCI38hysP5OZUgypWewm3PhKUaVx9pMv3dTUlyT5sbS"
		secretKey = "iPP0IRusNqUhKtVyl0gSteRnTEpMUXttUWSekH2MeqljcLkfzwyJ6J8nmUyUOxhn"
	} else if robot == "mainapi" {
		// mainapi
		apiKey = "mCXfycRaEiffizOajnB1VsVxytyUFnaA1tK4eX8QyuM8G565Weq5s4QXoyhkzwdE"
		secretKey = "wvRdYxo9O4IeBywbDCZgGhflwDwv2ERUbdQHUgoZ8JXTpUDGvFsTnXtzQOHxL9XW"
	}
	return binance.NewClient(apiKey, secretKey)
}

func Test_order(t *testing.T) {
	// mainapi
	client := getClient("oneplat")
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

func Test_sell(t *testing.T) {
	var sidetype = binance.SideTypeSell
	var symbol = "FILFDUSD"
	var quantity = "4957"
	var price = "4.159"

	client := getClient("oneplat")
	order, err := client.NewCreateOrderService().Type(binance.OrderTypeLimit).TimeInForce(binance.TimeInForceTypeGTC).
		Symbol(symbol).Side(sidetype).Quantity(quantity).Price(price).Do(context.Background())

	t.Log(order, err)
}
