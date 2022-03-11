package main

import (
	"github.com/gin-gonic/gin"
	"go-chat/routers"
	"go-chat/util"
)

/**
go的通讯聊天,基于频道号
*/
func main() {
	util.Manager.Start()
	router := gin.Default()
	routers.LoadRouter(router)
	router.Run(":8080")
}
