package main

import (
	"fmt"
	grm "go-redis-model"
	"time"

	"github.com/redis/go-redis/v9"
)

type User struct {
	ID       uint
	Name     string
	Email    *string
	Age      uint8
	Birthday *time.Time
	// CreatedAt time.Time
	// UpdatedAt time.Time
	ignored string // 非导出字段会被忽略
}

func main() {
	// 连接 Redis
	config := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	db, err := grm.Open(config)
	if err != nil {
		panic(err)
	}

	// 保存记录
	email := "alice@example.com"
	user := User{
		ID:    15,
		Name:  "Alice",
		Email: &email,
		Age:   30,
	}

	if err := db.Set(&user); err != nil {
		panic(err)
	}

	// 读取记录
	fetchedUser := User{ID: 15}
	if err := db.Get(&fetchedUser); err != nil {
		panic(err)
	}
	fmt.Println(fetchedUser.Name) // 输出 "Alice"

	// //删除记录
	// if err := db.Delete(&fetchedUser); err != nil {
	// 	panic(err)
	// }
}
