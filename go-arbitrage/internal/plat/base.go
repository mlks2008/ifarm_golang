package plat

import (
	"github.com/pkg/errors"
	"goarbitrage/internal/model"
	"sync"
	"time"
)

type ReqStat struct {
	Time  int64
	Count int64
}

var (
	reqStat sync.Map
	reqTest sync.Map
)

type Base struct {
	waitLock  sync.Mutex //声明一个全局互斥锁
	waitLock2 sync.Mutex //声明一个全局互斥锁
}

// 请求超出每秒上限时等一秒
func (this *Base) wait(reqkey string, api string, limitSecond int64, symbol string) {
	this.waitLock.Lock()
	this.reqstat(reqkey, api, limitSecond, symbol)
	this.waitLock.Unlock()
}

// 请求超出每秒上限时等一秒
func (this *Base) wait2(reqkey string, api string, limitSecond int64, symbol string) {
	this.waitLock2.Lock()
	this.reqstat(reqkey, api, limitSecond, symbol)
	this.waitLock2.Unlock()
}

func (this *Base) reqstat(reqkey string, apiname string, limitSecond int64, symbol string) {
	if val, ok := reqStat.Load(reqkey); ok {
		reqstat := val.(*ReqStat)
		if reqstat.Time == time.Now().Unix() {
			reqstat.Count += 1
		} else {
			reqstat.Time = time.Now().Unix()
			reqstat.Count = 1
		}
		if reqstat.Count >= limitSecond {
			time.Sleep(time.Millisecond * 1000)
		}

		if true {
			val, _ := reqTest.Load(reqkey)
			reqm := val.(*ReqStat)
			if time.Unix(reqm.Time, 0).Format("2006-01-02 15:04") == time.Now().Format("2006-01-02 15:04") {
				reqm.Count += 1
			} else {
				//utils.PrintLog(api.PrintLog, "req_test"+symbol, fmt.Sprintf("%v", time.Now().Unix()/(60*30)), fmt.Sprintf("\t\t%v请求统计:%v Time:%v Count:%v", symbol, reqkey+"-"+apiname, time.Unix(reqm.Time, 0).Format("2006-01-02 15:04"), reqm.Count))
				reqm.Time = time.Now().Unix()
				reqm.Count = 1
			}
		}
	} else {
		reqStat.Store(reqkey, &ReqStat{Time: time.Now().Unix(), Count: 1})
		reqTest.Store(reqkey, &ReqStat{Time: time.Now().Unix(), Count: 1})
	}
}

func (this *Base) PlatCode() string {
	return ""
}

// 所有交易对
func (this *Base) GetSymbols() (*model.SymbolsReturn, error) {
	return nil, errors.New("继续类未实现")
}

// 支持的交易区
func (this *Base) GetScope() (map[string]string, error) {
	return nil, errors.New("继续类未实现")
}

// 深度
func (this *Base) GetMarketDepth(symbol string) (model.MarketDepthReturn, error) {
	return model.MarketDepthReturn{}, errors.New("继续类未实现")
}

// 下单
func (this *Base) OrdersPlace(symbol string, price string, quantity string, side string) (string, error) {
	return "", errors.New("继续类未实现")
}
