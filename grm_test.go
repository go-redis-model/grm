package grm

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

// 定义测试模型
type TestUser struct {
	ID   uint32
	Name string
}

// 启动一个内存 Redis 服务器供测试使用
func setupTestRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return s
}

// 测试连接 Redis
func TestOpenConnection(t *testing.T) {
	s := setupTestRedis()
	defer s.Close()

	// 使用模拟 Redis 的地址
	config := &Options{
		Addr: s.Addr(),
	}
	db, err := Open(config)
	assert.NoError(t, err)
	assert.NotNil(t, db.client)
}

// 测试 Set 和 Get 操作
func TestSetAndGet(t *testing.T) {
	s := setupTestRedis()
	defer s.Close()

	config := &Options{Addr: s.Addr()}
	db, _ := Open(config)

	// 创建用户
	user := TestUser{
		ID:   1,
		Name: "Alice",
	}
	err := db.Set(&user)
	assert.NoError(t, err)

	// 读取用户
	fetched := TestUser{ID: 1}
	err = db.Get(&fetched)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", fetched.Name)

	// 读取不存在用户
	fetched.ID = 2
	err = db.Get(&fetched)
	assert.Error(t, err)
	assert.Equal(t, "partial error (1 failures)", err.Error())
}

func TestSetWithTTL(t *testing.T) {
	s := setupTestRedis()
	defer s.Close()

	config := &Options{Addr: s.Addr()}
	db, _ := Open(config)

	user := TestUser{ID: 1, Name: "Alice"}

	// 设置 TTL 为 1 秒
	err := db.Set(&user, WithTTL(time.Second))
	assert.NoError(t, err)

	// 验证 Key 存在
	err = db.Get(&user)
	assert.NoError(t, err)

	// 等待 TTL 过期
	s.FastForward(time.Second) // miniredis 提供的快进时间功能
	err = db.Get(&user)
	assert.Error(t, err) // Key 应已过期
}

// 测试 Delete 操作
func TestDelete(t *testing.T) {
	s := setupTestRedis()
	defer s.Close()

	config := &Options{Addr: s.Addr()}
	db, _ := Open(config)

	user := TestUser{ID: 1, Name: "Bob"}
	_ = db.Set(&user)

	// 删除用户
	err := db.Delete(&user)
	assert.NoError(t, err)

	// 验证用户不存在
	err = db.Get(&user)
	assert.Error(t, err) // 应返回 redis.Nil 错误
	key, _ := getKey(&user)
	assert.Equal(t, &PartialError{Errors: map[string]error{key: fmt.Errorf("key not found")}}, err)
}

func TestBatchOperations(t *testing.T) {
	s := setupTestRedis()
	defer s.Close()

	config := &Options{Addr: s.Addr()}
	db, _ := Open(config)

	// 定义批量用户
	users := []TestUser{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	// 批量保存
	err := db.Set(&users)
	assert.NoError(t, err)

	// 批量读取
	fetched := make([]TestUser, 2)
	fetched[0].ID = 1
	fetched[1].ID = 2
	err = db.Get(&fetched)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", fetched[0].Name)
	assert.Equal(t, "Bob", fetched[1].Name)

	// 批量删除
	err = db.Delete(&fetched)
	assert.NoError(t, err)
}

// 测试 Key 生成逻辑
func TestKeyGeneration(t *testing.T) {
	type Product struct {
		ID uint
	}
	p := Product{ID: 42}

	key, err := getKey(&p)
	assert.NoError(t, err)
	assert.Equal(t, "grm:products:42", key)
}

// 测试无效模型（缺少 ID 字段）
func TestInvalidModel(t *testing.T) {
	type Invalid struct {
		Name string
	}
	inv := Invalid{Name: "test"}

	_, err := getKey(&inv)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must have an 'ID' field")
}

// 基准测试 Set 操作
func BenchmarkSet(b *testing.B) {
	s := miniredis.RunT(b) // 集成测试中的临时 Redis
	config := &Options{Addr: s.Addr()}
	db, _ := Open(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := TestUser{ID: uint32(i), Name: "BenchUser"}
		_ = db.Set(&user)
	}
}

// 基准测试 Get 操作
func BenchmarkGet(b *testing.B) {
	s := miniredis.RunT(b)
	config := &Options{Addr: s.Addr()}
	db, _ := Open(config)

	// 预先插入数据
	user := TestUser{ID: 1, Name: "BenchUser"}
	_ = db.Set(&user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fetched := TestUser{ID: 1}
		_ = db.Get(&fetched)
	}
}

// 并发性能测试
func BenchmarkConcurrentSet(b *testing.B) {
	s := miniredis.RunT(b)
	config := &Options{Addr: s.Addr()}
	db, _ := Open(config)

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			counter++
			user := TestUser{ID: uint32(counter), Name: "ConcurrentUser"}
			_ = db.Set(&user)
		}
	})
}
