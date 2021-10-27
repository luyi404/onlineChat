package main
import (
	"fmt"
	"net"
	"flag"
	"io"
	"os"
)
type Client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
	flag int
}

func NewClient(serverIp string, serverPort int) *Client {
	//先创建服务器对象
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		flag: 999,
	}
	//链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return nil
	}
	client.conn = conn
	return client
}

func (client *Client) menu() bool{
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	var flag int
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法的参数")
		return false
	}
}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>>请输入聊天内容, exit表示退出公聊模式")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {

		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err: ", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>>请输入聊天内容, exit表示退出公聊模式")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) PrivateChat() {
	_, err := client.conn.Write([]byte("who\n"))
	if err != nil {
		fmt.Println("conn Write err: ", err)
		return
	}
	var RemoteUserName string
	fmt.Println("请输入聊天对象的用户名")
	fmt.Scanln(&RemoteUserName)
	if RemoteUserName != "exit"{
		var chatMsg string
		fmt.Println(">>>>>请输入聊天内容, exit表示退出")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {

			if len(chatMsg) != 0 {
				sendMsg := "to|" + RemoteUserName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err: ", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>>请输入聊天内容, exit表示退出私聊模式")
			fmt.Scanln(&chatMsg)
		}
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名：")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil{
		fmt.Println("conn.Write error: ", err)
		return false
	}
	return true
}

func(client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}


func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			//公聊
			client.PublicChat()
			break
		case 2:
			//私聊
			client.PrivateChat()
			break
		case 3:
			//更新用户名
			client.UpdateName()
			break
		}
	}
}


var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址（默认为127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口（默认为8888）")
}



func main(){
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		 fmt.Println(">>>>>> Client 链接服务器失败...")
		 return
	}

	go client.DealResponse()


	fmt.Println(">>>>>> Client 链接服务器成功...")
	client.Run()
	// select{}
}