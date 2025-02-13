# GRM (go-redis-model) 

[![tests](https://github.com/go-redis-model/grm/actions/workflows/test.yml/badge.svg)](https://github.com/go-redis-model/grm/actions)[![codecov](https://codecov.io/gh/go-redis-model/grm/branch/main/graph/badge.svg?token=afb57ac7-039f-407c-ba10-921a0e63c385)](https://codecov.io/gh/go-redis-model/grm)[![badge](https://goreportcard.com/badge/github.com/go-redis-model/grm)](https://goreportcard.com/report/github.com/go-redis-model/grm)

English | [中文](./README_zh.md)

**A Lightweight Redis Model Layer for Go, Inspired by GORM Style**

**GRM** is a minimalist ORM-style library built on top of `go-redis`, designed to simplify Redis data caching with **struct serialization** while keeping the magic to a minimum. Perfect for scenarios where you need lightweight model persistence without the complexity of a full ORM.

## ✨ Features
- **Auto Key Management**: Generates Redis keys using struct names (snake_case pluralized) and `ID` field.  
  Example: `User` struct → `grm:users:15`
- **Timestamp Automation**: Auto-populates `CreatedAt` and `UpdatedAt` fields.
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

    // Read
    fetched := User{ID: 1}
    db.Get(&fetched) 

    // Delete
    db.Delete(&fetched)
}
```

## 🔧 When to Use?
- Cache player states, leaderboards, or session data in game backends.
- Replace boilerplate Redis serialization code with struct-driven operations.
- Add lightweight caching without adopting a heavy ORM.

## 🔖 License

Licensed under [MIT License](./LICENSE)

## 💡 Contributors

See the [list of contributors](https://github.com/go-redis-model/grm/graphs/contributors).