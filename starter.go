package redisStarter

import (
	"github.com/dennesshen/photon-core-starter/core"
	"github.com/dennesshen/photon-redis-starter/redis"
)

func init() {
	core.RegisterAddModule(redis.Start)
}
