package redis

import (
	"context"
	
	"github.com/dennesshen/photon-core-starter/log"
	
	"github.com/redis/go-redis/v9"
)

type StartRedisAction func(ctx context.Context, redis *redis.ClusterClient) (err error)

var customAction []StartRedisAction

func RegisterRedisCustomize(action StartRedisAction) {
	customAction = append(customAction, action)
}

var rdb *redis.ClusterClient

func Redis() *redis.ClusterClient {
	return rdb
}

func Start(ctx context.Context) (err error) {
	rdb = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    config.Redis.Hosts,
		Password: config.Redis.Password,
	})
	
	// customize config
	for _, action := range customAction {
		if err = action(ctx, rdb); err != nil {
			log.Logger().Error(ctx, "failed to customize master database", "error", err)
			return
		}
	}
	
	// check connection
	if _, err = rdb.Ping(ctx).Result(); err != nil {
		log.Logger().Error(ctx, "connect to redis failure", "error", err)
		return
	}
	return
}
