package component

import (
	"gamefi_equipment/internal/component/kafka"
	"gamefi_equipment/internal/component/redis"
	"gamefi_equipment/internal/component/redistool"
	"gamefi_equipment/internal/component/sdks"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(sdks.NewSdks, kafka.NewKafka, redis.NewRedis, redistool.NewRedisTool)
