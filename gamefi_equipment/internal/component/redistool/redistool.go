package redistool

import (
	"context"
	"fmt"
	"gamefi_equipment/internal/component/redis"
	"gamefi_equipment/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type SceneId int32

const (
	Scene_ItemWithdraw SceneId = 1
)

type RedisTool struct {
	*redis.Redis
	log *log.Helper
}

func NewRedisTool(redis *redis.Redis, logger log.Logger) *RedisTool {
	return &RedisTool{
		Redis: redis,
		log:   log.NewHelper(logger),
	}
}

// 查找key是否存在，不存在时写入
func (r *RedisTool) Exist(user_id string, scene_id SceneId, ttl_second int64) bool {
	var key = fmt.Sprintf("%v:%v:%v:%v", conf.Prefix, "exist", user_id, scene_id)
	val := r.RedisCli.Exists(context.Background(), key).Val()

	if val == 0 {
		//不存在
		r.RedisCli.Set(context.Background(), key, 1, time.Duration(ttl_second)*time.Second)
		return false
	} else {
		//还在有效期
		return true
	}
}
