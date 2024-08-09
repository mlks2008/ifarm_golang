package biz

var (
	DB_DELIMITER = ":"
)

// TYPE HASH
var (
	DB_SET_USER_EQUIPMENT          = "%s:equips:%v" //玩家装备列表 %1:equips:%2(%1prefix %2uid) [equipment_id...]
	DB_HASH_EQUIPMENT              = "%s:equip:%v"  //玩家装备详情 %1:equip:%2(%1prefix %2装备id)
	DB_FIELD_EQUIPMENT_ID          = "id"           //装备自增ID
	DB_FIELD_EQUIPMENT_BASE_ID     = "baseid"       //装备静态表ID
	DB_FIELD_EQUIPMENT_CREATE_TIME = "time"         //装备获取时间
)
