package main
import (
	// "fmt"
	"net"
)
type User struct {
	Name string
	Addr string
	C	chan string
	conn net.Conn

	server *Server
}

//创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

//用户的上线业务
func (user *User) Online() {
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	user.server.BroadCast(user, "已上线")
}

//用户的下线业务
func (user *User) Offline(){
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	user.server.BroadCast(user, "下线")
}

func (user *User) SendMsg(msg string){
	user.conn.Write([]byte(msg))
}

func (user *User) DoMessage(msg string){
	if msg == "who" {
		//查询当前在线用户都有哪些
		user.server.mapLock.Lock()
		user.SendMsg("在线的用户有：\n")
		for _, eachuser := range user.server.OnlineMap {
			onlineMsg := "[" + eachuser.Addr + "]" + user.Name + ":" + " 在线...\n"
			user.SendMsg(onlineMsg)
		}
		user.server.mapLock.Unlock()
	} else {
		user.server.BroadCast(user, msg)
	}
}


//监听channel
func (this *User) ListenMessage() {
	for {
		msg := <- this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}