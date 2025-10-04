package main

import (
	"context"
	"fmt"
	"os"
	v1 "user/api/user/v1"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var userClient v1.UserClient

var conn *grpc.ClientConn

func main() {
	Init()
	TestCreateUser()
	if conn != nil {
		_ = conn.Close()
	}
}

// Init 初始化 grpc 链接
func Init() {
	var err error
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	}

	target := os.Getenv("USER_GRPC_TARGET")
	if target == "" {
		if consulAddr := os.Getenv("CONSUL_ADDR"); consulAddr != "" {
			target = fmt.Sprintf("consul://%s/%s?wait=14s", consulAddr, "shop.users.service")
		} else {
			target = "127.0.0.1:50051"
		}
	}

	conn, err = grpc.Dial(target, opts...)
	if err != nil {
		panic("grpc link err" + err.Error())
	}
	userClient = v1.NewUserClient(conn)
}

func TestCreateUser() {
	rsp, err := userClient.CreateUser(context.Background(), &v1.CreateUserInfo{
		Mobile:   "13888888888",
		Password: "lucien",
		NickName: fmt.Sprintf("YWWW%d", 1),
	})
	if err != nil {
		panic("grpc 创建用户失败" + err.Error())
	}
	fmt.Println(rsp.Id)
}

func TestGetUser() {
	rsp, err := userClient.GetUserByMobile(context.Background(), &v1.MobileRequest{
		Mobile: "15283660176",
	})
	if err != nil {
		panic("grpc 获取用户失败" + err.Error())
	}
	fmt.Println(rsp)
}
