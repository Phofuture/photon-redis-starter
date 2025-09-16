package redisStarter

import (
	"github.com/Phofuture/photon-core-starter/core"
	"github.com/Phofuture/photon-redis-starter/redis"
)

func init() {
	core.RegisterAddModule(redis.Start)
}
