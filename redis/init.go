package redis

import (
	"context"
	"fmt"

	"github.com/Phofuture/photon-core-starter/log"

	"github.com/redis/go-redis/v9"
)

// RedisClient 統一接口，支持所有 Redis 客戶端類型
type RedisClient interface {
	redis.Cmdable
	Ping(ctx context.Context) *redis.StatusCmd
	Close() error
}

type StartRedisAction func(ctx context.Context, redis RedisClient) (err error)

var customAction []StartRedisAction

func RegisterRedisCustomize(action StartRedisAction) {
	customAction = append(customAction, action)
}

var rdb RedisClient

func Redis() RedisClient {
	return rdb
}

// NewRedisClient 工廠方法創建 Redis 客戶端
func NewRedisClient(clientType ClientType, addrs []string, password string) (RedisClient, error) {
	switch clientType {
	case ClientTypeStandalone:
		if len(addrs) == 0 {
			return nil, fmt.Errorf("standalone client requires at least one address")
		}
		return redis.NewClient(&redis.Options{
			Addr:     addrs[0],
			Password: password,
		}), nil
	case ClientTypeCluster:
		if len(addrs) == 0 {
			return nil, fmt.Errorf("cluster client requires at least one address")
		}
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			Password: password,
		}), nil
	default:
		return nil, fmt.Errorf("unsupported client type: %s", clientType)
	}
}

func Start(ctx context.Context) (err error) {
	// 決定客戶端類型，默認為 cluster 保持向後兼容
	clientType := ClientTypeCluster
	if config.Redis.Type != "" {
		clientType = ClientType(config.Redis.Type)
	}

	// 創建客戶端
	rdb, err = NewRedisClient(clientType, config.Redis.Hosts, config.Redis.Password)
	if err != nil {
		log.Logger().Error(ctx, "failed to create redis client", "error", err)
		return
	}

	// customize config
	for _, action := range customAction {
		if err = action(ctx, rdb); err != nil {
			log.Logger().Error(ctx, "failed to customize redis client", "error", err)
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
