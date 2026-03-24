package redis

import "github.com/redis/go-redis/v9"

const (
	RedisNil = redis.Nil

	// Stream ID 常數（用於 XRead / XReadBlock / XRange）
	StreamIDStart  = "0" // 從 stream 開頭讀取
	StreamIDNewest = "$" // 只讀新訊息（用於追蹤最新）
	StreamIDMin    = "-" // XRange 最小 ID
	StreamIDMax    = "+" // XRange 最大 ID
)
