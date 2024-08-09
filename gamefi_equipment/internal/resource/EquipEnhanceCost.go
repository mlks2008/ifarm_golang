package resource

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

var EquipEnhanceCostConfigRes = new(equipEnhanceCostConfigList)

type (
	equipEnhanceCostConfig struct {
		Id          int    `json:"id"`
		Level       int    `json:"level"`
		Univer      string `json:"univer"`
		Quality     int    `json:"quality"`
		EnhanceCost []struct {
			Type    int `json:"type"`
			GoodsId int `json:"goodsId"`
			Num     int `json:"num"`
		} `json:"enhanceCost"`
	}

	equipEnhanceCostConfigList struct {
		dataMap map[string]*equipEnhanceCostConfig
	}
)

func (c *equipEnhanceCostConfigList) Name() string {
	return "Equip_enhance_cost"
}

func (c *equipEnhanceCostConfigList) OnReload(data []byte) error {
	dataList := make([]*equipEnhanceCostConfig, 0)
	err := jsoniter.Unmarshal(data, &dataList)
	if err != nil {
		panic(fmt.Sprintf("[Config.Parser] cfg[%s] unmarshal data error:%v", c.Name(), err))
	}

	dataMap := make(map[string]*equipEnhanceCostConfig)
	for _, item := range dataList {
		key := fmt.Sprintf("%v#%v#%v", item.Quality, item.Level, item.Univer)
		if _, ok := dataMap[key]; ok {
			panic(fmt.Sprintf("%v.json配制表存在重复数据:%v-%v-%v-%v", c.Name(), item.Id, item.Quality, item.Level, item.Univer))
		} else {
			dataMap[key] = item
		}
	}

	c.dataMap = dataMap
	return nil
}

func (c *equipEnhanceCostConfigList) GetEquipEnhanceCost(quility, level int32, univer string) *equipEnhanceCostConfig {
	key := fmt.Sprintf("%v#%v#%v", quility, level, univer)
	return c.dataMap[key]
}
