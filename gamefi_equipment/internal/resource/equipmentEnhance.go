package resource

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

var EquipmentEnhanceConfigRes = new(equipmentEnhanceConfigList)

type (
	equipmentEnhanceConfig struct {
		Id            int64 `json:"Id"`
		Quility       int64 `json:"quility"`
		Level         int64 `json:"level"`
		Hp            int64 `json:"hp"`
		Atk           int64 `json:"atk"`
		Def           int64 `json:"def"`
		Speed         int64 `json:"speed"`
		PickRegion    int64 `json:"pickRegion"`
		Energy        int64 `json:"energy"`
		Wind          int64 `json:"wind"`
		Fire          int64 `json:"fire"`
		Water         int64 `json:"water"`
		Earth         int64 `json:"earth"`
		Thunder       int64 `json:"thunder"`
		Shadow        int64 `json:"shadow"`
		Holy          int64 `json:"holy"`
		Arcane        int64 `json:"arcane"`
		AntiWind      int64 `json:"antiWind"`
		AntiFire      int64 `json:"antiFire"`
		AntiWater     int64 `json:"antiWater"`
		AntiEarth     int64 `json:"antiEarth"`
		AntiThunder   int64 `json:"antiThunder"`
		AntiShadow    int64 `json:"antiShadow"`
		AntiHoly      int64 `json:"antiHoly"`
		AntiArcane    int64 `json:"antiArcane"`
		SucceedWeight int64 `json:"succeedWeight"`
		FailWeight    int64 `json:"failWeight"`
	}

	equipmentEnhanceConfigList struct {
		dataMap map[string]*equipmentEnhanceConfig
	}
)

func (c *equipmentEnhanceConfigList) Name() string {
	return "EquipmentEnhance"
}

func (c *equipmentEnhanceConfigList) OnReload(data []byte) error {
	dataList := make([]*equipmentEnhanceConfig, 0)
	err := jsoniter.Unmarshal(data, &dataList)
	if err != nil {
		panic(fmt.Sprintf("[Config.Parser] cfg[%s] unmarshal data error:%v", c.Name(), err))
	}

	dataMap := make(map[string]*equipmentEnhanceConfig)
	for _, enhance := range dataList {
		key := fmt.Sprintf("%v#%v", enhance.Quility, enhance.Level)
		if _, ok := dataMap[key]; ok {
			panic(fmt.Sprintf("%v.json配制表存在重复数据:%v-%v-%v", c.Name(), enhance.Id, enhance.Quility, enhance.Level))
		} else {
			dataMap[key] = enhance
		}
	}

	c.dataMap = dataMap
	return nil
}

func (c *equipmentEnhanceConfigList) GetEquipmentEnhance(quility, level int32) *equipmentEnhanceConfig {
	key := fmt.Sprintf("%v#%v", quility, level)
	return c.dataMap[key]
}
