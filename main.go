package main

import (
	"fmt"
	User "go_im_room/client"
	"go_im_room/server"
	"os"
)

func main() {


	if len(os.Args) < 2{
		// 启动服务端：go run main.go s
		// 启动客户端：go run main.go c
		fmt.Println("需要参数启动参数，例如：go run main.go  s ")
		os.Exit(1)
	}
	parasm := os.Args[1:]
	for _, v := range parasm {
		if v == "s"{
			server.NewServer("127.0.0.1",8888).Start()
		}
		if v== "c"{
			User.NewUserClient("127.0.0.1",8888)
		}
	}
}
