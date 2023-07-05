
go build -o chat_server_win.exe code_server_win.go
go build -o chat_client_win.exe code_client_win.go

set GOARCH=amd64
set GOOS=linux
go build -o chat_server_linux code_server_linux.go
go build -o chat_client_linux code_client_linux.go
