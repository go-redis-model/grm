# GRM (go-redis-model) 

[![tests](https://github.com/go-redis-model/grm/actions/workflows/test.yml/badge.svg)](https://github.com/go-redis-model/grm/actions)[![codecov](https://codecov.io/gh/go-redis-model/grm/branch/main/graph/badge.svg?token=afb57ac7-039f-407c-ba10-921a0e63c385)](https://codecov.io/gh/go-redis-model/grm)[![Go Report Card](https://goreportcard.com/badge/github.com/go-redis-model/grm)](https://goreportcard.com/report/github.com/go-redis-model/grm)

English | [中文](./README_zh.md)

**A Lightweight Redis Model Layer for Go, Inspired by GORM Style**

**GRM** is a minimalist ORM-style library built on top of `go-redis`, designed to simplify Redis data caching with **struct serialization** while keeping the magic to a minimum. Perfect for scenarios where you need lightweight model persistence without the complexity of a full ORM.

## ✨ Features
- **Auto Key Management**: Generates Redis keys using struct names (snake_case pluralized) and `ID` field.  
  Example: `User` struct → `grm:users:15`
- **Timestamp Automation**: Auto-populates `CreatedAt` and `UpdatedAt` fields.
- **Custom Serialization**: Supports custom serialization and has built-in serialization for JSON, MessagePack, and Protobuf.
- **CRUD Simplicity**: Intuitive `Set`, `Get`, and `Delete` operations for models.
- **Decoupled Config**: Use simplified options or extend with native `redis.Options` for advanced needs.
- **Battle-Tested**: Full unit and benchmark tests ensure reliability and performance.

## 🚀 Quick Start
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
    // Connect
    db, _ := grm.Open(&grm.Options{Addr: "localhost:6379"})

    // Create/Update
    user := User{ID: 1, Name: "Alice"}
    db.Set(&user) // Key: "grm:users:1"

    // Read（需预填充 ID）
    fetched := User{ID: 1}
    db.Get(&fetched) 

    // Delete
    db.Delete(&fetched)

    // 批量保存
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

## 🔄 Serialization Support
GRM supports custom serialization and includes the following built-in serializers:

| protocol       | Constant                     | Explanation                          |
|------------|--------------------------|-------------------------------|
| JSON       | `grm.JSONSerializer`     | By default, it is highly compatible.                |
| MessagePack| `grm.MessagePackSerializer` | High-performance binary format    |
| Protobuf   | `grm.ProtobufSerializer` | Please generate the model code in advance.            |

Switching to Protobuf requires generating model code in advance. Please refer to the example in `examples/protobuf` for details.

### Example: Switch to MessagePack
```go
db, _ := grm.Open(
    config,
    grm.WithSerializer(grm.MessagePackSerializer),
)
```

## 🔖 License

Licensed under [MIT License](./LICENSE)

## 💡 Contributors

See the [list of contributors](https://github.com/go-redis-model/grm/graphs/contributors).