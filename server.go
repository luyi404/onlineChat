package main

import (
	"net"
	"fmt"
	"sync"
	"io"
)

type Server struct {
	Ip string
	Port int

	//在线user表
	OnlineMap map[string]*User
	mapLock sync.RWMutex

	//广播的channel
	Message chan string
}
//创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:	ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("链接建立成功")

	//用户上线了，而且应该直接广播
	//加入表中
	user := NewUser(conn, this)
	
	user.Online()

	//接受客户端发送的信息
	go func(){
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n==0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			//提取用户的消息
			msg := string(buf[:n-1])	//去掉了\n回车

			user.DoMessage(msg)
		}
	}()

}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

//监听Massage
func (this *Server) ListenMessager(){
	for {
		msg := <- this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap{
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}



//启动服务器的借口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//监听messager
	go this.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()	//表示上线了
		if err != nil{
			fmt.Println("listener accept err:", err)
			continue
		}
 
		//do handler
		go this.Handler(conn)
	}
}