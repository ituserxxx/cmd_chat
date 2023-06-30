package User

import (
	"cmd_chat/comm"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	GoOut      chan string
	Msg        chan comm.MsgInfo
}

var msgF = `
用户：%s		%s
%s
`
var infoF = `~~notice 用户%s	%s	%s`

func NewUserClient(serIp string, serPo int, cname string) {
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
		GoOut:      make(chan string),
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serIp, serPo))
	if err != nil {
		fmt.Println("client Dial fail:", err.Error())
		return
	}
	cl.Conn = conn
	// 监听用户输入
	go cl.handleUserAcceptMsg()
	go cl.listenSendChan()
	cl.Msg <- comm.MsgInfo{
		Event:comm.EventInitName,
		Data:  cname,
	}

	select {
	case <-cl.GoOut:
		return
	case <-time.After(time.Second * 5):
		break
	}

	var chatMsg string
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			cl.Msg <- comm.MsgInfo{
				Event: comm.EventPublicMsg,
				Data:  cl.Name+":" +chatMsg,
			}

		}
		chatMsg = ""
		_, _ = fmt.Scanln(&chatMsg)

	}
	_ = cl.Conn.Close()
	close(cl.Msg)
}



func (c *Client) listenSendChan() {
	for {
		m := <-c.Msg
		v, _ := json.Marshal(m)
		_, err := c.Conn.Write(v)
		if err != nil {
			fmt.Println("client send Data fail", err.Error())
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
		err = json.Unmarshal([]byte(msg), &d)
		if err != nil {
			fmt.Println("unmarshal failed!")
			continue
		}
		c.DoMessage(d)
	}

}
func (c *Client) DoMessage(msg *comm.MsgInfo) {
	if msg.Event == comm.EventInitName {
		if msg.Code != 0{
			fmt.Println(msg.Data)
			c.GoOut<-"1"
			return
		}
		c.Name = msg.Data
		fmt.Println(fmt.Sprintf( `(我)：%s`, c.Name))
		return
	}
	if msg.Event == comm.EventAllUsers {
		fmt.Println(msg.Data)
		//fmt.Println(fmt.Sprintf(myMsgF, c.Name))
		return
	}
	if msg.Event == comm.EventAT {
		if  msg.Code != 0{
			fmt.Println(msg.Data)
			return
		}
	}
	fmt.Println("\n"+msg.Data)
	if msg.Event != comm.EventSysInfo{
		fmt.Print(`(我)：`)
	}
	//fmt.Println(fmt.Sprintf(myMsgF, c.Name))

}
