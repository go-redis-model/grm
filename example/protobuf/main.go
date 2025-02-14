package main

import (
	"fmt"

	"github.com/go-redis-model/grm"
	pb "github.com/go-redis-model/grm/example/protobuf/pb"
)

func main() {
	// 连接 Redis
	config := &grm.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	db, err := grm.Open(config, grm.WithSerializer(grm.ProtobufSerializer))
	if err != nil {
		panic(err)
	}

	// 保存记录
	email := "alice@example.com"
	user := pb.User{
		ID:    15,
		Name:  "Alice",
		Email: email,
		Age:   30,
	}

	if err := db.Set(&user); err != nil {
		panic(err)
	}

	// 读取记录
	fetchedUser := pb.User{ID: 15}
	if err := db.Get(&fetchedUser); err != nil {
		panic(err)
	}
	fmt.Println(fetchedUser.Name) // 输出 "Alice"

	//删除记录
	if err := db.Delete(&fetchedUser); err != nil {
		panic(err)
	}

	// 批量保存用户
	users := []pb.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	db.Set(&users)
	fmt.Println(users)

	// 批量读取（需预填充 ID）
	fetched := []pb.User{{ID: 1}, {ID: 2}}
	db.Get(&fetched)

	// 批量删除
	db.Delete(&fetched)
}
