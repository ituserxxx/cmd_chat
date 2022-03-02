
go build -o server main.go server.go user.go

执行 server

linux 环境测试客户端连接
nc 127.0.0.1 8888

win 环境测试客户端链接
telnet 127.0.0.1 8888
