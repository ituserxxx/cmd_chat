package server

import (
	"cmd_chat/comm"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type ChatServer struct {
	ip              string
	port            string
	GbMsg           chan string      //广播消息
	onlineUserTotal int64            //在线用户总数
	onlineMap       map[string]*User //在线用户
	mapLock         sync.RWMutex
	PrintChan       chan string //大厅日志
}

var IMserver *ChatServer

func NewServer(ip ,port string) {

	ser := &ChatServer{
		ip:        ip,
		port:      port,
		onlineMap: make(map[string]*User),
		GbMsg:     make(chan string),
		PrintChan: make(chan string),
	}
	IMserver = ser
	ser.Start()
}

func (s *ChatServer) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.ip, s.port))
	if err != nil {
		println(err.Error())
		fmt.Println("服务开启失败，请检查参数~~")
		return
	}
	defer listener.Close()
	//监听广播消息
	go s.GuangboMsg()

	// 大厅日志输出
	go s.printLog()

	s.PrintChan <- "~~服务大厅已开启~~"

	for {
		//监听用户连接
		con, err := listener.Accept()
		if err != nil {
			fmt.Println("listen accept err", err.Error())
			continue
		}
		//处理用户逻辑
		go s.handlerUserAccept(con)
	}
}

func (s *ChatServer) handlerUserAccept(conn net.Conn) {
	//初始化 用户
	u := CreateNewUser("", conn)
	//广播用户下线
	defer func() {
		u.Downline()
	}()
	// 开启消息接收监听
	go u.HandleMsg()
	u.Online()
	//接收客户端发送的消息
	buf := make([]byte, 1024)
	for {
		l, err := conn.Read(buf)
		//合法关闭
		if l == 0 {
			break
		}
		if err != nil && err != io.EOF {
			fmt.Println("read Err :", err.Error())
			continue
		}
		// 获取收到的消息
		msg := strings.TrimSpace(string(buf[:l]))

		var d *comm.MsgInfo
		err = json.Unmarshal(comm.B64Encry(msg), &d)
		if err != nil {
			fmt.Println("unmarshal failed!")
			continue
		}
		u.DoMessage(d)

	}

}

func (s *ChatServer) GuangboMsg() {
	for {
		//取出广播消息
		m := <-s.GbMsg
		//遍历发送给每个用户的消息 channel
		for _, user := range s.onlineMap {
			user.C <- comm.MsgInfo{
				Event: comm.EventGuangbo,
				Data:  m,
			}
		}
	}
}
func (s *ChatServer) printLog() {
	for {
		v, ok := <-s.PrintChan
		if ok {
			fmt.Printf("\n[log->%s] %s", time.Now().Format("2006-01-02 15:04:05"), v)
		}
	}
}

func (s *ChatServer) GuangboMsgToAllUser(event, msg string) {

	for _, user := range s.onlineMap {

		user.C <- comm.MsgInfo{
			Event: event,
			Data:  msg,
		}
	}
}
func (s *ChatServer) GuangboMsgToOtherUser(event, uId, msg string) {
	for _, user := range s.onlineMap {
		if user.ID != uId {
			user.C <- comm.MsgInfo{
				Event: event,
				Data:  msg,
			}
		}
	}
}
