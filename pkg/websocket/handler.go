package websocket

import (
	"github.com/gin-gonic/gin"
)

var hubSingleton *Hub

func InitHub() {
	hubSingleton = NewHub()
	go hubSingleton.Run()
}

func ServeWs(ctx *gin.Context, roomId string) {
	serveWS(ctx, roomId, hubSingleton)
}
