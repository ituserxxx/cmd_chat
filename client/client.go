package User

import (
	"flag"
	"fmt"
	"io"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	Flag       int //当前用户模式
	Out bool
	Msg chan string
}

func NewUserClient(serIp string, serPo int,cname string) {
	flag.Parse()
	cl := &Client{
		Name: cname,
		ServerIp:   serIp,
		ServerPort: serPo,
		Conn:       nil,
		Flag:       999,
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
	cl.Msg<-"sys|"+cname
	cl.listenUserInput()
	_ = cl.Conn.Close()
}


func (c *Client) SelectOnlineUser()  {
	m := "who\n"
	_, err := c.Conn.Write([]byte(m))
	if err != nil {
		fmt.Println("SelectOnlineUser  send fail", err.Error())
		return
	}
}
func (c *Client) PrviteChat()  {
	var remoteName string
	var chatMsg string

	c.SelectOnlineUser()

	fmt.Println(">>>>>选择聊天对象：")
	_, _ = fmt.Scanln(&remoteName)
	for remoteName != "exit"{
		fmt.Println("<<<<<输入消息：")
		_, _ = fmt.Scanln(&chatMsg)
		for chatMsg != "exit"{
			if len(chatMsg) != 0{
				m := "to|"+remoteName+"|"+chatMsg
				_, err := c.Conn.Write([]byte(m))
				if err != nil {
					fmt.Println("PrviteChat send fail", err.Error())
					break
				}
			}
			chatMsg = ""
			fmt.Println("<<<<<<输入消息：")
			_, _ = fmt.Scanln(&chatMsg)
		}
		//从私聊退出后 再次选择用户
		c.SelectOnlineUser()
		remoteName = ""
		fmt.Println(">>>>选择聊天对象：")
		_, _ = fmt.Scanln(&remoteName)
	}

}

func (c *Client) UpdateName() bool {
	fmt.Println("请输入用户名：")
	_, _ = fmt.Scanln(&c.Name)
	m := "rename|" + c.Name
	c.Msg<-m
	return true
}
func (c *Client) listenUserInput()  {
	var chatMsg string
	_, _ = fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0{
			c.Msg<-chatMsg
		}
		chatMsg = ""
		_, _ = fmt.Scanln(&chatMsg)
	}
	return
}
func (c *Client) listenSendChan() {
	for {
		m := <-c.Msg
		_, err := c.Conn.Write([]byte("\n"+m))
		if err != nil {
			fmt.Println("client send Msg fail", err.Error())
		}
	}
}
//监听客户端输入，然后在发送给服务端
func (c *Client) handleUserAcceptMsg() {
	myMsgF := `(我)：%s`
	fmt.Println(fmt.Sprintf(myMsgF,c.Name))
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
		// 打印收到的消息
		msg := string(buf[:l])
		fmt.Println(msg)
		fmt.Println(fmt.Sprintf(myMsgF,c.Name))
	}

}