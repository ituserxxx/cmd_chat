
go build -o server main.go server.go user.go

执行 go run main.go

linux 环境测试连接
nc 127.0.0.1 8888

win10 环境测试连接
tenlet 127.0.0.1 8888

client连接

cd ./client 
执行 go run client.go



开启服务
go run main.go

开启客户端：-na指定名称
go run main.go -s=false -na=xx1

