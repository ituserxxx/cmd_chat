package server

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	C    chan string
	Con  net.Conn
	Ser  *Server
}

func NewUser(name string, con net.Conn, server *Server) *User {
	//初始化用户
	u := &User{
		Name: name,
		C:    make(chan string),
		Con:  con,
		Ser:  server,
	}

	// 开启消息接收监听
	go u.HandleMsg()

	return u

}

// 监听当前User channel ，一旦有消息，直接发送给当前用户
func (u *User) HandleMsg() {
	for {
		m := <-u.C
		_, err := u.Con.Write([]byte(m))
		if err != nil {
			fmt.Println("send msg fail", err.Error())
		}
	}

}
func (u *User) DoMessage(msg string) {
	//输入who 则查询在线的所有人
	if msg == "who" {
		u.Ser.mapLock.Lock()
		for _, us := range u.Ser.onlineMap {
			u.C <- "用户：" + us.Name + "在线\n"
		}
		u.Ser.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		name := strings.Split(msg, "|")[1]

		u.Ser.mapLock.Lock()
		_, ok := u.Ser.onlineMap[name]
		if ok {
			u.C <- "用户名已存在\n"
			u.Ser.mapLock.Unlock()
			return
		}
		delete(u.Ser.onlineMap, u.Name)
		u.Ser.onlineMap[name] = u
		u.Name = name
		u.C <- "更名成功->" + name + "\n"
		u.Ser.mapLock.Unlock()
	} else {
		u.Ser.msg <- "[user :说----->" + u.Name + msg + "\n"
	}

}
func (u *User) Downline() {
	u.Ser.msg <- "[user :说----->" + u.Name + "]---xia线了\n"
	u.Ser.mapLock.Lock()
	delete(u.Ser.onlineMap, u.Name)
	u.Ser.mapLock.Unlock()
}

func (u *User) Online() {
	//加入在线用户列表
	u.Ser.mapLock.Lock()
	u.Ser.onlineMap[u.Name] = u
	u.Ser.mapLock.Unlock()
	// 广播消息给所有用户
	u.Ser.msg <- "[user :说----->" + u.Name + "]---上线了\n"
}
