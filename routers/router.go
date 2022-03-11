package routers

import (
	"github.com/gin-gonic/gin"
	"go-chat/routers/chat"
)

func LoadRouter(app *gin.Engine) {
	//发送
	app.GET("/chat/:nickname/:roomId", chat.Chat)
}
