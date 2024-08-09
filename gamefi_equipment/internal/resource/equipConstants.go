package resource

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

var EquipConstantsConfigRes = new(equipConstantsConfig)

type (
	equipConstantsConfig struct {
		WhiteMainNum  int `json:"whiteMainNum"`
		WhiteVicNum   int `json:"whiteVicNum"`
		BlueMainNum   int `json:"blueMainNum"`
		BlueVicNum    int `json:"blueVicNum"`
		PurpleMainNum int `json:"purpleMainNum"`
		PurpleVicNum  int `json:"purpleVicNum"`
		OrangeMainNum int `json:"orangeMainNum"`
		OrangeVicNum  int `json:"orangeVicNum"`
		RedMainNum    int `json:"redMainNum"`
		RedVicNum     int `json:"redVicNum"`
	}
)

func (c *equipConstantsConfig) Name() string {
	return "EquipConstants"
}

func (c *equipConstantsConfig) OnReload(data []byte) error {
	dataList := make([]*equipConstantsConfig, 0)
	err := jsoniter.Unmarshal(data, &dataList)
	if err != nil {
		panic(fmt.Sprintf("[Config.Parser] cfg[%s] unmarshal data error:%v", c.Name(), err))
	}

	for _, data := range dataList {
		EquipConstantsConfigRes = data
		break
	}

	return nil
}
