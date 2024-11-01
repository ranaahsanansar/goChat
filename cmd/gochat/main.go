package main

import (
	"time"

	"github.com/gin-contrib/cors"
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

	// Configure CORS to allow all origins
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	hub := websockets.NewHub()
	wshandler := websockets.NewHandler(hub)
	go hub.Run()
	// Module-wise routes
	user.RegisterRoutes(r)
	chat.ChatRoutes(r, wshandler)

	// Start server
	r.Run(":3000")
}
