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

// DBOption 是数据库级别的配置选项
type DBOption func(*DB)

func WithSerializer(s Serializer) DBOption {
	return func(db *DB) {
		db.serializer = s
	}
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
