package User

import (
	"bufio"
	"cmd_chat/comm"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

type Client struct {
	ServerIp   string
	ServerPort string
	Name       string
	Conn       net.Conn
	Status       string
	Msg        chan comm.MsgInfo
}

func NewUserClient(serIp ,serPo , cname string) {
	if len(cname) == 0 {
		fmt.Println("昵称一定要有哦")
	}
	flag.Parse()
	cl := &Client{
		Name:       cname,
		ServerIp:   serIp,
		ServerPort: serPo,
		Conn:       nil,
		Msg:        make(chan comm.MsgInfo),
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", serIp, serPo))
	if err != nil {
		fmt.Println("client Dial fail:", err.Error())
		return
	}
	defer func() {
		conn.Close()
		close(cl.Msg)
	}()
	cl.Conn = conn
	// 监听用户输入
	go cl.handleUserAcceptMsg()
	go cl.listenSendChan()
	cl.Msg <- comm.MsgInfo{
		Event: comm.EventInitName,
		Data:  strings.TrimSpace(cname),
	}
	var i int
	for {
		if i > 5{
			break
		}
		switch cl.Status {
		case "fial":
			return
		case "online":
			break
		default:

		}
		time.Sleep(time.Second)
		i++
	}


	var chatMsg string
	rd := bufio.NewReader(os.Stdin)
	rdMsg, _, err := rd.ReadLine()
	if err != nil {
		println(err.Error())
		return
	}
	chatMsg = strings.TrimSpace(string(rdMsg))
	for chatMsg != "exit" {
		var prMyIsTruc bool
		if len(chatMsg) != 0 {
			inputMsg := comm.MsgInfo{
				Event: comm.EventPublicMsg,
			}
			if chatMsg == comm.EventInputAllUsers{
				inputMsg.Event = comm.EventInputAllUsers
			}else if chatMsg[:1] == comm.EventInputAT{
				inputMsg.Event = comm.EventInputAT
				prMyIsTruc = true
				chatMsg = chatMsg[1:]
				mhIndex := strings.Index(chatMsg,":")
				toUserName :=strings.TrimSpace( chatMsg[:mhIndex])
				toMsg := strings.TrimSpace(chatMsg[mhIndex+1:])

				if len(toUserName) == 0||len(toMsg)==0{
					fmt.Println("系统提示:消息格式不对，示例：@xxx:你在干嘛？")
					fmt.Print("(我):")
					return
				}
				inputMsg.Data = fmt.Sprintf("%s:%s",toUserName,toMsg)
			}else{
				prMyIsTruc = true
				inputMsg.Data = fmt.Sprintf("%s 说：%s", cl.Name, chatMsg)
			}

			cl.Msg <- inputMsg
		}
		if prMyIsTruc {
			fmt.Print("(我):")
		}
		chatMsg = ""
		rd = bufio.NewReader(os.Stdin)
		rdMsg, _, err = rd.ReadLine()
		if err != nil {
			println(err.Error())
			continue
		}
		chatMsg = string(rdMsg)
	}
}

func (c *Client) listenSendChan() {
	for {
		if c.Status == "fial"{
			return
		}
		select {
		case m := <-c.Msg:
			v, _ := json.Marshal(m)
			_, err := c.Conn.Write([]byte(comm.B64Encode(v)))
			if err != nil {
				fmt.Println("client send Data fail", err.Error())
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

//监听客户端输入，然后在发送给服务端
func (c *Client) handleUserAcceptMsg() {
	//一旦client.conn 有数据就copy 到stdout 标准输出上，永久阻塞监听
	//io.Copy(os.Stdout, c.Conn)
	buf := make([]byte, 1024)
	for {
		l, err := c.Conn.Read(buf)
		//合法关闭
		if l == 0 {
			return
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
		c.DoMessage(d)
	}

}
func (c *Client) DoMessage(msg *comm.MsgInfo) {
	//fmt.Printf("msg %#v", msg)
	if msg.Event == comm.EventInitName {
		if msg.Code != 0 {
			fmt.Println(msg.Data)
			c.Status = "fial"
			return
		}
		c.Status = "online"
		c.Name = msg.Data
		fmt.Println("我" + msg.Data + "又回来啦 >_< ")
		fmt.Print(`(我)：`)
		return
	}
	if msg.Event == comm.EventInputAllUsers {
		fmt.Println(msg.Data)
		fmt.Print(`(我)：`)
		return
	}
	if msg.Event == comm.EventInputAT {
		if msg.Code != 0 {
			fmt.Println(msg.Data)
			return
		}
	}
	fmt.Println("\n" + msg.Data)
	if msg.Code != 0 {
		fmt.Println(msg.Data)
		return
	}
	fmt.Print(`(我)：`)
}
