package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	ip        string
	port      int
	msg       chan string
	onlineMap map[string]*User
	mapLock   sync.RWMutex
}

//初始化服务连接
func NewServer(ip string, port int) *Server {
	return &Server{
		ip:        ip,
		port:      port,
		onlineMap: make(map[string]*User), //在线用户
		msg:       make(chan string),      //广播消息
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
	if err != nil {
		fmt.Println("listen err", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("房间已开启~~")

	//监听广播消息
	go s.GuangboMsg()

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

func (s *Server) handlerUserAccept(conn net.Conn) {
	//初始化 用户
	u := NewUser("xxx1:="+conn.RemoteAddr().String(), conn, s)

	//用户上线
	u.Online()

	// 在线状态
	isLive := make(chan bool)
	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			l, err := conn.Read(buf)
			//合法关闭
			if l == 0 {
				//广播用户下线
				u.Downline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("read Err :", err.Error())
				return
			}
			//获取出消息
			msg := string(buf[:l-1])
			//发送消息
			u.DoMessage(msg)
			isLive<-true
		}
	}()

	//超时强踢------还没测试成功-----
	for {
		select {
		case <-isLive:
		case <-time.After(time.Second * 5):
			u.C <- "你已经被剔下线\n"
			close(u.C)
			_ = conn.Close()
			//s.mapLock.Lock()
			//delete(s.onlineMap, u.Name)
			//s.mapLock.Unlock()
		}
	}
	//当前handler阻塞
	//select {}
}

//处理用户输入消息
func (s *Server) handlerUserInputMsg(conn net.Conn, u *User, isLive chan bool) {

}
func (s *Server) GuangboMsg() {
	for {
		//取出广播消息
		m := <-s.msg
		s.mapLock.Lock()

		//遍历发送给每个用户的消息 channel
		for _, user := range s.onlineMap {
			user.C <- m
		}
		s.mapLock.Unlock()
	}
}
