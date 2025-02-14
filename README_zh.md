# GRM (go-redis-model) 

[![tests](https://github.com/go-redis-model/grm/actions/workflows/test.yml/badge.svg)](https://github.com/go-redis-model/grm/actions)[![codecov](https://codecov.io/gh/go-redis-model/grm/branch/main/graph/badge.svg?token=afb57ac7-039f-407c-ba10-921a0e63c385)](https://codecov.io/gh/go-redis-model/grm)[![Go Report Card](https://goreportcard.com/badge/github.com/go-redis-model/grm)](https://goreportcard.com/report/github.com/go-redis-model/grm)

[English](./README.md) | 中文

**基于 Go 的轻量级 Redis 模型层，灵感来自 GORM 风格**

**GRM** 是一个基于 `go-redis` 构建的极简 ORM 风格库，旨在通过 **结构体序列化** 简化 Redis 数据缓存，同时保持零魔法。特别适合需要轻量级模型持久化但不想引入复杂完整 ORM 的场景。

## ✨ 特性
- **自动 Key 管理**: 使用结构体名称（复数形式 + snake_case）和 `ID` 字段生成 Redis Key。  
  示例：`User` 结构体 → `grm:users:15`
- **时间戳自动化**: 自动维护 `CreatedAt` 和 `UpdatedAt` 字段。
- **自定义序列化**: 支持自定义序列化，并且内置JSON、MessagePack、Protobuf序列化。
- **CRUD 简化**: 直观的 `Set`、`Get` 和 `Delete` 操作方法。
- **解耦配置**: 支持精简配置或直接使用原生 `redis.Options` 满足高级需求。
- **严格测试**: 完整的单元测试和性能测试保障可靠性。

## 🚀 快速开始
```go
package main

import (
    "github.com/go-redis-model/grm"
)

type User struct {
    ID        uint
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func main() {
    // 连接
    db, _ := grm.Open(&grm.Options{Addr: "localhost:6379"})

    // 创建/更新
    user := User{ID: 1, Name: "Alice"}
    db.Set(&user) // Key: "grm:users:1"

    // 读取（需预填充 ID）
    fetched := User{ID: 1}
    db.Get(&fetched) 

    // 删除
    db.Delete(&fetched)

    // 批量保存用户
    users := []User{
        {ID: 1, Name: "Alice"},
        {ID: 2, Name: "Bob"},
    }
    db.Set(&users)

    // 批量读取（需预填充 ID）
    fetched := []User{{ID: 1}, {ID: 2}}
    db.Get(&fetched)

    // 批量删除
    db.Delete(&fetched)
}
```

## 🔄 序列化支持
GRM 支持自定义序列化，并且内置以下序列化器：

| 协议       | 常量                     | 说明                          |
|------------|--------------------------|-------------------------------|
| JSON       | `grm.JSONSerializer`     | 默认，兼容性好                |
| MessagePack| `grm.MessagePackSerializer` | 高性能二进制格式    |
| Protobuf   | `grm.ProtobufSerializer` | 需提前生成模型代码            |

### 示例：切换为 MessagePack
```go
db, _ := grm.Open(
    config,
    grm.WithSerializer(grm.MessagePackSerializer),
)
```
切换为 protobuf 需要提前生成模型代码，具体请参考 `examples/protobuf` 示例。

## 🔖 License

Licensed under [MIT License](./LICENSE)

## 💡 Contributors

See the [list of contributors](https://github.com/go-redis-model/grm/graphs/contributors).