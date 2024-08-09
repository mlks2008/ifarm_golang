package biz

import (
	credis "components/database/redis"
	clog "components/log/zaplogger"
	"context"
	"fmt"
	"gamefi_equipment/internal/component/redis"
	"gamefi_equipment/internal/conf"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func getRedis() *redis.Redis {
	c := config.New(
		config.WithSource(
			file.NewSource(fmt.Sprintf("../../configs/dev")),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	data, _, _ := redis.NewRedis(bc.Data, initLogger(bc.Log))
	return data
}
func initLogger(c *conf.Log) log.Logger {
	caller := func(depth int) log.Valuer {
		return func(context.Context) interface{} {
			_, file, line, _ := runtime.Caller(depth)
			idx := strings.LastIndexByte(file, '/')
			if idx == -1 {
				return file[idx+1:] + ":" + strconv.Itoa(line)
			}
			for i := 0; i < 2; i++ {
				idx = strings.LastIndexByte(file[:idx], '/')
			}
			return file[idx+1:] + ":" + strconv.Itoa(line)
		}
	}
	return log.With(clog.NewLoggerWithName(c.Dir, c.Name, "", c.Debug),
		"ts", log.DefaultTimestamp,
		"caller", caller(4),
		"service.id", "id",
		"service.name", c.Name,
		"service.version", "Version",
	)
}

// 缓存设计示例
func Test_redis(t *testing.T) {
	redis := getRedis()

	baseId := 11111
	userId := 10000

	//添加装备
	err := redis.RedisCli.SAdd(context.Background(),
		fmt.Sprintf(DB_SET_USER_EQUIPMENT, conf.Prefix, userId),
		1, 2,
	).Err()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("sadd ok")
	}
	//用户装备id列表
	idlist, err := redis.RedisCli.SMembers(context.Background(),
		fmt.Sprintf(DB_SET_USER_EQUIPMENT, conf.Prefix, userId),
	).Result()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("smembers ok", idlist)
	}

	//缓存用户装备详细
	for _, id := range idlist {
		err = redis.RedisCli.HMSet(context.Background(),
			fmt.Sprintf(DB_HASH_EQUIPMENT, conf.Prefix, id),
			DB_FIELD_EQUIPMENT_ID, id,
			DB_FIELD_EQUIPMENT_BASE_ID, baseId,
			DB_FIELD_EQUIPMENT_CREATE_TIME, time.Now().Unix(),
		).Err()
		if err != nil {
			t.Error(err)
		} else {
			t.Log("hmset ok", id)
		}
	}
	//用户装备详细列表
	//for _, id := range idlist {
	//	val, err := redis.RedisCli.HMGet(context.Background(),
	//		fmt.Sprintf(DB_HASH_EQUIPMENT, conf.Prefix, id),
	//		DB_FIELD_EQUIPMENT_ID,
	//		DB_FIELD_EQUIPMENT_BASE_ID,
	//		DB_FIELD_EQUIPMENT_CREATE_TIME,
	//	).Result()
	//	if err != nil {
	//		t.Error(err)
	//	} else {
	//		t.Log("hmget ok", fmt.Sprintf("%s", val[0]), fmt.Sprintf("%s", val[1]), fmt.Sprintf("%s", val[2]))
	//	}
	//}
	//用户装备详细列表（性能）
	ctx := context.Background()
	pipeline := redis.RedisCli.Pipeline()
	result := make([]*credis.SliceCmd, 0)
	for _, id := range idlist {
		key := fmt.Sprintf(DB_HASH_EQUIPMENT, conf.Prefix, id)
		result = append(result, pipeline.HMGet(ctx, key,
			DB_FIELD_EQUIPMENT_ID,
			DB_FIELD_EQUIPMENT_BASE_ID,
			DB_FIELD_EQUIPMENT_CREATE_TIME))
	}
	_, _ = pipeline.Exec(ctx)
	for _, r := range result {
		val, err := r.Result()
		if err != nil {
			t.Error(err)
		} else {
			t.Log("hmget ok", fmt.Sprintf("%s", val[0]), fmt.Sprintf("%s", val[1]), fmt.Sprintf("%s", val[2]))
		}
	}
}
