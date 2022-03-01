package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	ip         string
	port       int
	msg        chan string
	allUserMap map[string]*User
	mapLock    sync.RWMutex
}

func NewServer(ip string,port int)*Server  {
	return &Server{
		ip:         ip,
		port:       port,
		allUserMap: make(map[string]*User),
		msg:        make(chan string),
	}
}

func (s *Server)Start()  {
	listener,err := net.Listen("tcp",fmt.Sprintf("%s:%d",s.ip,s.port))
	if err != nil {
		fmt.Println("listen err",err.Error())
		return
	}
	defer listener.Close()

	//监听上线
	go s.hanlerMsg()

	//监听用户连接
	for  {
		con,err := listener.Accept()
		if err != nil {
			fmt.Println("listen accept err",err.Error())
			continue
		}
		go s.handlerUserAccept(con)
	}
}
func (s *Server)hanlerMsg()  {
	for  {
		m:= <-s.msg
		s.mapLock.Lock()
		for _, user := range s.allUserMap {
			user.Send(m)
		}
		s.mapLock.Unlock()
	}
}
func (s *Server)handlerUserAccept(conn net.Conn)  {
	u := &User{
		Name: "韩信偷塔"+conn.RemoteAddr().String(),
		Con:  conn,
	}
	s.mapLock.Lock()
	s.allUserMap[u.Name] = u
	s.mapLock.Unlock()
	s.msg <- "user :"+u.Name+"---上线了\n"
}