package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User

	maplock sync.RWMutex

	Message chan string
}

// listen Message
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.maplock.Lock()
		for _, cil := range this.OnlineMap {
			cil.C <- msg
		}
		this.maplock.Unlock()
	}
}

// func broadcast
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + msg

	this.Message <- sendMsg
}

//连接当前业务

func (this *Server) Handler(conn net.Conn) {
	//fmt.Println("链接成功")

	user := NewUser(conn, this)

	user.Online()

	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn err:", err)
				return
			}

			msg := string(buf[:n-1])

			user.DoMessage(msg)
		}
	}()

	select {}

}

//创建一个server接口

func NewServer(ip string, port int) *Server {
	Server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return Server
}

// 定义socket连接
func (this *Server) Start() {
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	defer Listener.Close()

	go this.ListenMessage()

	//accept
	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		//do Handler
		go this.Handler(conn)
	}
}
