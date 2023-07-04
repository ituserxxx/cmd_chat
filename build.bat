
go build -o chat_server_win.exe code_server.go
go build -o chat_client_win.exe code_client.go

set GOARCH=amd64
set GOOS=linux
go build -o chat_server_linux code_server.go
go build -o chat_client_linux code_client.go
