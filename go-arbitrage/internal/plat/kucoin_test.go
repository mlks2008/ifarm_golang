package plat

import (
	"encoding/json"
	"fmt"
	"github.com/Kucoin/kucoin-go-sdk"
	"testing"
	"time"
)

func Test_Kucoin_Client(t *testing.T) {
	plat := NewKuCoin()

	rsp, err := plat.client.Symbols("")
	if err != nil {
		fmt.Println(err)
		return
	}
	l := kucoin.SymbolsModel{}
	if err := rsp.ReadData(&l); err != nil {
		fmt.Println(err)
		return
	}
	m, _ := json.Marshal(l)
	fmt.Println(fmt.Sprintf("%s", m))

	time.Sleep(time.Second)
}
