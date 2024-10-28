package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ranaahsanansar/gochat/initializers"
	"github.com/ranaahsanansar/gochat/internal/chat"
	"github.com/ranaahsanansar/gochat/internal/user"
	"github.com/ranaahsanansar/gochat/internal/websockets"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	hub := websockets.NewHub()
	wshandler := websockets.NewHandler(hub)
	go hub.Run()
	// Module-wise routes
	user.RegisterRoutes(r)
	chat.ChatRoutes(r, wshandler)

	// Start server
	r.Run(":3000")
}
