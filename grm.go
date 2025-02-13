package grm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/kenshaw/snaker"
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

func (db *DB) Set(input interface{}) error {
	elements, err := processBatch(input)
	if err != nil {
		return err
	}

	// 收集键值对（格式: [key1, value1, key2, value2, ...]）
	ctx := context.Background()
	keyValues := make([]interface{}, 0, len(elements)*2)

	for _, elem := range elements {
		if elem.Kind() != reflect.Struct {
			return errors.New("element must be a struct")
		}

		model := elem.Addr().Interface()
		updateTimestamps(elem)

		key, err := getKey(model)
		if err != nil {
			return err
		}

		data, err := json.Marshal(model)
		if err != nil {
			return err
		}

		keyValues = append(keyValues, key, data)
	}

	return db.client.MSet(ctx, keyValues...).Err()
}

func (db *DB) Get(input interface{}) error {
	elements, err := processBatch(input)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(elements))
	for _, elem := range elements {
		model := elem.Addr().Interface()
		key, err := getKey(model)
		if err != nil {
			return err
		}
		keys = append(keys, key)
	}

	ctx := context.Background()
	values, err := db.client.MGet(ctx, keys...).Result()
	if err != nil {
		return err
	}

	errors := make(map[string]error)
	for i, val := range values {
		key := keys[i]
		if val == nil {
			errors[key] = fmt.Errorf("key not found")
			continue
		}

		data := []byte(val.(string))
		if err := json.Unmarshal(data, elements[i].Addr().Interface()); err != nil {
			errors[key] = fmt.Errorf("decode error: %v", err)
		}
	}

	if len(errors) > 0 {
		return &PartialError{Errors: errors}
	}
	return nil
}

func (db *DB) Delete(input interface{}) error {
	elements, err := processBatch(input)
	if err != nil {
		return err
	}

	ctx := context.Background()
	keys := make([]string, 0, len(elements))

	for _, elem := range elements {
		model := elem.Addr().Interface()
		key, err := getKey(model)
		if err != nil {
			return err
		}
		keys = append(keys, key)
	}

	return db.client.Del(ctx, keys...).Err()
}

// getKey 生成 Redis Key，格式为 "struct_prefix:id"
func getKey(model interface{}) (string, error) {
	// 获取结构体类型名称（如 "User"）
	t := reflect.TypeOf(model).Elem()
	structName := t.Name()

	// 转换为 snake_case 并复数化（如 "users"）
	prefix := snaker.CamelToSnake(structName) + "s"

	// 提取 ID 字段的值
	v := reflect.ValueOf(model).Elem()
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return "", fmt.Errorf("model must have an 'ID' field")
	}
	id := fmt.Sprintf("%v", idField.Interface())

	return fmt.Sprintf("grm:%s:%s", prefix, id), nil
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

// processBatch 解析输入，返回可遍历的反射值切片
func processBatch(input interface{}) ([]reflect.Value, error) {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Ptr {
		return nil, errors.New("input must be a pointer to a struct or slice/array")
	}

	elem := v.Elem()
	switch elem.Kind() {
	case reflect.Struct:
		// 单个结构体，包装成单元素切片
		return []reflect.Value{elem}, nil
	case reflect.Slice, reflect.Array:
		// 切片或数组，提取所有元素
		length := elem.Len()
		results := make([]reflect.Value, 0, length)
		for i := 0; i < length; i++ {
			results = append(results, elem.Index(i))
		}
		return results, nil
	default:
		return nil, errors.New("input must be a pointer to struct or slice/array")
	}
}
