package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/ranaahsanansar/gochat/internal/websockets"
	"github.com/ranaahsanansar/gochat/middleware"
)

func ChatRoutes(r *gin.Engine, wshandler *websockets.Handler) {
	chatRoutes := r.Group("/chats")
	{
		chatRoutes.GET("/join-room/:roomId", wshandler.JoinRoom)
		chatRoutes.POST("/create-group", middleware.UserGuard, CreateGroup)
		chatRoutes.POST("/auto-create-group/:userId", middleware.UserGuard, AutoCreateGroup)
		chatRoutes.POST("/add-member/:groupId/:userId", middleware.UserGuard, AddMemberToGroup)
		chatRoutes.GET("/get-groups", middleware.UserGuard, GetGroups)
	}
}

// ws://localhost:3001/chats/join-room/123
