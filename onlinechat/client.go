package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct{
	Name string
	ServerIp string
	ServerPort int
	conn  net.Conn
	flag   int 
}

var ServerIp string
var ServerPort int

func init() {
	flag.StringVar(&ServerIp,"ip","127.0.0.1","查看服务器默认地址")
	flag.IntVar(&ServerPort,"port",8888,"查看服务器默认端口")
}


func NewClient(ServerIp string, ServerPort int) *Client {

	client := &Client{
		ServerIp: ServerIp,
		ServerPort: ServerPort,
		flag :      999,
	}

	conn,err := net.Dial("tcp", fmt.Sprintf("%s:%d",ServerIp,ServerPort))

	if err != nil {
		fmt.Println("net.Dial err:",err)
		return  nil
	}
	client.conn = conn

	return client

}

func main() {

	flag.Parse()

    client := NewClient(ServerIp , ServerPort)

    if client == nil {
        fmt.Println("服务器链接失败......")
        return
    }

    fmt.Println("服务器链接成功......")

    
    go client.DealRespond()
    client.Run()

	
}

func (client *Client) UpdateName() bool{

	fmt.Println("请输入用户名")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	
	_,err := client.conn.Write([]byte(sendMsg))

	if err != nil {
		fmt.Println("conn Write error:",err)
		return false
	}
	return  true
}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println("已可以进行公聊,输入exit退出\n")
	fmt.Scanln(&chatMsg)

	for chatMsg !="exit" {
		sendMsg := chatMsg +"\n"
		_,err := client.conn.Write([]byte(sendMsg))
		 if err !=  nil {
			fmt.Println(" conn Write err:",err)
			return
		 }

		chatMsg = ""
		fmt.Println("已可以进行公聊,输入exit退出\n")
		fmt.Scanln(&chatMsg)

	}
	
}

func (client *Client) selectUser() {
	sendMsg := "who\n"
	_,err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err:",err)
	}
} 

func (client *Client) PraviteChat () {
	var chatMsg string
	var RemoteName string

	client.selectUser()
	fmt.Println("选择一个用户进行聊天,输入exit退出")
	fmt.Scanln(&RemoteName)

	for RemoteName != "exit" {
		fmt.Println("已可以发送消息,输入exit退出\n")
	    fmt.Scanln(&chatMsg)

		for chatMsg !="exit" {
			sendMsg := "to|" + RemoteName +"|" + chatMsg +"\n"
			_,err := client.conn.Write([]byte(sendMsg))
		 	if err !=  nil {
				fmt.Println(" conn Write err:",err)
				return
		 	}
			chatMsg = ""
			fmt.Println("已可以发送消息,输入exit退出\n")
			fmt.Scanln(&chatMsg)
	
		}
		client.selectUser()
		fmt.Println("选择一个用户进行聊天,输入exit退出")
	    fmt.Scanln(&RemoteName)
		
	} 

}

func (client *Client) DealRespond () {
	io.Copy( os.Stdout,client.conn )
}

func (client *Client) menu () bool {
	var flag int 

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更改用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >=0 && flag <=3 {
		client.flag =flag
		return true
	} else {
		fmt.Println("请输入合法数字")
		return false
	}

}

func (client *Client) Run() {
	for client.flag !=0 { 
		for client.menu() != true {

		}
		switch client.flag {
		case 1:
			fmt.Println("公聊模式")
			client.PublicChat()
			break

		case 2:
			fmt.Println("私聊模式")
			client.PraviteChat()
			break

		case 3:
			client.UpdateName()
			break
			
		}
	}
}

