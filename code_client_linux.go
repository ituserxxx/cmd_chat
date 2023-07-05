package main

import (
	User "cmd_chat/client"
	"os"
)


func main() {
	parasm := os.Args[1:]

	// 客户端
	if len(parasm) == 3{
		User.NewUserClient(parasm[0],parasm[1],parasm[2])
		return
	}
	println("启动失败,请检查参数 \n 客户端启动示例 ：client ip port username")
}
