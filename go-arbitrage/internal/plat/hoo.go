package plat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"goarbitrage/internal/platsdk/hoosdk"
	"log"
)

type Hoo struct {
	Base
	client *hoosdk.Client
}

func NewHoo() *Hoo {
	obj := &Hoo{}
	obj.client = hoosdk.NewClient(viper.GetString("plats.hoo.apiaddr"), viper.GetString("plats.hoo.access_key"), viper.GetString("plats.hoo.secret_key"))
	return obj
}
func (this *Hoo) PlatCode() string {
	return P_Hoo
}

func (this *Hoo) OrdersPlace(symbol string, price string, quantity string, side string) (string, error) {
	body, err := this.client.OrdersPlace(symbol, price, quantity, side)
	if err != nil {
		log.Println(err)
	}
	res := &hoosdk.OrdersPlace{}
	json.Unmarshal([]byte(body), &res)
	if res.Code != 0 {
		return "", errors.New(res.Msg)
	}
	return fmt.Sprintf("%v-%v", res.Data.OrderId, res.Data.TradeNo), nil
}
