package server

import (
	"cmd_chat/comm"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

type User struct {
	ID   string
	Name string
	C    chan comm.MsgInfo
	Conn net.Conn
	Status string
}

func CreateNewUser(name string, con net.Conn) *User {
	//初始化用户
	u := &User{
		ID:   comm.Krand(),
		Name: name,
		C:    make(chan comm.MsgInfo),
		Conn: con,
	}

	return u

}

// HandleMsg 监听当前User channel ，一旦有消息，直接发送给当前用户
func (u *User) HandleMsg() {
	for {
		if u.Status == "out"{
			u.Conn.Close()
			break
		}
		select {
		case m :=<-u.C:
			v, _ := json.Marshal(m)
			_, err := u.Conn.Write(v)
			if err != nil {
				fmt.Println("send GbMsg fail", err.Error())
			}
		default:
			time.Sleep(time.Second)
		}
	}

}

// DoMessage 处理交互的消息
func (u *User) DoMessage(msg *comm.MsgInfo) {

	if msg.Event == comm.EventInitName {
		newName := msg.Data
		info := comm.MsgInfo{
			Event: comm.EventInitName,
		}
		for _, us := range IMserver.onlineMap {
			if us.Name == newName {
				info.Code = 1
				info.Data = fmt.Sprintf("系统提示:昵称 '%s' 已经被人用了，换个吧", newName)
				u.C<-info
				u.Status ="out"
				return
			}
		}
		u.Name = newName
		u.Status ="online"
		info.Data = newName

		IMserver.GuangboMsgToOtherUser(comm.EventSysInfo,   u.ID,"系统提示:"+u.Name+"上线了~~")
		u.C <- info
		return
	}

	if msg.Event == comm.EventInputAllUsers {
		au := fmt.Sprintf("\n------------------\n--在线总人数：%d 人", IMserver.onlineUserTotal)
		for _, us := range IMserver.onlineMap {
			au += "\n--" + us.Name
		}
		au += "\n" + "------------------"
		u.C <- comm.MsgInfo{
			Event: comm.EventInputAllUsers,
			Data:  au,
		}
		return
	}

	if msg.Event == comm.EventInputAT {
		//msg:{Event:"@", Data:"xx1:ok", Code:0}
		mhIndex := strings.Index(msg.Data,":")
		toUserName := msg.Data[:mhIndex]
		toMsg := msg.Data[mhIndex+1:]
		var toUser *User
		for _, toUser = range IMserver.onlineMap {
			if toUser.Name == toUserName {
				break
			}
		}
		if toUser == nil {
			u.C <- comm.MsgInfo{
				Code: 1,
				Data: "对方用户已经不在线了~~",
			}
			return
		}

		toUser.C <- comm.MsgInfo{
			Event: comm.EventInputAT,
			Data:  "悄悄话：用户：" + u.Name + " 对我说：" + toMsg,
		}
		return
	}
	IMserver.GuangboMsgToOtherUser(comm.EventPublicMsg, u.ID,msg.Data)

}
func (u *User) Downline() {
	IMserver.mapLock.Lock()
	delete(IMserver.onlineMap, u.ID)
	IMserver.onlineUserTotal--
	IMserver.mapLock.Unlock()

	if u.Status == "online"{
		IMserver.GuangboMsgToOtherUser(comm.EventSysInfo, u.ID,  "系统提示:"+u.Name+"下线了~~")
	}
	IMserver.PrintChan <- fmt.Sprintf("\nuser down:%s   当前总人数 total:%d", u.ID, IMserver.onlineUserTotal)
	_ = u.Conn.Close()
}

func (u *User) Online() {
	//加入在线用户列表
	IMserver.mapLock.Lock()
	IMserver.onlineMap[u.ID] = u
	IMserver.onlineUserTotal++
	IMserver.mapLock.Unlock()

	IMserver.PrintChan <- fmt.Sprintf("\n user online:%s   当前总人数 total:%d", u.ID, IMserver.onlineUserTotal)
}
