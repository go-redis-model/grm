package grm

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Options struct {
	// 精简字段（常用配置）
	Addr     string // 覆盖 redis.Options 的同名字段
	Password string // 覆盖 redis.Options 的同名字段
	DB       int    // 覆盖 redis.Options 的同名字段

	// 嵌入完整的 Redis 配置（支持高级配置）
	RedisOptions redis.Options
}

type DB struct {
	client *redis.Client
}

// Open 连接 Redis，返回 GRM 的 DB 实例
func Open(config *Options) (*DB, error) {
	// 如果精简字段有值，覆盖 RedisOptions 的对应字段
	if config.Addr != "" {
		config.RedisOptions.Addr = config.Addr
	}
	if config.Password != "" {
		config.RedisOptions.Password = config.Password
	}
	if config.DB != 0 {
		config.RedisOptions.DB = config.DB
	}

	client := redis.NewClient(&config.RedisOptions)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &DB{client: client}, nil
}

func (db *DB) Set(model interface{}) error {
	v := reflect.ValueOf(model).Elem()

	// 自动更新时间戳
	updateTimestamps(v)

	// 生成 Redis Key
	key, err := getKey(model)
	if err != nil {
		return err
	}

	// 序列化为 JSON
	data, err := json.Marshal(model)
	if err != nil {
		return err
	}

	// 存储到 Redis
	ctx := context.Background()
	return db.client.Set(ctx, key, data, 0).Err()
}

func (db *DB) Get(model interface{}) error {
	// 生成 Redis Key
	key, err := getKey(model)
	if err != nil {
		return err
	}

	// 从 Redis 读取数据
	ctx := context.Background()
	data, err := db.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	// 反序列化到结构体
	return json.Unmarshal(data, model)
}

func (db *DB) Delete(model interface{}) error {
	// 生成 Redis Key
	key, err := getKey(model)
	if err != nil {
		return err
	}

	// 删除 Key
	ctx := context.Background()
	return db.client.Del(ctx, key).Err()
}

// getKey 生成 Redis Key，格式为 "struct_prefix:id"
func getKey(model interface{}) (string, error) {
	// 获取结构体类型名称（如 "User"）
	t := reflect.TypeOf(model).Elem()
	structName := t.Name()

	// 转换为 snake_case 并复数化（如 "users"）
	prefix := pluralize(toSnakeCase(structName))

	// 提取 ID 字段的值
	v := reflect.ValueOf(model).Elem()
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return "", fmt.Errorf("model must have an 'ID' field")
	}
	id := fmt.Sprintf("%v", idField.Interface())

	return fmt.Sprintf("grm:%s:%s", prefix, id), nil
}

// 辅助函数：转换为 snake_case
func toSnakeCase(s string) string {
	// 简化实现，实际需要更复杂的处理（如 "MemberNumber" → "member_number"）
	return strings.ToLower(s)
}

// 辅助函数：复数化（简单加 "s"）
func pluralize(s string) string {
	return s + "s"
}

func updateTimestamps(v reflect.Value) {
	now := time.Now()

	// 设置 CreatedAt（仅当为零值时）
	createdAt := v.FieldByName("CreatedAt")
	if createdAt.IsValid() && createdAt.CanSet() {
		if createdAt.Interface().(time.Time).IsZero() {
			createdAt.Set(reflect.ValueOf(now))
		}
	}

	// 始终更新 UpdatedAt
	updatedAt := v.FieldByName("UpdatedAt")
	if updatedAt.IsValid() && updatedAt.CanSet() {
		updatedAt.Set(reflect.ValueOf(now))
	}
}
