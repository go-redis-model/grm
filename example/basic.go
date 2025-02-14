package main

import (
	"fmt"
	"time"

	"github.com/go-redis-model/grm"
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
	config := &grm.Options{
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

	//删除记录
	if err := db.Delete(&fetchedUser); err != nil {
		panic(err)
	}

	// 批量保存用户
	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	db.Set(&users)
	fmt.Println(users)

	// 批量读取（需预填充 ID）
	fetched := []User{{ID: 1}, {ID: 2}}
	db.Get(&fetched)

	// 批量删除
	db.Delete(&fetched)
}
