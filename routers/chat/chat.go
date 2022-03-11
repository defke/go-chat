package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-chat/util"
	"log"
	"net/http"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

/**开启连接**/
func Chat(c *gin.Context) {
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("err:", err)
	}
	roomId := c.Param("roomId")
	nickname := c.Param("nickname")
	client := &util.Client{
		RoomId: roomId,
		Id:     nickname,
		Socket: conn,
		Send:   make(chan util.Message),
	}
	util.Manager.Register <- client
	go client.Read()
	go client.Write()
}
