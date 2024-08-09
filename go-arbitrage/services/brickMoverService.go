/**
 * brickMoverService.go
 * ============================================================================
 * 交易所间搬砖
 * ============================================================================
 * author: peter.wang
 */

package services

import (
	"goarbitrage/internal/plat"
)

type BrickMoverService struct {
	p1 plat.Plat
	p2 plat.Plat
}

func NewBrickMoverService(platcode1, platcode2 string) *BrickMoverService {
	ser := &BrickMoverService{
		p1: plat.Get(platcode1),
		p2: plat.Get(platcode2),
	}
	return ser
}

func (this *BrickMoverService) Start() {
	//如果plat1的买一 > plat2的卖一，计算价差高于两边手续费比例，则plat1卖，plat2买
	//如果plat2的买一 > plat1的卖一，计算价差高于两边手续费比例，则plat2卖，plat1买
}
