package main

import (
	"bufio"
	User "cmd_chat/client"
	"cmd_chat/server"
	"flag"
	"os"
)

func main1() {
	//var chatMsg string
	//_, _ = fmt.Scan(&chatMsg)
	rd  := bufio.NewReader(os.Stdin)
	res,_,err := rd.ReadLine()
	if err != nil {
		println(err.Error())
		return
	}
	println("------",string(res))

}
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
