package main

import (
	User "cmd_chat/client"
	"cmd_chat/server"
	"flag"
)

func main() {

	var isSer bool
	var cname string
	flag.BoolVar(&isSer, "s", true, "开启服务端,默认true")
	flag.StringVar(&cname, "na", "", "客户名称,不能为空")
	flag.Parse()

	if isSer{
		server.NewServer("127.0.0.1",8888)
	}else if cname!=""{
		User.NewUserClient("127.0.0.1",8888,cname)
	}

}
