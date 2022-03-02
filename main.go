package main

import "go_im_room/server"

func main() {
	server.NewServer("127.0.0.1",8888).Start()
}
