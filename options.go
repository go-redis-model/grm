package grm

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// redis 配置
type Options struct {
	// 精简字段（常用配置）
	Addr     string // 覆盖 redis.Options 的同名字段
	Password string // 覆盖 redis.Options 的同名字段
	DB       int    // 覆盖 redis.Options 的同名字段

	// 嵌入完整的 Redis 配置（支持高级配置）
	RedisOptions redis.Options
}

type SetOption func(*setConfig)

type setConfig struct {
	ttl time.Duration
}

func WithTTL(d time.Duration) SetOption {
	return func(cfg *setConfig) {
		cfg.ttl = d
	}
}
