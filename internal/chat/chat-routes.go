package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/ranaahsanansar/gochat/middleware"
)

func ChatRoutes(r *gin.Engine) {
	chatRoutes := r.Group("/chats")
	{
		chatRoutes.GET("/join-room/:roomid", middleware.UserGuard, JoinRoom)
	}
}
