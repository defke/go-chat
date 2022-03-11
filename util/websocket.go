package util

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

//消息管理
var Manager = ClientManager{
	Broadcast:  make(chan Message),         //发送的消息通道
	Register:   make(chan *Client),         //注册的用户连接
	Unregister: make(chan *Client),         //断开连接的用户信息
	Clients:    make(map[string][]*Client), //根据用户id连接的用户
}

type Message struct {
	//发送人
	Sender string `json:"sender,omitempty"`
	//发送的房间号
	RoomId string `json:"roomId,omitempty"`
	//发送的内容信息
	Content string `json:"content"`
	Ts      string `json:"ts"`
}

type Client struct {
	Id     string //当前用户id
	RoomId string //用户连接的6位序列号
	//Avatar string    //用户的头像
	Socket *websocket.Conn
	//给用户发送的信息
	Send chan Message
}

type ClientManager struct {
	//根据用户id查找连接的用户
	Clients map[string][]*Client
	//发送的信息都到这里，在选择转发发送人
	Broadcast chan Message
	//注册的连接
	Register chan *Client
	//断开的连接
	Unregister chan *Client
}

//程序开启时启动监听
func (manager *ClientManager) Start() {
	go func() {
		for {
			log.Println("<---管道通信--->")
			select {
			case client := <-manager.Register:
				log.Printf("新用户加入:%v", client.RoomId)
				clientList := manager.Clients[client.RoomId]
				clientList = append(clientList, client)
				manager.Clients[client.RoomId] = clientList
			case client := <-manager.Unregister:
				//关闭发送channel，避免资源浪费
				close(client.Send)
				clientList := manager.Clients[client.RoomId]
				for i, existClient := range clientList {
					if existClient.Id == client.Id {
						//移除
						clientList = append(clientList[:i], clientList[i+1:]...)
						break
					}
				}
				manager.Clients[client.RoomId] = clientList
			case msg := <-manager.Broadcast:
				if clientList, ok := manager.Clients[msg.RoomId]; ok {
					for _, client := range clientList {
						client.Send <- msg
					}
				} else {
					log.Println("用户未在线")
				}
			}
		}
	}()
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.WriteMessage(websocket.CloseMessage, []byte("发送信息错误"))
		c.Socket.Close()
		fmt.Println("退出读")
	}()
	for {
		c.Socket.PongHandler()
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		msg := Message{}
		msg.Sender = c.Id
		msg.Content = string(message)
		msg.RoomId = c.RoomId
		msg.Ts = time.Now().Format("2006-01-02 15:04:05")
		Manager.Broadcast <- msg
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
		fmt.Println("退出写")
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			fmt.Println("send:", msg)
			if !ok {
				return
			}
			if err := c.Socket.WriteJSON(msg); err != nil {
				log.Println("send error:", err)
			}
		}
	}
}
