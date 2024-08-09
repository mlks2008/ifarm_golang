package resource

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

var EquipBreakDownConfigRes = new(equipBreakDownConfigList)

type (
	equipBreakDownConfig struct {
		Id              int    `json:"id"`
		Level           int    `json:"level"`
		Univer          string `json:"univer"`
		Quality         int    `json:"quality"`
		BreakDownReward []struct {
			Type    int `json:"type"`
			GoodsId int `json:"goodsId"`
			Num     int `json:"num"`
		} `json:"breakDownReward"`
		BreakDownCost []struct {
			Type    int `json:"type"`
			GoodsId int `json:"goodsId"`
			Num     int `json:"num"`
		} `json:"breakDownCost"`
	}

	equipBreakDownConfigList struct {
		dataMap map[string]*equipBreakDownConfig
	}
)

func (c *equipBreakDownConfigList) Name() string {
	return "Equip_breakDown"
}

func (c *equipBreakDownConfigList) OnReload(data []byte) error {
	dataList := make([]*equipBreakDownConfig, 0)
	err := jsoniter.Unmarshal(data, &dataList)
	if err != nil {
		panic(fmt.Sprintf("[Config.Parser] cfg[%s] unmarshal data error:%v", c.Name(), err))
	}

	dataMap := make(map[string]*equipBreakDownConfig)
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

func (c *equipBreakDownConfigList) GetEquipBreakDown(quility, level int32, univer string) *equipBreakDownConfig {
	key := fmt.Sprintf("%v#%v#%v", quility, level, univer)
	return c.dataMap[key]
}
