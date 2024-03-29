### 服务端部署

拉代码

```shell
git clone https://github.com/ituserxxx/cmd_chat.git
go mod tidy
```

### windows 环境运行


打包成 windows 环境执行文件

```shell
# 服务端
go build -o chat_server_win.exe code_server.go

# 客户端
go build -o chat_client_win.exe code_client.go
```

双击运行服务端

```shell
chat_server_win.exe
```
![img](./img/server_win_1.png)
![img](./img/server_win_2.png)

双击运行客户端

```shell
chat_client_win.exe
```
![img](./img/img_1.png)
![img](./img/img_2.png)

### 本地开发

运行服务端

```shell
go run code_server_linux.go ip port
```
![img](./img/img.png)

运行客户端

```shell
go run code_client_linux.go ip port name
```
![img](./img/win_chat.png)

### 打包服务端到 Linux 环境

打包

```shell
set GOARCH=amd64
set GOOS=linux
# 打包服务端
go build -o chat_server_linux code_server_linux.go
# 打包客户端
go build -o chat_client_linux code_client_linux.go
```

运行服务端

```shell
chat_server_linux ip port
```

运行客户端

```shell
chat_client_linux ip port name
```


