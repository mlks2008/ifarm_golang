package main

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"time"
)

type CexBinanceCli struct {
	client *binance.Client
}

func NewCexBinanceCli(apikey, secretKey string) *CexBinanceCli {
	client := binance.NewClient(apikey, secretKey)

	return &CexBinanceCli{client: client}
}

func (this *CexBinanceCli) SpotBalances() (res *binance.Account, err error) {
	var fetch = func() (res *binance.Account, err error) {
		account, err := this.client.NewGetAccountService().Do(context.Background())
		return account, err
	}

	for i := 0; i < 3; i++ {
		res, err = fetch()
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		} else {
			return res, nil
		}
	}
	return nil, err
}

func (this *CexBinanceCli) CurMonthTransferHistory() (res []*binance.SubAccountTransferHistory, err error) {
	now := time.Now()

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Unix() * 1000
	currentTime := now.Unix() * 1000

	var fetch = func() (res []*binance.SubAccountTransferHistory, err error) {
		transfers, err := this.client.NewSubAccountTransferHistoryService().
			StartTime(startOfMonth).
			EndTime(currentTime).
			Do(context.Background())

		return transfers, err
	}

	for i := 0; i < 3; i++ {
		res, err = fetch()
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		} else {
			return res, nil
		}
	}
	return nil, err

}
