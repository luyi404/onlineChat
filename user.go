package main

import (
	// "fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

//创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
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
func (user *User) Offline() {
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	user.server.BroadCast(user, "下线")
}

func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

func (user *User) DoMessage(msg string) {
	if msg == "who" {
		//查询当前在线用户都有哪些
		user.server.mapLock.Lock()
		user.SendMsg("在线的用户有：\n")
		for _, eachuser := range user.server.OnlineMap {
			onlineMsg := "[" + eachuser.Addr + "]" + eachuser.Name + ":" + " 在线...\n"
			user.SendMsg(onlineMsg)
		}
		user.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		//判断是不是重名了
		_, ok := user.server.OnlineMap[newName]
		if ok {
			user.SendMsg("此用户名已被使用\n")
		} else {
			user.server.mapLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[newName] = user
			user.server.mapLock.Unlock()

			user.Name = newName
			user.SendMsg("宁已经变更用户名为: " + newName + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//先获取用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			user.SendMsg("私聊消息格式不对\n")
			return
		}
		//再找User对象
		remoteUser, ok := user.server.OnlineMap[remoteName]
		if !ok {
			user.SendMsg("用户不存在，请使用 \"who\" 命令来查询在线用户\n")
			return
		}
		//获取消息内容并发送
		content := strings.Split(msg, "|")[2]
		if content == "" {
			user.SendMsg("无消息内容，请重发\n")
			return
		}
		remoteUser.SendMsg(user.Name + "私聊你说： " + content + "\n")
	} else {
		user.server.BroadCast(user, msg)
	}
}

//监听channel
func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n"))
	}
}
