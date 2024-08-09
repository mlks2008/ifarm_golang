package conf

const (
	RUN_DEV  = "dev"
	RUN_QA   = "qa"
	RUN_YZ   = "yz"
	RUN_PROD = "prod"
)

const (
	ConditionType_universe int32 = 1 // 宇宙
	ConditionType_quality  int32 = 2 // 品质
	ConditionType_star     int32 = 3 // 星级
	ConditionType_game     int32 = 4 // 游戏
)

const (
	Quality_white  = 1 //白
	Quality_blue   = 2 //蓝
	Quality_purple = 3 //紫
	Quality_orange = 4 //橙
	Quality_red    = 5 //红
)

const (
	//主属性
	Attr_hp         = "hp"
	Attr_atk        = "atk"
	Attr_def        = "def"
	Attr_speed      = "speed"
	Attr_pickRegion = "pickRegion"
	Attr_energy     = "energy"
	//附加属性
	Attr_wind        = "wind"
	Attr_fire        = "fire"
	Attr_water       = "water"
	Attr_earth       = "earth"
	Attr_thunder     = "thunder"
	Attr_shadow      = "shadow"
	Attr_holy        = "holy"
	Attr_arcane      = "arcane"
	Attr_antiWind    = "antiWind"
	Attr_antiFire    = "antiFire"
	Attr_antiWater   = "antiWater"
	Attr_antiEarth   = "antiEarth"
	Attr_antiThunder = "antiThunder"
	Attr_antiShadow  = "antiShadow"
	Attr_antiHoly    = "antiHoly"
	Attr_antiArcane  = "antiArcane"
)
