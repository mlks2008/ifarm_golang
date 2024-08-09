package ent

import (
	"github.com/go-pg/pg/v10/orm"
	"time"
)

var (
	TableEquipment = (*EquipmentEnt)(nil)
)

type Attr struct {
	Attr  string `json:"attr"`
	Value int64  `json:"value"`
}

type EquipmentEnt struct {
	tableName    struct{} `pg:"platform.user_equipment_info"`
	Id           int64    `pg:"id,pk" json:"id"`
	BaseId       int64    `pg:"base_id" json:"base_id"`
	Type         int32    `pg:"type" json:"type"`
	Position     int32    `pg:"position" json:"position"`
	Quality      int32    `pg:"quality" json:"quality"`
	UserId       int64    `pg:"user_id" json:"user_id"`
	HeroId       int64    `pg:"hero_id" json:"hero_id"`
	UserHeroId   int64    `pg:"user_hero_id" json:"user_hero_id"`
	Level        int32    `pg:"level" json:"level"`
	Star         int32    `pg:"star" json:"star"`
	Status       int32    `pg:"status" json:"status"`
	DropAttrs    []*Attr  `pg:"drop_attrs" json:"drop_attrs"`
	UpgradeAttrs []*Attr  `pg:"upgrade_attrs" json:"upgrade_attrs"`
	//Hp          int64            `pg:"hp" json:"hp"`
	//Atk         int64            `pg:"atk" json:"atk"`
	//Def         int64            `pg:"def" json:"def"`
	//Speed       int64            `pg:"speed" json:"speed"`
	//PickRegion  int64            `pg:"pick_region" json:"pick_region"`
	//Wind        int64            `pg:"wind" json:"wind"`
	//Fire        int64            `pg:"fire" json:"fire"`
	//Water       int64            `pg:"water" json:"water"`
	//Earth       int64            `pg:"earth" json:"earth"`
	//Thunder     int64            `pg:"thunder" json:"thunder"`
	//Shadow      int64            `pg:"shadow" json:"shadow"`
	//Holy        int64            `pg:"holy" json:"holy"`
	//Arcane      int64            `pg:"arcane" json:"arcane"`
	//AntiWind    int64            `pg:"anti_wind" json:"antiWind"`
	//AntiFire    int64            `pg:"anti_fire" json:"antiFire"`
	//AntiWater   int64            `pg:"anti_water" json:"antiWater"`
	//AntiEarth   int64            `pg:"anti_earth" json:"antiEarth"`
	//AntiThunder int64            `pg:"anti_thunder" json:"antiThunder"`
	//AntiShadow  int64            `pg:"anti_shadow" json:"antiShadow"`
	//AntiHoly    int64            `pg:"anti_holy" json:"antiHoly"`
	//AntiArcane  int64            `pg:"anti_arcane" json:"antiArcane"`
	CreateTime time.Time `pg:"create_time,default:now()"`
	Updatetime time.Time `pg:"update_time"`
}

func (e *EquipmentEnt) Name() string {
	model, _ := orm.NewModel(e)
	return string(model.(orm.TableModel).Table().SQLName)
}
