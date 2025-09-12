package redisStarter

import (
	"github.com/Phofuture/photon-redis-starter/redis"
	"github.com/dennesshen/photon-core-starter/core"
)

func init() {
	core.RegisterAddModule(redis.Start)
}
