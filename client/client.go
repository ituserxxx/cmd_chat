package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	Flag       int //当前用户模式
}

//命令解析：
var serIp string
var serPo int

func init() {
	flag.StringVar(&serIp, "ip", "127.0.0.1", "设置服务器ip(默认127.0.0.1)")
	flag.IntVar(&serPo, "port", 8888, "设置服务器port(默认8888)")
}

func main() {
	flag.Parse()
	cl := NewClient(serIp, serPo)
	if cl == nil {
		fmt.Println("client Dial fail:222222")
		return
	}
	fmt.Println("连接成功啦")

	//处理server回复的消息
	go cl.res()
	//启动菜单
	cl.Run()
}

func NewClient(serIp string, serPort int) *Client {
	cl := &Client{
		ServerIp:   serIp,
		ServerPort: serPort,
		Conn:       nil,
		Flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serIp, serPort))
	if err != nil {
		fmt.Println("client Dial fail:", err.Error())
		return nil
	}
	cl.Conn = conn
	return cl
}

func (c *Client) menu() bool {
	var flagMark int
	fmt.Println("1.公聊")
	fmt.Println("2.私聊")
	fmt.Println("3.改名聊")
	fmt.Println("0.退出")
	_, _ = fmt.Scanln(&flagMark)
	if flagMark >= 0 && flagMark < 4 {
		c.Flag = flagMark
		return true
	}
	return false
}

//循环菜单选择
func (c *Client) Run() {
	for c.Flag != 0 {
		for c.menu() == true {
			switch c.Flag {
			case 0:
				return
			case 1:
				fmt.Println("公聊模式开启")
				c.PublicChat()
			case 2:
				c.PrviteChat()
				fmt.Println("私聊模式开启")
			case 3:
				fmt.Println("改名开启")
				c.UpdateName()
			}
		}
	}
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
	_, err := c.Conn.Write([]byte(m))
	if err != nil {
		fmt.Println("send msg fail", err.Error())
		return false
	}
	return true
}
func (c *Client) PublicChat()  {
	var chatMsg string
	fmt.Println("你想说什么：")
	_, _ = fmt.Scanln(&chatMsg)
	fmt.Println("------------>"+chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0{
			_, err := c.Conn.Write([]byte(chatMsg+"\n"))
			if err != nil {
				fmt.Println("send msg fail", err.Error())
				break
			}
		}
		chatMsg = ""
		fmt.Println("你想说什么：")
		_, _ = fmt.Scanln(&chatMsg)
	}
	return
}

//监听客户端输入，然后在发送给服务端
func (c *Client) res() {
	//一旦client.conn 有数据就copy 到stdout 标准输出上，永久阻塞监听
	io.Copy(os.Stdout, c.Conn)
}
