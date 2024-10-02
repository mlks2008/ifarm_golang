package main

import (
	"components/log/log"
	"time"

	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"goarbitrage/pkg/utils"

	"strconv"
	"strings"

	"github.com/adshao/go-binance/v2"
)

func placeOrder(side string, quantity, price string) (int64, error) {
	//qtyStr := strconv.FormatFloat(quantity, 'f', 0, 64)
	//priceStr := strconv.FormatFloat(price, 'f', 5, 64)

	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideType(side)).
		Type(binance.OrderTypeLimit).
		Quantity(quantity).
		Price(price).
		TimeInForce(binance.TimeInForceTypeGTC).
		Do(context.Background())
	if err != nil {
		//log.Logger.Error("Error placing %s order: %v", side, err)
		return 0, err
	}
	log.Logger.Debugf("Order placed: %s %s at %s(%s), orderId:%v", order.Symbol, order.Side, order.Price, quantity, order.OrderID)
	return order.OrderID, nil
}

func openOrders() ([]*binance.Order, error) {
	openOrders, err := client.NewListOpenOrdersService().Symbol(symbol).Do(context.Background())
	if err != nil {
		log.Logger.Error("Error fetching open orders:", err)
		return nil, err
	}
	return openOrders, nil
}

func cancelOrders(tye binance.SideType, symbol string, openOrders []*binance.Order) {
	for _, order := range openOrders {
		if order.Symbol != symbol {
			continue
		}
		if tye == "all" || order.Side == tye {
			_, err := client.NewCancelOrderService().
				Symbol(symbol).
				OrderID(order.OrderID).
				Do(context.Background())
			if err != nil {
				log.Logger.Error("Error cancelOrders orders:", err)
			} else {
				log.Logger.Debugf("Order %s %d canceled", order.Side, order.OrderID)
			}
		}
	}
}

func cancelOrder(orderId int64) error {
	_, err := client.NewCancelOrderService().
		Symbol(symbol).
		OrderID(orderId).
		Do(context.Background())
	if err != nil {
		log.Logger.Error("Error cancelOrder order:", orderId, err)
		return err
	} else {
		log.Logger.Debugf("Order %d canceled", orderId)
		return nil
	}
}

func getOrderStatus(symbol string, orderid int64) (bool, string, error) {
	order, err := client.NewGetOrderService().Symbol(symbol).OrderID(orderid).Do(context.Background())
	if err != nil {
		//log.Logger.Errorf("getOrderStatus %v %v", orderid, err)
		return false, "", err
	}
	if order.Status == binance.OrderStatusTypeFilled && order.OrigQuantity == order.ExecutedQuantity {
		return true, order.OrigQuantity, nil
	}
	return false, "", nil
}

func getCurrentPrice() (float64, error) {
	res, err := client.NewListPricesService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return 0, err
	}
	price, err := strconv.ParseFloat(res[0].Price, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func getBalances() (decimal.Decimal, decimal.Decimal, decimal.Decimal, error) {
	balances, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Logger.Error(err)
		return decimal.Zero, decimal.Zero, decimal.Zero, err
	}
	var filBalance, fdusdBalance, stopBalance decimal.Decimal
	for _, b := range balances.Balances {
		if b.Asset == "FIL" {
			b1, _ := decimal.NewFromString(b.Free)
			b2, _ := decimal.NewFromString(b.Locked)
			filBalance = b1.Add(b2)
		} else if b.Asset == "FDUSD" {
			b1, _ := decimal.NewFromString(b.Free)
			b2, _ := decimal.NewFromString(b.Locked)
			fdusdBalance = b1.Add(b2)
		} else if b.Asset == "CRV" {
			b1, _ := decimal.NewFromString(b.Free)
			b2, _ := decimal.NewFromString(b.Locked)
			stopBalance = b1.Add(b2)
		}
	}
	return filBalance, fdusdBalance, stopBalance, nil
}

// ------------------------------------------ files -------------------------------------------------

func RunGetFilCost(robot, symbol string) (decimal.Decimal, decimal.Decimal, int64) {
	cost := utils.ReadFile(fmt.Sprintf("files/down2_dcaservice.%v.%v.cost", robot, symbol))
	cost = strings.Replace(cost, "\n", "", -1)
	if cost == "" {
		currentFIL, currentUSDT, _, err := getBalances()
		if err != nil {
			panic(err)
		}
		RunSetFilCost(robot, symbol, currentUSDT.String(), currentFIL.String())
		return currentUSDT, currentFIL, time.Now().Unix()
	} else {
		v := strings.Split(cost, ",")
		costQuote, _ := decimal.NewFromString(v[0])
		costostBase, _ := decimal.NewFromString(v[1])
		updateTime, _ := strconv.ParseInt(v[2], 10, 64)
		return costQuote, costostBase, updateTime
	}
}

func RunSetFilCost(robot, symbol string, usdtbalance string, btcbalance string) {
	utils.UpdateFile(fmt.Sprintf("files/down2_dcaservice.%v.%v.cost", robot, symbol), usdtbalance+","+btcbalance+","+fmt.Sprintf("%v", time.Now().Unix()))
}

func RunGetInitPrice(robot, symbol string) (float64, error) {
	initprice := utils.ReadFile(fmt.Sprintf("files/down2_dcaservice.%v.%v.initprice", robot, symbol))
	initprice = strings.Replace(initprice, "\n", "", -1)
	if initprice == "" || initprice == "0" {
		return 0, nil
	} else {
		val, err := strconv.ParseFloat(initprice, 10)
		return val, err
	}
}

func RunSetInitPrice(robot, symbol string, price float64) {
	utils.UpdateFile(fmt.Sprintf("files/down2_dcaservice.%v.%v.initprice", robot, symbol), fmt.Sprintf("%v", price))
}

func RunGetInt64(robot, symbol string, filed string) int64 {
	val := utils.ReadFile(fmt.Sprintf("files/down2_dcaservice.%v.%v.%v", robot, symbol, filed))
	val = strings.Replace(val, "\n", "", -1)
	if val == "" || val == "0" {
		return 0
	} else {
		val, _ := strconv.ParseInt(val, 10, 64)
		return val
	}
}

func RunSetInt64(robot, symbol string, filed string, value int64) {
	utils.UpdateFile(fmt.Sprintf("files/down2_dcaservice.%v.%v.%v", robot, symbol, filed), fmt.Sprintf("%v", value))
}
