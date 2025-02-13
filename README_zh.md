# GRM (go-redis-model) 

[![Tests](https://github.com/go-redis-model/grm/actions/workflows/test.yml/badge.svg)](https://github.com/go-redis-model/grm/actions)[![codecov](https://codecov.io/gh/go-redis-model/grm/branch/main/graph/badge.svg?token=afb57ac7-039f-407c-ba10-921a0e63c385)](https://codecov.io/gh/go-redis-model/grm)

[English](./README.md) | 中文

**基于 Go 的轻量级 Redis 模型层，灵感来自 GORM 风格**

**GRM** 是一个基于 `go-redis` 构建的极简 ORM 风格库，旨在通过 **结构体序列化** 简化 Redis 数据缓存，同时保持零魔法。特别适合需要轻量级模型持久化但不想引入复杂完整 ORM 的场景。

## ✨ 特性
- **自动 Key 管理**: 使用结构体名称（复数形式 + snake_case）和 `ID` 字段生成 Redis Key。  
  示例：`User` 结构体 → `grm:users:15`
- **时间戳自动化**: 自动维护 `CreatedAt` 和 `UpdatedAt` 字段。
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

    // 读取
    fetched := User{ID: 1}
    db.Get(&fetched) 

    // 删除
    db.Delete(&fetched)
}
```

## 🔧 适用场景
- 游戏服务端中缓存玩家状态、排行榜或会话数据。
- 用结构体驱动操作替代重复的 Redis 序列化代码。
- 为项目添加轻量级缓存而不引入重型 ORM。


## 🔖 License

Licensed under [MIT License](./LICENSE)

## 💡 Contributors

See the [list of contributors](https://github.com/go-redis-model/grm/graphs/contributors).