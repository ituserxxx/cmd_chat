package server

import (
	"cmd_chat/utils"
	"fmt"
	"net"
	"strings"
)

type User struct {
	ID   string
	Name string
	C    chan string
	Conn net.Conn
}

func CreateNewUser(name string, con net.Conn) *User {
	//初始化用户
	u := &User{
		ID:   utils.Krand(),
		Name: name,
		C:    make(chan string),
		Conn: con,
	}

	// 开启消息接收监听
	go u.HandleMsg()
	u.Online()
	return u

}

// HandleMsg 监听当前User channel ，一旦有消息，直接发送给当前用户
func (u *User) HandleMsg() {
	for {
		m := <-u.C
		_, err := u.Conn.Write([]byte("\n"+m))
		if err != nil {
			fmt.Println("send GbMsg fail", err.Error())
		}
	}

}

// DoMessage 处理交互的消息
func (u *User) DoMessage(msg string) {
	//输入who 则查询在线的所有人
	if msg == "who" {//当前在线用户
		IMserver.mapLock.Lock()

		for _, us := range IMserver.onlineMap {
			u.C <- "|" + us.Name + "|在线"
		}
		IMserver.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {//改名功能
		name := strings.Split(msg, "|")[1]
		IMserver.mapLock.Lock()
		_, ok := IMserver.onlineMap[name]
		if ok {
			u.C <- "用户名已存在\n"
			IMserver.mapLock.Unlock()
			return
		}
		delete(IMserver.onlineMap, u.Name)
		IMserver.onlineMap[name] = u
		IMserver.mapLock.Unlock()
		u.Name = name
		u.C <- "更名成功->" + name + "\n"
	}else if len(msg) > 4 && msg[:2] == "to"{ //私聊功能
		toU := strings.Split(msg, "|")[1]
		toUser,ok := IMserver.onlineMap[toU]
		if !ok{
			u.C <- "当前用户不存在\n"
			return
		}
		toMsg := strings.Split(msg, "|")[2]
		if toMsg == ""{
			u.C <- "消息格式不对，示例：to|张三|你好\n"
			return
		}
		toUser.C<-"用户："+u.Name+"对你说："+toMsg+"\n"

	}else if msg[:3]=="sys"{

	}else {
		IMserver.GuangboMsgToOtherUser( u.ID,msg)
	}

}
func (u *User) Downline() {
	IMserver.mapLock.Lock()
	delete(IMserver.onlineMap, u.ID)
	IMserver.mapLock.Unlock()
	IMserver.onlineUserTotal--
	IMserver.PrintChan<-fmt.Sprintf("user down:%s   user total:%d", u.ID, IMserver.onlineUserTotal)
	IMserver.UserOnlineAndDownline( u.ID,"下线")
	_= u.Conn.Close()

}

func (u *User) Online() {
	//加入在线用户列表
	IMserver.mapLock.Lock()
	IMserver.onlineMap[u.ID] = u
	IMserver.mapLock.Unlock()
	// 在线总数+1
	IMserver.onlineUserTotal++
	IMserver.UserOnlineAndDownline( u.ID,"上线")
	IMserver.PrintChan<-fmt.Sprintf("user online:%s   user total:%d", u.ID, IMserver.onlineUserTotal)
}
