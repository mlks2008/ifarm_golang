package resource

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

var EquipmentConfigRes = new(equipmentConfigList)

type (
	AttValue struct {
		Type   string `json:"type"`
		Value  []int  `json:"value"`
		Weight int    `json:"weight"`
	}

	equipmentConfig struct {
		Id           int64 `json:"id"`
		EffectiveAtt []struct {
			Type  int         `json:"type"`
			Value interface{} `json:"value"`
		} `json:"effectiveAtt"` // 生效条件
		Type         int32  `json:"type"`
		PartType     int32  `json:"partType"` // 部位
		Quality      int32  `json:"quality"`
		Universe     string `json:"universe"`
		Legend       int64  `json:"legend"` // 专属
		UseCondition []struct {
			Type  int         `json:"type"`
			Value interface{} `json:"value"`
		} `json:"useCondition"` // 穿戴条件
		RunesNum int        `json:"runesNum"`
		GameAtt  []AttValue `json:"gameAtt"`  // 主属性
		SecndAtt []AttValue `json:"secndAtt"` // 附加属性
	}

	equipmentConfigList struct {
		dataList []*equipmentConfig
		dataMap  map[int64]*equipmentConfig
	}
)

func (c *equipmentConfigList) Name() string {
	return "Equipment"
}

func (c *equipmentConfigList) OnReload(data []byte) error {
	dataList := make([]*equipmentConfig, 0)
	err := jsoniter.Unmarshal(data, &dataList)
	if err != nil {
		panic(fmt.Sprintf("[Config.Parser] cfg[%s] unmarshal data error:%v", c.Name(), err))
	}

	dataMap := make(map[int64]*equipmentConfig)
	for _, equipment := range dataList {
		dataMap[equipment.Id] = equipment
	}

	c.dataList = dataList
	c.dataMap = dataMap
	return nil
}

func (c *equipmentConfigList) GetById(id int64) *equipmentConfig {
	return c.dataMap[id]
}
