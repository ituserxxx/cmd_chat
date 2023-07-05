package main

import (
	"cmd_chat/server"
	"fmt"
	"os"
)

func main() {
	parasm := os.Args[1:]
	fmt.Printf("[%#v]", parasm)
	//服务端
	if len(parasm) == 2 {
		server.NewServer(parasm[0], parasm[1])
		return
	}
	println("启动失败,请检查参数 \n 服务端启动示例 ：server ip port")

}
