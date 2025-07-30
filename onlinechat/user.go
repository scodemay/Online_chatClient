package main

import "net"

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

//封装上下线，及消息处理
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

	if msg == "who"{
		for _,cil := range this.server.OnlineMap{
			OnlineMessage := "[" + cil.Addr + "]" + cil.Name + "在线...\n"
			this.SendMsg(OnlineMessage)
		}
	}else{

		this.server.BroadCast(this, msg)

	}
	
}

func (this *User) SendMsg (msg string) {
	this.conn.Write([]byte(msg))
}