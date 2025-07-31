package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

//build Api

func NewUser(conn net.Conn, server *Server) *User {

	UserAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   UserAddr,
		Addr:   UserAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//listen user channel
	go user.ListenMessage()

	return user
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}

// 封装上下线，及消息处理
func (this *User) Online() {

	this.server.maplock.Lock()

	this.server.OnlineMap[this.Name] = this

	this.server.maplock.Unlock()

	this.server.BroadCast(this, "已上线")
}

func (this *User) Offline() {
	this.server.maplock.Lock()

	delete(this.server.OnlineMap, this.Name)

	this.server.maplock.Unlock()

	this.server.BroadCast(this, "已下线")
}

func (this *User) DoMessage(msg string) {

	if msg == "who" {
		for _, cil := range this.server.OnlineMap {
			OnlineMessage := "[" + cil.Addr + "]" + cil.Name + " 在线...\n"
			this.SendMsg(OnlineMessage)
		}
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		NewName := strings.Split(msg, "|")[1]

		_, ok := this.server.OnlineMap[NewName]
		if ok {
			this.SendMsg("名字已被占用")
		} else {
			this.server.maplock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[NewName] = this
			this.server.maplock.Unlock()
			this.Name = NewName
			this.SendMsg("已更新名字：" + this.Name + "\n")
		}

	} else if len(msg) > 3 && msg[:3] == "to|" {

		RemoteName := strings.Split(msg, "|")[1]
		if RemoteName == "" {
			this.SendMsg("格式不正确\n")
			return
		}

		RemoteUser, ok := this.server.OnlineMap[RemoteName]

		if !ok {
			this.SendMsg("该用户名不存在\n")
			return
		}

		content := strings.Split(msg, "|")[2]

		if content == "" {
			this.SendMsg("请输入有效内容\n")
			return
		}

		RemoteUser.SendMsg(this.Name + " 对您说 " + content + " \n")

	} else {
		this.server.BroadCast(this, msg)
	}

}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}
